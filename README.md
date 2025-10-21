[![Tests](https://github.com/kaatinga/commit/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kaatinga/commit/actions/workflows/test.yml)
[![GitHub release](https://img.shields.io/github/release/kaatinga/commit.svg)](https://github.com/kaatinga/commit/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/commit/blob/main/LICENSE)
[![lint workflow](https://github.com/kaatinga/commit/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/commit/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/commit/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# commit CLI Tool

`commit` is an AI-powered commit message generator that helps you streamline your version control workflow.

## Overview

Instead of racking your brains for a descriptive commit message, simply run the `commit` tool. It will automatically
stage all changed files in your git repository and craft a concise, meaningful commit message for you.

## Features:

- Leverages OpenAI to generate commit messages based on your code changes.
- Automatically stages all changed files in the git repository.
- Commits with the AI-generated message.
- Optional `push` command to commit and push in one step.
- Dry-run mode to preview commit messages without committing.
- Works with user-provided API keys and repository paths.

## Installation:

    go install github.com/kaatinga/commit@latest

## Usage:

    commit [global options] command [command options] [arguments...]

### Commands:

- `commit` (default): Generate a commit message and commit changes.
- `push`: Generate a commit message, commit changes, and push to remote.

### Global Options:

- `--key value`: Provide a valid key to work with OpenAI API.
- `--path value`: Specify the path to your git repository (default: current directory).
- `--help, -h`: Display help information.
- `--version, -v`: Print the current version of the tool.

## Environment Variables:

- `OPENAI_API_KEY`: The tool can optionally read from this environment variable to interface with the OpenAI API. Set this variable if you don't intend to use the `--key` option directly.
- `COMMIT_DRYRUN`: Set to `true` to enable dry-run mode, which will display the generated commit message without actually committing.

## Examples:

### Basic usage (uses OPENAI_API_KEY environment variable):

    commit

### With explicit API key:

    commit --key sk-your-api-key-here

### Specify a different repository path:

    commit --path /home/user/my-repo

### Commit and push in one command:

    commit push

### Dry-run mode (preview commit message without committing):

    export COMMIT_DRYRUN=true
    commit

For further help or to view a list of commands, run:

    commit help
