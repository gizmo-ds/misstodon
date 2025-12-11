[![爱发电](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fafdian.com%2Fapi%2Fuser%2Fget-profile%3Fuser_id%3D75e549844b5111ed8df552540025c377&query=%24.data.user.name&label=%E7%88%B1%E5%8F%91%E7%94%B5&color=%23946ce6)](https://afdian.com/a/gizmo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gizmo-ds/misstodon?style=flat-square)
[![Build images](https://img.shields.io/github/actions/workflow/status/gizmo-ds/misstodon/images.yaml?branch=main&label=docker%20image&style=flat-square)](https://github.com/gizmo-ds/misstodon/actions/workflows/images.yaml)
[![License](https://img.shields.io/github/license/gizmo-ds/misstodon?style=flat-square)](./LICENSE)

# Contribute to misstodon

## Prepare your environment

Prerequisites:

- [Go 1.20+](https://go.dev/doc/install)
- [git](https://git-scm.com/)
- [Bun](https://bun.sh/docs/installation) or [Deno](https://deno.land/manual/getting_started/installation) or [Node.js](https://nodejs.org/)

## Fork and clone misstodon

First you need to fork this project on GitHub.

Clone your fork:

```shell
git clone git@github.com:<you>/misstodon.git
```

## Prerequisite before build

### mfm.js

Misskey uses a non-standard Markdown implementation, which they call MFM (Misskey Flavored Markdown).

> Misskey has an [mfm.rs](https://github.com/misskey-dev/mfm.rs). If it's completed, I will attempt to compile it to WebAssembly and replace the current implementation.

If you are using [Bun](https://bun.sh/docs/installation):

```shell
bun install
bun run build
```

If you are using [Deno](https://deno.land/manual/getting_started/installation):

```shell
deno task build
```

If you are using [Node.js](https://nodejs.org/):

```shell
corepack prepare pnpm@latest --activate
corepack enable
pnpm install
pnpm run build
```

## Test your change

Currently, misstodon lacks proper unit tests. You can create test cases in the `pkg/misstodon/provider/misskey` directory.

```shell
go test github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey -v -run TestTimelinePublic
```

Another simple approach is to run the misstodon server and use tools like [Insomnia](https://insomnia.rest/) to test the API.

Start the misstodon server:

```shell
cp config_example.toml config.toml
go run cmd/misstodon/main.go start --fallbackServer=misskey.io
```

Request the API:

```shell
curl --request GET --url "http://localhost:3000/nodeinfo/2.0" | jq .
```

## Create a commit

Commit messages should be well formatted, and to make that "standardized", we are using Conventional Commits.

You can follow the documentation on [their website](https://www.conventionalcommits.org).

## Submit a pull request

Push your branch to your `misstodon` fork and open a pull request against the main branch.
