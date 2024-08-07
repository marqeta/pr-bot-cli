package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
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

	setupGHEClients()

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

func setupGHEClients() (*github.Client, *githubv4.Client) {
	log.Info().Msg("Setting up GHE clients")
	tok := os.Getenv("INPUT_GITHUB_TOKEN")
	if tok == "" {
		log.Error().Msg("GITHUB_TOKEN not set")
		os.Exit(1)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok},
	)
	httpClient := oauth2.NewClient(context.Background(), ts)

	// Initialize v3 client
	v3 := github.NewClient(httpClient)
	// Test the connection by fetching the authenticated user's details
	user, _, err := v3.Users.Get(context.Background(), "")
	if err != nil {
		log.Error().Msgf("Error testing GitHub v3 client connection: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Connected to GitHub as user: %s", *user.Login)

	// Initialize v4 client
	v4 := githubv4.NewClient(httpClient)
	// Test the connection by fetching the authenticated user's details
	user, _, err = v3.Users.Get(context.Background(), "")
	if err != nil {
		log.Error().Msgf("Error testing GitHub v4 client connection: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Connected to GitHub as user: %s", *user.Login)

	return v3, v4
}
