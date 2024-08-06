package pullrequest

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot-cli/internal/id"
	"github.com/marqeta/pr-bot-cli/internal/opa"
)

//go:generate mockery --name EventHandler --testonly
type EventHandler interface {
	EvalAndReview(ctx context.Context, id id.PR, event *github.PullRequestEvent) error
}

type eventHandler struct {
	evaluator opa.Evaluator
}

func NewEventHandler(evaluator opa.Evaluator) EventHandler {
	return &eventHandler{
		evaluator: evaluator,
	}
}

func (eh *eventHandler) EvalAndReview(ctx context.Context, id id.PR, event *github.PullRequestEvent) error {
	//ghe := input.ToGHE(event)

	// test message
	fmt.Printf("evaluating PR. \n")
	fmt.Printf("event payload: %v\n", event)

	//tags := id.ToTags()
	//opaResult, err := eh.evaluator.Evaluate(ctx, ghe)
	//// todo: use zerolog
	////oplog.Err(err).Interface("decision", opaResult).Msg("opa evaluation complete")
	//if err != nil {
	//	fmt.Println("opa evaluation error")
	//	return err
	//}

	//if !opaResult.Track {
	//	fmt.Println("track=false, skipping review")
	//	return nil
	//}
	//
	//switch opaResult.Review.Type {
	//
	//case types.Approve:
	//	return eh.reviewer.Approve(ctx, id, opaResult.Review.Body, review.ApproveOptions{
	//		AutoMergeEnabled: event.GetPullRequest().GetBase().GetRepo().GetAllowAutoMerge(),
	//		DefaultBranch:    event.Repo.GetDefaultBranch(),
	//		MergeMethod:      mergeMethod(event),
	//	})
	//case types.RequestChanges:
	//	return eh.reviewer.RequestChanges(ctx, id, opaResult.Review.Body)
	//case types.Comment:
	//	return eh.reviewer.Comment(ctx, id, opaResult.Review.Body)
	//default:
	//	oplog.Info().Msg("skipping review")
	//}
	return nil
}

//func mergeMethod(event *github.PullRequestEvent) githubv4.PullRequestMergeMethod {
//	rebase := event.PullRequest.GetBase().GetRepo().GetAllowRebaseMerge()
//	squash := event.PullRequest.GetBase().GetRepo().GetAllowSquashMerge()
//	fc := event.PullRequest.GetChangedFiles()
//	// TODO: let policy specify what merge method to use.
//	// when rebasing empty commits on to main,
//	// no new commit is created, therefore no triggers would be fired.
//	// use squash to force a new commit to be created. when merging empty PRs
//	if rebase && fc > 0 {
//		return githubv4.PullRequestMergeMethodRebase
//	}
//	if squash {
//		return githubv4.PullRequestMergeMethodSquash
//	}
//	return githubv4.PullRequestMergeMethodMerge
//
//}
