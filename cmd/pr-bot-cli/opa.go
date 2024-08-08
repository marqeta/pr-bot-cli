package main

import (
	"bytes"
	"context"
	"fmt"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/oci"
	"github.com/marqeta/pr-bot/opa"
	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/input/plugins"
	"github.com/open-policy-agent/opa/sdk"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	serviceName = "pr-bot"
	env         = "github-action"
	bundleRoot  = "/opt/app/bundles"
	bundleFile  = "bundle.tar.gz"
)

func setUpOPAEvaluator(api gh.API) opa.Evaluator {
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
func setUpInputFactory(api gh.API) input.Factory {
	log.Info().Msg("Setting up input factory")
	branchProtection := plugins.NewBranchProtection(api)
	// 100KB size limit
	filesChanged := plugins.NewFilesChanged(api, 100*1000)
	return input.NewFactory(branchProtection, filesChanged)
}
