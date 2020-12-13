package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v4"
	"golang.org/x/sync/errgroup"
)

var (
	messages = []string{
		"I'm A Survivor",
		"The end is near",
		"It Needs Correction",
		"Simply Snap My Finger",
		"Fine. I'll do it myself",
		"I know what it's like to lose",
		"You should have gone for the head",
		"Your optimism is misplaced, Asgardian",
		"Perfectly balanced, as all things should be",
		"The hardest choices require the strongest wills",
		"Fun isn’t something one considers when balancing the universe. But this… does put a smile on my face.",
		"I used the stones to destroy the stones. And it nearly killed me. But the work is done, it always will be. I am inevitable",
	}
)

func snap(containers []types.Container) []types.Container {
	sample := append([]types.Container(nil), containers...)

	rand.Shuffle(len(containers), func(i, j int) {
		sample[i], sample[j] = sample[j], sample[i]
	})

	sample = sample[:len(containers)/2]

	return sample
}

func run(ctx context.Context) error {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	containers, err := c.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
	})
	if err != nil {
		return err
	}

	s := snap(containers)
	if len(containers) == 1 {
		s = containers
	}
	if len(s) == 0 {
		return nil
	}

	eg, ctx := errgroup.WithContext(ctx)

	p := mpb.New(mpb.WithWidth(len(s) * 2))
	bar := p.AddBar(int64(len(s)*2),
		mpb.BarStyle("[💀💀🐳]<+"),
	)

	for _, container := range s {
		container := container
		eg.Go(func() error {
			err := c.ContainerKill(ctx, container.ID, "9")
			bar.Increment(1)
			return err
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	bar.Abort(false)
	p.Wait()

	fmt.Println(messages[rand.Intn(len(messages))])

	return nil
}

func main() {
	root := cobra.Command{
		Use: "thanos",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context())
		},
	}

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
