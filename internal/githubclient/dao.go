package githubclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v50/github"
	"github.com/marqeta/pr-bot/errors"
	gh "github.com/marqeta/pr-bot/github"
	"github.com/marqeta/pr-bot/id"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/shurcooL/githubv4"
)

const (
	ApprovalTemplate = `
<details>
<summary>
%v
<br/>
</summary>

~~~json
%v
~~~

</details>
`
)

type SimpleDao struct {
	delegate gh.API
	v3       *github.Client
}

var _ gh.API = (*SimpleDao)(nil)

func NewAPI(v3 *github.Client, v4 *githubv4.Client, metrics metrics.Emitter) gh.API {
	return &SimpleDao{
		delegate: gh.NewAPI("", 0, v3, v4, metrics),
		v3:       v3,
	}
}

func (s *SimpleDao) ListReviews(ctx context.Context, id id.PR) ([]*github.PullRequestReview, error) {
	return s.delegate.ListReviews(ctx, id)
}

func (s *SimpleDao) AddReview(ctx context.Context, id id.PR, summary, event string) error {
	// override the base implementation
	msg := gh.ApprovalMessage{
		RequestID: middleware.GetReqID(ctx),
	}
	b, e := json.MarshalIndent(msg, "", "  ")
	if e != nil {
		return e
	}
	body := fmt.Sprintf(ApprovalTemplate, summary, string(b))
	_, _, err := s.v3.PullRequests.CreateReview(ctx, id.Owner, id.Repo, id.Number,
		&github.PullRequestReviewRequest{
			Body:  &body,
			Event: &event,
		})
	return err
}

func (s *SimpleDao) EnableAutoMerge(ctx context.Context, id id.PR, method githubv4.PullRequestMergeMethod) error {
	return s.delegate.EnableAutoMerge(ctx, id, method)
}

func (s *SimpleDao) IssueComment(ctx context.Context, id id.PR, comment string) error {
	return s.delegate.IssueComment(ctx, id, comment)
}

func (s *SimpleDao) IssueCommentForError(ctx context.Context, id id.PR, err errors.APIError) error {
	return s.delegate.IssueCommentForError(ctx, id, err)
}

func (s *SimpleDao) ListAllTopics(ctx context.Context, id id.PR) ([]string, error) {
	return s.delegate.ListAllTopics(ctx, id)
}

func (s *SimpleDao) ListRequiredStatusChecks(ctx context.Context, id id.PR, branch string) ([]string, error) {
	return s.delegate.ListRequiredStatusChecks(ctx, id, branch)
}

func (s *SimpleDao) ListFilesInRootDir(ctx context.Context, id id.PR, branch string) ([]string, error) {
	return s.delegate.ListFilesInRootDir(ctx, id, branch)
}

func (s *SimpleDao) ListFilesChangedInPR(ctx context.Context, id id.PR) ([]*github.CommitFile, error) {
	return s.delegate.ListFilesChangedInPR(ctx, id)
}

func (s *SimpleDao) GetBranchProtection(ctx context.Context, id id.PR, branch string) (*github.Protection, error) {
	return s.delegate.GetBranchProtection(ctx, id, branch)
}
