package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

func snap(containers []types.Container) []types.Container {
	sample := append([]types.Container(nil), containers...)
	rand.Shuffle(len(containers), func(i, j int) {
		sample[i], sample[j] = sample[j], sample[i]
	})
	sample = sample[:len(containers)/2]
	return sample
}

func run(dockerCli command.Cli) {
	cli := dockerCli.Client()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic("You should have gone for the head.")
	}

	s := snap(containers)
	for _, container := range s {
		err := dockerCli.Client().ContainerKill(context.Background(), container.ID, "9")
		if err != nil {
			panic("I know what it's like to lose.")
		}
	}

	fmt.Println("Fun isn’t something one considers when balancing the universe. But this… does put a smile on my face.")
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
			Run: func(cmd *cobra.Command, args []string) {
				run(dockerCli)
			},
		}
		return cmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "rumpl",
		Version:       "0.5.0",
	})
}
