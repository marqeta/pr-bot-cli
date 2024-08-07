package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"net/http"
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
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName == "" {
		fmt.Println("GITHUB_EVENT_NAME not set")
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

	// Get the PR number, repo owner and repo name from the event
	prNumber := event.GetPullRequest().GetNumber()
	repoOwner := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()

	fmt.Printf("Event name: %s\n, PR number: %d, owner: %s, repoName:%s \n", eventName, prNumber, repoOwner, repoName)

	// Set up the GitHub clients
	v3Client, _ := setupGHEClients()

	// Create a comment
	comment := &github.IssueComment{Body: github.String("👋 Thanks for opening this pull request! PR Bot will auto-approve if it can.")}
	_, _, err = v3Client.Issues.CreateComment(context.Background(), repoOwner, repoName, prNumber, comment)
	if err != nil {
		log.Error().Msgf("Error creating comment: %v", err)
		os.Exit(1)
	}
}

func setupGHEClients() (*github.Client, *githubv4.Client) {
	log.Info().Msg("Setting up GHE clients")
	tok := os.Getenv("GITHUB_TOKEN")
	if tok == "" {
		log.Error().Msg("GITHUB_TOKEN not set")
		os.Exit(1)
	}
	fmt.Printf("for testing: %s....%s\n", tok[:5], tok[len(tok)-5:])

	// Create a custom HTTP client with SSL verification disabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	tc.Transport = tr

	// Initialize v3 client
	v3 := github.NewClient(tc)

	// Initialize v4 client
	v4 := githubv4.NewClient(tc)
	return v3, v4
}
