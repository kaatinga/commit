# commit CLI Tool

`commit` is an automatic commit message generator that harnesses the power of AI to help you streamline your version
control workflow.

## Overview

Instead of racking your brains for a descriptive commit message, simply run the `commit` tool. It will automatically
stage all changed files in your git repository and craft a concise, meaningful commit message for you.

## Features:

- **AI-Powered**: Leverages OpenAI to generate commit messages.
- **Automatic Staging**: Stages all changed files in the git repository.
- **Customizable**: Works with user-provided keys and paths.

## Installation:

    go install github.com/commitdev/commit@latest

## Usage:

    commit [global options] command [command options] [arguments...]

### Global Options:

- `--key value`: Provide a valid key to work with ChatGPT.
- `--path value`: Specify the path to your git repository.
- `--help, -h`: Display help information.
- `--version, -v`: Print the current version of the tool.

## Environment Variable:

The tool can optionally read from the `OPENAI_API_KEY` environment variable to interface with the OpenAI API. Ensure you
set this variable if you don't intend to use the `--key` option directly.

## Example:

    commit --key abcdefg --path /home/user/my-repo

or just:

    commit

For further help or to view a list of commands, run:

    commit help
