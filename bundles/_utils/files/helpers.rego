package _utils.files

import future.keywords

# true if every file changed in the PR are in the list of files_in_scope.
# Does a case insenstive match on filenames
# Args:
# files_in_scope: array of allowed filenames
# files_changed: array from input.plugins.files_changed
# Example:
# valid_file_change(["one.go", "two.go","three.go"], [{"filename":"one.og"...}{"filename":"three.go"}]) == true
#
valid_file_change(files_in_scope, files_changed) if {
	lower_files_in_scope := [l |
		some f in files_in_scope
		l := lower(f)
	]

	every file in files_changed {
		filename := lower(file.filename)
		some file_in_scope in files_in_scope
		contains(filename, file_in_scope)
	}
}

valid_file_change_match_case(files_in_scope, files_changed) if {
	every file in files_changed {
		filename := file.filename
		some file_in_scope in files_in_scope
		contains(filename, file_in_scope)
	}
}
