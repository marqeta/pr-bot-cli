package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot-cli/internal/githubclient"
	"github.com/marqeta/pr-bot-cli/internal/metrics"
	pgithub "github.com/marqeta/pr-bot/github"
	pid "github.com/marqeta/pr-bot/id"
	pmetrics "github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var configFile string

func main() {
	rootCmd := &cobra.Command{
		Use:   "hello",
		Short: "PR-bot CLI",
		Run: func(_ *cobra.Command, _ []string) {
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

func evaluatePullRequest(cmd *cobra.Command, _ []string) {
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

	id := pid.PR{
		Owner:        event.GetRepo().GetOwner().GetLogin(),
		Repo:         event.GetRepo().GetName(),
		Number:       event.GetPullRequest().GetNumber(),
		NodeID:       event.GetPullRequest().GetNodeID(),
		RepoFullName: event.GetRepo().GetFullName(),
		Author:       event.GetPullRequest().GetUser().GetLogin(),
		URL:          event.GetPullRequest().GetHTMLURL(),
	}
	fmt.Printf("Event name: %s\n, PR number: %d, owner: %s, repoName:%s \n", eventName, id.Number, id.Owner, id.RepoFullName)

	log.Info().Msg("Setting up GHE clients")
	tok := os.Getenv("GITHUB_TOKEN")
	if tok == "" {
		log.Error().Msg("GITHUB_TOKEN not set")
		os.Exit(1)
	}
	v3Client, v4Client := githubclient.CreateGithubClients(cmd.Context(), tok)
	emitter := metrics.NewEmitter()
	ghAPI := pgithub.NewAPI("localhost", 8080, v3Client, v4Client, emitter)
	reviewer := setupReviewer(ghAPI, emitter)

	err = reviewer.Comment(cmd.Context(), id, "👋 Thanks for opening this pull request! PR Bot will auto-approve if it can.")
	if err != nil {
		log.Error().Msgf("Error creating comment: %v", err)
		os.Exit(1)
	}
}

func setupReviewer(api pgithub.API, emitter pmetrics.Emitter) review.Reviewer {
	log.Info().Msg("Setting up reviewer")
	// todo: mutex -> dedup -> precond -> rate limited -> reviewer
	base := review.NewReviewer(api, emitter)
	precond := review.NewPreCondValidationReviewer(base)
	return precond
}
