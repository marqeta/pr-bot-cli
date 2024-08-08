package simple.allow_empty_commits

import future.keywords

# version of pr bot schema version the policies are authored in
schema := "v1"

default track := false

track if {
	input.pull_request.changed_files == 0
	input.event == "pull_request"

	not has_ignore_topic
}

has_ignore_topic if {
	# allows repo owners to skip evaluation from this module
	# convention is to follow the pattern pr-bot-<package name>-ignore
	# 50 char limit; valid chars: small case letters, numbers and hypens
	some topic in input.repository.topics
	topic == "pr-bot-mq-empty-commits-ignore"
}