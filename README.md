# go-valtown

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-valtown.svg)](https://pkg.go.dev/github.com/psyb0t/go-valtown)
[![CI](https://github.com/psyb0t/go-valtown/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-valtown/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-valtown/badges/version.svg)](https://github.com/psyb0t/go-valtown/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-valtown/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-valtown/badges/importers.svg)](https://github.com/psyb0t/go-valtown/blob/badges/importers.md)

Unofficial Go client for the [Val Town API](https://docs.val.town/api/), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Val Town publishes. I am not affiliated with Val Town.

The package is `valtown`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Contents

- [Install](#install)
- [Usage](#usage)
- [Auth](#auth)
- [Regenerating](#regenerating)
- [Development](#development)
- [Agent integrations](#agent-integrations)
- [License](#license)

## Install

```sh
go get github.com/psyb0t/go-valtown
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields. The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-valtown) and in `go doc github.com/psyb0t/go-valtown`.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-valtown"
)

func main() {
	auth := valtown.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("VALTOWN_API_TOKEN"))

		return nil
	})

	client, err := valtown.NewClientWithResponses("https://api.val.town", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.MeGetWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

Listing vals takes a params struct, and enum values have typed constants:

```go
privacy := valtown.ValsList2ParamsPrivacyPublic

res, err := client.ValsList2WithResponse(ctx, &valtown.ValsList2Params{
	Limit:   10,
	Privacy: &privacy,
})
// res.JSON200.Data is []valtown.Val, res.JSON200.Links has the pagination links
```

Running SQL against your Val Town SQLite database:

```go
var statement valtown.SqliteExecuteJSONBody_Statement
if err := statement.FromSqliteExecuteJSONBodyStatement0("SELECT 1 AS one"); err != nil {
	panic(err)
}

res, err := client.SqliteExecuteWithResponse(ctx, valtown.SqliteExecuteJSONRequestBody{
	Statement: statement,
})
// res.JSON200 is a *valtown.ResultSet with Columns and Rows
```

A statement is either a plain SQL string (`FromSqliteExecuteJSONBodyStatement0`) or SQL with arguments (`FromSqliteExecuteJSONBodyStatement1`).

Method names come from the spec's operation IDs. Some carry a version number from the API path (`V1blobsList`, `V2blobsList`) and the v2 val endpoints end in `2` (`ValsList2`, `ValsGet2`, `BranchesCreate2`). Deleting a val is `ProjectsDelete2`. The three v3 endpoints without an operation ID are named after their path: `GetV3ValViewer`, `PostV3ConnectionsSlackToken` and `PostV3ConnectionsGoogleDocsToken`.

## Auth

Val Town authenticates with a bearer token: `Authorization: Bearer <token>`. Create an API token in the Val Town settings and set it with a request editor, as in the example above.

The base URL is `https://api.val.town`.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://api.val.town/openapi.json` to `api/openapi.json`.
- `scripts/fix-spec.sh` runs it through `scripts/fix-valtown-spec.jq` and writes `api/openapi.fixed.json`. The spec declares several string enums (privacy, val type, SQLite batch mode and others) with no `type`. oapi-codegen turns those into `interface{}` with no constants, so the filter adds `type: string` to every enum whose values are all strings.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64`, because the spec gives numbers no format and oapi-codegen would otherwise use `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## Agent integrations

The [skill](.agents/skills/go-valtown) teaches a coding agent how to use go-valtown when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-valtown@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-valtown@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-valtown
```

## License

[MIT](LICENSE). The generated code comes from Val Town's published API specification. The specification itself is not redistributed here.
