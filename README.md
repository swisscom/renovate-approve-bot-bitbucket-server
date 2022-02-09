# renovate-approve-bot-bitbucket-server

## Introduction

This is a small application, similar to 
[renovatebot/renovate-approve-bot-bitbucket-cloud](https://github.com/renovatebot/renovate-approve-bot-bitbucket-cloud/)
to automatically approve PRs where the current user is added as a reviewer.  
  
The idea behind this tool is to be able to auto-approve PRs from [Renovate](https://renovatebot.com/)
so that, if they have `automerge` enabled they can be automerged by Renovate itself.  
  
This tool doesn't automatically merge any PR, it just approves them.

## Requirements

- Docker
- Make

## Building

### Docker Image

```bash
make REGISTRY= IMAGE=your-username/bb-approve-bot docker-build
make REGISTRY= IMAGE=your-username/bb-approve-bot docker-run
```

**Warning: the default configuration assumes one of our internal Docker registries**

### Locally

```
make build
./approve-bot
```

## Usage

```bash
Usage: approve-bot [--debug] --username USERNAME --password PASSWORD --endpoint ENDPOINT [--author-filter AUTHOR-FILTER]

Options:
  --debug, -D
  --username USERNAME, -u USERNAME [env: BITBUCKET_USERNAME]
  --password PASSWORD, -p PASSWORD [env: BITBUCKET_PASSWORD]
  --endpoint ENDPOINT, -e ENDPOINT [env: BITBUCKET_ENDPOINT]
  --author-filter AUTHOR-FILTER, -a AUTHOR-FILTER [env: BITBUCKET_AUTHOR_FILTER]
  --help, -h             display this help and exit

```

#### Quick Local Run

```bash
export BITBUCKET_USERNAME=your-username
read -s -r BITBUCKET_PASSWORD
# Type password and press enter
export BITBUCKET_PASSWORD
export BITBUCKET_ENDPOINT=https://bitbucket.example.com/rest
export BITBUCKET_AUTHOR_FILTER=renovate-bot # Only approve PRs created by this user
make REGISTRY= IMAGE=your-username/bb-approve-bot docker-run
```

