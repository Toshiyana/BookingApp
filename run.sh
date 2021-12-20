#!/bin/zsh

# before implement run.sh, implement "chmod +x run.sh" in terminal

# "go build" does not include test files, but "go run" include test files.
# a && b: when a was successd, implement b

go build -o go-bookingapp cmd/web/*.go && ./go-bookingapp