package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v4"
)

var (
	messages = [...]string{
		"The hardest choices require the strongest wills",
		"Perfectly balanced, as all things should be",
		"I used the stones to destroy the stones. And it nearly killed me. But the work is done, it always will be. I am inevitable",
		"I'm A Survivor",
		"You should have gone for the head",
		"It Needs Correction",
		"Simply Snap My Finger",
		"I know what it's like to lose",
		"Fine. I'll do it myself",
		"The end is near",
		"Fun isn’t something one considers when balancing the universe. But this… does put a smile on my face.",
		"Your optimism is misplaced, Asgardian",
	}
)

type ThanosError struct{}

func (t ThanosError) Error() string {
	return messages[rand.Intn(len(messages))]
}

func snap(containers []types.Container) []types.Container {
	sample := append([]types.Container(nil), containers...)
	rand.Shuffle(len(containers), func(i, j int) {
		sample[i], sample[j] = sample[j], sample[i]
	})
	sample = sample[:len(containers)/2]
	return sample
}

func kill(wg *sync.WaitGroup, dockerCli command.Cli, ID string, errCh chan error) {
	defer wg.Done()
	err := dockerCli.Client().ContainerKill(context.Background(), ID, "9")
	errCh <- err
}

func run(dockerCli command.Cli) error {
	cli := dockerCli.Client()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return ThanosError{}
	}

	s := snap(containers)
	if len(containers) == 1 {
		s = containers
	}
	if len(s) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWidth(len(s) * 2))
	bar := p.AddBar(int64(len(s)*2),
		mpb.BarStyle("[💀💀🐳]<+"),
	)
	errCh := make(chan error)
	for _, container := range s {
		wg.Add(1)
		go kill(&wg, dockerCli, container.ID, errCh)
	}
	hasError := false
	n := 0
	go func() {
		for v := range errCh {
			bar.Increment(1)
			n += 1
			if v != nil {
				hasError = true
			}
		}
	}()

	wg.Wait()
	close(errCh)
	// Wait a little bit so that the progess bar has the time to render everything
	time.Sleep(100 * time.Millisecond)
	bar.Abort(false)
	p.Wait()

	if hasError {
		return ThanosError{}
	}
	fmt.Println(messages[rand.Intn(len(messages))])
	return nil
}

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Short: "",
			Long:  "",
			Use:   "thanos",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				if err := plugin.PersistentPreRunE(cmd, args); err != nil {
					return err
				}
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return run(dockerCli)
			},
		}
		return cmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "rumpl",
		Version:       "1.0.0",
	})
}
