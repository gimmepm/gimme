# gimme

[![CircleCI](https://circleci.com/gh/gimmepm/gimme/tree/master.svg?style=svg)](https://circleci.com/gh/gimmepm/gimme/tree/master)

Gimme is end user software management for open source software. Created from a frustration of being blind to new releases for open source software that we like to use.

Here's an illustration of pre-gimme times...

1. **You find `awesome-opensource-software`** and it fits your need
1. You **download the latest release** of `awesome-opensource-software` (which doesn't have any published packages in package managers, or you want to be on the latest release without relying on package maintainers)
1. A week later, a **new release** of `awesome-opensource-software` comes out... but **you didn't know** that
1. More time passes, and **more releases are published... still you don't know**
1. Now you are running outdated (possibly insecure) software, all out of the "install-and-forget" pattern that OSS can introduce

Life with gimme...

1. **You find `awesome-opensource-software`** and it fits your need
1. **You star `awesome-opensource-software`** on GitHub
1. You **download the latest release** of `awesome-opensource-software`
1. A week later, a **new release** of `awesome-opensource-software` comes out... and **you run `gimme get updates`** and notice there is a new release
1. You **get the latest release** for `awesome-opensource-software`

## Installation

Gimme is released as both a bin ([Download here](https://github.com/gimmepm/gimme/releases)) and a [Docker image](https://hub.docker.com/r/gimmepm/gimme/tags/).

### Bin

1. Download the latest gimme release from GitHub
1. Extract the tarball
1. Copy/move the `gimme` bin to a location in your `$PATH` (e.g. `/usr/local/bin`)

### Container

*Note: this requires Docker to be installed locally*

Run the following:

```
$ docker run --rm -it \
    -e GIMME_GITHUB_TOKEN \
    -v /etc/localtime:/etc/localtime:ro \
    gimmepm/gimme get updates --since=-3days
```

## Configuration

Gimme requires a GitHub access token to make API calls. You can retrieve this by doing the following:

1. Navigate your browser to GitHub
1. Go into *Settings*
1. Click *Developer settings*
1. Go into *Personal access tokens*
1. Create a new token (requires no special access, just create the token)

You can either pass this token into gimme with the `--token` param, or you can set the `GIMME_GITHUB_TOKEN` environment variable to this token (recommended).

## Usage

### Get all updates in the past day

```
$ gimme get updates --since=-1days
```

### Get all repos that you have starred

```
$ gimme get repos
```

## Future/roadmap

- Caching on the local machine
- Updates since last retrieval
- Custom manipulation (add or delete) of watches repositories
- Support GitLab repos
- Automated install, update, and delete of the software
