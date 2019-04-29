package main

import (
	"context"
	"fmt"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"math/rand"
)

func snap(containers []types.Container) []types.Container {
	sample := append([]types.Container(nil), containers...)
	rand.Shuffle(len(containers), func(i, j int) {
		sample[i], sample[j] = sample[j], sample[i]
	})
	sample = sample[:len(containers) / 2]
	return sample
}

func run() {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		panic("I know what it's like to lose.")
	}
	err = dockerCli.Initialize(cliflags.NewClientOptions())
	if err != nil {
		panic("I know what it's like to lose.")
	}
	cli := dockerCli.Client()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic("I know what it's like to lose.")
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
func main () {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Short: "",
			Long:  "",
			Use:   "thanos",
			Run: func(cmd *cobra.Command, args []string) {
				run()
			},
		}
		return cmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "rumpl",
		Version:       "0.5.0",
	})
}
