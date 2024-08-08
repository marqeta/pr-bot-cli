package simple.allow_empty_commits

import future.keywords

lgtm := "LGTM!! :100: :tada:"

marker1 := "re-run"

marker2 := "pr-bot-please"

default review["type"] := "SKIP"

has_required_status_check if {
	count(input.plugins.base_branch_protection.required_status_checks.contexts) > 0
}

has_re_run_marker if {
	some label in input.pull_request.labels
	lower(label.name) == marker1
}

has_re_run_marker if {
	some label in input.pull_request.labels
	lower(label.name) == marker2
}

has_re_run_marker if {
	startswith(lower(input.pull_request.title), marker1)
}

has_re_run_marker if {
	startswith(lower(input.pull_request.title), marker2)
}

review["type"] := "APPROVE" if {
	# auto approve if: merging changes to default branch
	input.pull_request.base.ref == input.repository.default_branch

	# auto approve if: pr has empty commits
	input.pull_request.changed_files == 0
	input.pull_request.additions == 0
	input.pull_request.deletions == 0

	# verify there is required status check
	has_required_status_check

	# has re-run marker
	has_re_run_marker
}

review["body"] := lgtm if {
	review.type == "APPROVE"
}