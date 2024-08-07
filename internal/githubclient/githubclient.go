package githubclient

import (
	"context"

	"github.com/google/go-github/v50/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func CreateGithubClients(ctx context.Context, tok string) (*github.Client, *githubv4.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok},
	)

	httpClient := oauth2.NewClient(ctx, ts)

	// Initialize v3 client
	v3 := github.NewClient(httpClient)

	// Initialize v4 client
	v4 := githubv4.NewClient(httpClient)
	return v3, v4
}
