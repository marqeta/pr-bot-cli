package pullrequest

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot-cli/internal/id"
	"strings"
)

var ErrEventActionNotFound = errors.New("event action was empty or nil")
var ErrMismatchedEvent = errors.New("expected pull_request event")
var ErrLabelNotFound = errors.New("pull_request.Label was nil")
var ErrPRNotFound = errors.New("event.pull_request was nil")

// Actions are used to identify registered callbacks.
const (
	// EventName is the event name of github.EventName's
	EventName = "pull_request"

	OpenedAction         = "opened"
	ReopenedAction       = "reopened"
	EditedAction         = "edited"
	LabeledAction        = "labeled"
	UnlabeledAction      = "unlabeled"
	ReviewRequested      = "review_requested"
	ReviewRequestRemoved = "review_request_removed"
	AssignedAction       = "assigned"
	UnassignedAction     = "unassigned"
	SynchronizeAction    = "synchronize"
)

//go:generate mockery --name Dispatcher
type Dispatcher interface {
	Dispatch(ctx context.Context, eventName string, event *github.PullRequestEvent) error
}

type EventFilter interface {
	ShouldHandle(ctx context.Context, id id.PR) (bool, error)
}

type dispatcher struct {
	handler EventHandler
	filter  EventFilter
}

func NewDispatcher(eh EventHandler, ef EventFilter) Dispatcher {
	return &dispatcher{
		handler: eh,
		filter:  ef,
	}
}

func (d *dispatcher) Dispatch(ctx context.Context, eventName string, event *github.PullRequestEvent) error {
	//oplog := httplog.LogEntry(ctx)
	fmt.Printf("Dispath event name %s\n", eventName)

	var err error

	if eventName != EventName {
		//oplog.Err(ErrMismatchedEvent).Send()
		return parseError(ctx, ErrMismatchedEvent)
	}

	if event == nil || event.Action == nil || len(*event.Action) == 0 {
		//oplog.Err(ErrEventActionNotFound).Send()
		return parseError(ctx, ErrEventActionNotFound)
	}

	if event.PullRequest == nil {
		//oplog.Err(ErrPRNotFound).Send()
		return parseError(ctx, ErrPRNotFound)
	}

	visibility := aws.ToString(event.Repo.Visibility)
	if visibility != "public" {
		return nil
	}

	action := *event.Action
	//httplog.LogEntrySetField(ctx, "action", action)
	//httplog.LogEntrySetField(ctx, "repo", *event.Repo.FullName)
	//httplog.LogEntrySetField(ctx, "pr", fmt.Sprint(*event.PullRequest.Number))
	//oplog = httplog.LogEntry(ctx)

	id := id.PR{
		Owner:        *event.Repo.Owner.Login,
		Repo:         *event.Repo.Name,
		Number:       *event.PullRequest.Number,
		NodeID:       *event.PullRequest.NodeID,
		RepoFullName: *event.Repo.FullName,
		Author:       *event.PullRequest.User.Login,
		URL:          *event.PullRequest.HTMLURL,
	}

	// todo: add event filter back
	//shouldHandle, err := d.filter.ShouldHandle(ctx, id)
	shouldHandle := true
	if err != nil {
		return err
	}

	if !shouldHandle {
		return nil
	}

	switch action {

	case OpenedAction:
		if strings.HasPrefix(id.Author, "svc-") {
			//d.metrics.EmitDist(ctx, "openedPRs", 1.0, id.ToTags())
			fmt.Printf("opened PR\n")
		}
		fallthrough

	case ReopenedAction, EditedAction, LabeledAction,
		UnlabeledAction, ReviewRequested, ReviewRequestRemoved,
		AssignedAction, UnassignedAction, SynchronizeAction:
		return d.handler.EvalAndReview(ctx, id, event)

	default:
		//oplog.Info().Msgf("No Handlers registered for Event: %s and Action: %s", eventName, action)
		fmt.Printf("No Handlers registered for Event: %s and Action: %s\n", eventName, action)
	}

	return nil
}

func parseError(ctx context.Context, err error) error {
	// todo: custom error
	return err
}
