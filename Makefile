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

test_base: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/scenarios/pass_base \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\": \"at4\", \"tester_log_prefix\": \"stage-1\", \"title\": \"connect-to-port\"}, {\"slug\": \"ia4\", \"tester_log_prefix\": \"stage-2\", \"title\": \"respond-with-200\"}, {\"slug\": \"ih0\", \"tester_log_prefix\": \"stage-3\", \"title\": \"respond-with-404\"}, {\"slug\": \"cn2\", \"tester_log_prefix\": \"stage-4\", \"title\": \"respond-with-content\"}, {\"slug\": \"fs3\", \"tester_log_prefix\": \"stage-5\", \"title\": \"parse-headers\"}, {\"slug\": \"ej5\", \"tester_log_prefix\": \"stage-6\", \"title\": \"concurrent-connections\"}, {\"slug\": \"ap6\", \"tester_log_prefix\": \"stage-7\", \"title\": \"get-file\"}, {\"slug\": \"qv8\", \"tester_log_prefix\": \"stage-8\", \"title\": \"post-file\"}]" \
	dist/main.out

test_compression: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/scenarios/pass_base \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\": \"df4\", \"tester_log_prefix\": \"stage-9\", \"title\": \"content-encoding header\"}, {\"slug\": \"ij8\", \"tester_log_prefix\": \"stage-10\", \"title\": \"compression-multiple-schemes\"}, {\"slug\": \"cr8\", \"tester_log_prefix\": \"stage-11\", \"title\": \"compression-gzip\"}]" \
	dist/main.out

test_persistence: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/scenarios/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\": \"ag9\", \"tester_log_prefix\": \"stage-12\", \"title\": \"persistence-1\"}, {\"slug\": \"ul1\", \"tester_log_prefix\": \"stage-13\", \"title\": \"persistence-2\"}, {\"slug\": \"kh7\", \"tester_log_prefix\": \"stage-14\", \"title\": \"persistence-3\"}]" \
	dist/main.out

test_all: test_base test_compression test_persistence

test_release_locally:
	goreleaser release -f main.goreleaser.yml --clean --snapshot
