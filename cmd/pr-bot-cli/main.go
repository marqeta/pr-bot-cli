package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var configFile string

func main() {
	rootCmd := &cobra.Command{
		Use:   "hello",
		Short: "PR-bot CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Config file: %s\n", configFile)
			fmt.Println("blah blah pr bot!")
		},
	}

	evaluateCmd := &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluate a PR",
		Run:   evaluatePullRequest,
	}

	rootCmd.AddCommand(evaluateCmd)

	pflag.StringVarP(&configFile, "config", "c", "", "Path to the configuration file")
	pflag.Parse()
	rootCmd.Flags().AddFlagSet(pflag.CommandLine)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func evaluatePullRequest(cmd *cobra.Command, args []string) {
	// get the GitHub event path
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		fmt.Println("GITHUB_EVENT_PATH not set")
		os.Exit(1)
	}

	// Read the event payload from the json file
	eventPayload, err := os.ReadFile(eventPath)
	if err != nil {
		fmt.Println("Failed to read event payload:", err)
		os.Exit(1)
	}

	// parse the event payload into a PullRequestEvent
	var event github.PullRequestEvent
	if err := json.Unmarshal(eventPayload, &event); err != nil {
		fmt.Println("Failed to unmarshal event payload:", err)
		os.Exit(1)
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")

	fmt.Printf("event name %s\n event payload: %v\n", eventName, event)
	// todo: call dispatcher
}
