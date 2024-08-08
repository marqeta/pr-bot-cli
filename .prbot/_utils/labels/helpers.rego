package _utils.labels

import future.keywords

# approve if pr-bot-please label is present
has_please_label if {
	some tage in input.pull_request.labels
	lower(tage.name) == "pr-bot-please"
}
