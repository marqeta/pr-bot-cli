package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot-cli/internal/githubclient"
	"github.com/marqeta/pr-bot-cli/internal/metrics"
	"github.com/marqeta/pr-bot/configstore"
	pgithub "github.com/marqeta/pr-bot/github"
	pid "github.com/marqeta/pr-bot/id"
	pmetrics "github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/oci"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/input/plugins"
	"github.com/marqeta/pr-bot/pullrequest"
	"github.com/marqeta/pr-bot/pullrequest/review"
	"github.com/open-policy-agent/opa/sdk"
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
	log.Info().Interface("PR", id).Msg("Parsed PR")

	log.Info().Msg("Setting up GHE clients")
	tok := os.Getenv("GITHUB_TOKEN")
	if tok == "" {
		log.Error().Msg("GITHUB_TOKEN not set")
		os.Exit(1)
	}
	v3Client, v4Client := githubclient.CreateGithubClients(cmd.Context(), tok)
	emitter := metrics.NewEmitter()
	ghAPI := githubclient.NewAPI(v3Client, v4Client, emitter)

	eventFilter, err := setupEventFilter(&pullrequest.RepoFilterCfg{Allowlist: []string{".*"}}, ghAPI)
	if err != nil {
		log.Error().Msgf("Error setting up event filter: %v", err)
		os.Exit(1)
	}

	shouldHandle, err := eventFilter.ShouldHandle(cmd.Context(), id)
	if err != nil {
		log.Error().Msgf("Error invoking event filter: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("ShouldHandle: %v", shouldHandle)

	oap := setUpOPAEvaluator(ghAPI)
	ghe := input.ToGHE(&event)
	opaResult, err := oap.Evaluate(cmd.Context(), ghe)
	log.Info().Interface("decision", opaResult).Msg("opa evaluation complete")

	reviewer := setupReviewer(ghAPI, emitter)
	err = reviewer.Comment(cmd.Context(), id, "👋 Thanks for opening this pull request! PR Bot will auto-approve if it can.")
	if err != nil {
		log.Error().Msgf("Error creating comment: %v", err)
		os.Exit(1)
	}

}

func setupEventFilter(cfg *pullrequest.RepoFilterCfg, api pgithub.API) (pullrequest.EventFilter, error) {
	err := cfg.Update()
	if err != nil {
		return nil, err
	}

	cs, err := configstore.NewInMemoryStore[*pullrequest.RepoFilterCfg](cfg)
	if err != nil {
		return nil, err
	}
	return pullrequest.NewRepoFilter(cs, api), nil
}

func setupReviewer(api pgithub.API, emitter pmetrics.Emitter) review.Reviewer {
	log.Info().Msg("Setting up reviewer")
	// todo: mutex -> dedup -> precond -> rate limited -> reviewer
	base := review.NewReviewer(api, emitter)
	precond := review.NewPreCondValidationReviewer(base)
	return precond
}

const (
	serviceName = "pr-bot"
	env         = "github-action"
	bundleRoot  = "/opt/app/bundles"
	bundleFile  = "pr-bot-policy.tar.gz"
)

func setUpOPAEvaluator(api pgithub.API) opa.Evaluator {
	log.Info().Msg("Setting up OPA evaluator")
	client := setUpOPAClient()
	log.Info().Msg("Successfully setup OPA client")
	modules := FindOPAModules()
	policy := setUpOPAPolicies(client)
	factory := setUpInputFactory(api)
	return opa.NewEvaluator(modules, policy, factory, nil)
}

func setUpOPAClient() client.Client {
	log.Info().Msg("Setting up OPA client")
	config := fmt.Sprintf(`
	{
	   "labels": {
	      "app": "%s",
	      "region": "us-east-2",
	      "environment": "%s"
	   },
	   "bundles": {
	      "local": {
	         "resource": "file:///%s/%s"
	      }
	   }
	}`, serviceName, env, bundleRoot, bundleFile)
	// TODO this can block indefinitely, use channel to signal completion and set timeout
	opaSDK, err := sdk.New(context.Background(), sdk.Options{
		ID:     fmt.Sprintf("%s-%s", serviceName, env),
		Config: bytes.NewReader([]byte(config)),
	})
	if err != nil {
		log.Err(err).Msg("Error creating OPA SDK client")
		os.Exit(1)
	}
	log.Info().Msg("Successfully created OPA SDK client")
	return client.NewClient(opaSDK)
}

func FindOPAModules() []string {
	log.Info().Msg("Finding OPA modules")
	reader := oci.NewReader()
	filepath := fmt.Sprintf("%s/%s", bundleRoot, bundleFile)
	dirs, err := reader.ListDirs(context.Background(), filepath)
	if err != nil {
		log.Err(err).Msg("Error reading OPA bundle directories")
		os.Exit(1)
	}
	modules := reader.FilterModules(context.Background(), dirs)
	log.Info().Interface("Modules", modules).Msg("Found OPA modules")
	return modules
}

func setUpOPAPolicies(opaClient client.Client) opa.Policy {
	log.Info().Msg("Setting up OPA policies")
	v1 := opa.NewV1Policy(opaClient)
	return opa.NewVersionedPolicy(
		map[string]opa.Policy{"v1": v1},
		opaClient,
	)
}
func setUpInputFactory(api pgithub.API) input.Factory {
	log.Info().Msg("Setting up input factory")
	branchProtection := plugins.NewBranchProtection(api)
	// 100KB size limit
	filesChanged := plugins.NewFilesChanged(api, 100*1000)
	return input.NewFactory(branchProtection, filesChanged)
}
