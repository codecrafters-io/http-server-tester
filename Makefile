.PHONY: release build run

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

release:
	git tag v$(next_version_number)
	git push origin main v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester

test:
	go test -v ./internal/

copy_course_file:
	hub api \
		repos/codecrafters-io/build-your-own-http-server/contents/course-definition.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils

test_with_server: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/scenarios/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\": \"connect-to-port\", \"tester_log_prefix\": \"stage-1\", \"title\": \"connect-to-port\"}, {\"slug\": \"respond-with-200\", \"tester_log_prefix\": \"stage-2\", \"title\": \"respond-with-200\"}, {\"slug\": \"respond-with-404\", \"tester_log_prefix\": \"stage-3\", \"title\": \"respond-with-404\"}, {\"slug\": \"respond-with-content\", \"tester_log_prefix\": \"stage-4\", \"title\": \"respond-with-content\"}, {\"slug\": \"parse-headers\", \"tester_log_prefix\": \"stage-5\", \"title\": \"parse-headers\"}, {\"slug\": \"concurrent-connections\", \"tester_log_prefix\": \"stage-6\", \"title\": \"concurrent-connections\"}, {\"slug\": \"get-file\", \"tester_log_prefix\": \"stage-7\", \"title\": \"get-file\"}, {\"slug\": \"post-file\", \"tester_log_prefix\": \"stage-8\", \"title\": \"post-file\"}]" \
	dist/main.out
