---
name: go-valtown
description: Unofficial Go client for the Val Town API, generated with oapi-codegen. Package `valtown` (`github.com/psyb0t/go-valtown`) calls api.val.town to manage vals, branches and files, environment variables and variable groups, blob storage, the account's SQLite database, email, telemetry logs and traces, users, orgs and aliases. Every endpoint is a typed `...WithResponse` method on `ClientWithResponses`, auth is a bearer API token set through `WithRequestEditorFn`, and responses decode into `JSON200` and per-status fields. Use when Go code needs to talk to Val Town, or already imports go-valtown.
homepage: https://github.com/psyb0t/go-valtown
user-invocable: true
metadata:
  openclaw:
    emoji: "🏙️"
    requires:
      bins:
        - go
      env:
        - VALTOWN_API_TOKEN
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the client sends requests only to the base URL you pass, https://api.val.town."
---

# go-valtown

Typed Go client for the Val Town API at api.val.town.

For install details, the full operation list and pagination, see [references/setup.md](references/setup.md).

## Security & safety

- Read the API token from the environment (`VALTOWN_API_TOKEN`) and set it in a request editor. Never hardcode it, commit it or log it, and never log the outgoing `Authorization` header.
- Requests go only to Val Town. Blob contents, SQL statements, file contents and environment variable values you send are stored in the user's Val Town account.
- Deleting vals, branches, files, blobs or environment variables, writing SQL, and sending email change the user's account and cannot be undone through this client. Call those endpoints only when the user asked for that change, and only on resources the task owns.

## When to use

- Go code needs to create, read, update or delete vals, val files or branches on Val Town.
- Go code needs Val Town blob storage, the account's SQLite database, or Val Town email.
- The project already imports `github.com/psyb0t/go-valtown`.

## When NOT to use

- Code running inside a val. Vals run on Deno and use Val Town's own `std/` modules, not Go.
- Code that is not Go. Call the HTTP API directly.

## Client and auth

```go
import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-valtown"
)

auth := valtown.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VALTOWN_API_TOKEN"))

	return nil
})

client, err := valtown.NewClientWithResponses("https://api.val.town", auth)
if err != nil {
	return err
}
```

## Who am I

```go
res, err := client.MeGetWithResponse(ctx)
if err != nil {
	return err
}

if res.JSON200 == nil {
	return fmt.Errorf("me: status %d: %s", res.StatusCode(), res.Body)
}

fmt.Println(res.JSON200.Id, *res.JSON200.Username)
```

## List vals

```go
privacy := valtown.ValsList2ParamsPrivacyPublic

res, err := client.ValsList2WithResponse(ctx, &valtown.ValsList2Params{
	Limit:   10,
	Privacy: &privacy,
})
if err != nil {
	return err
}

for _, val := range res.JSON200.Data { // []valtown.Val
	fmt.Println(val.Id, val.Name, val.Privacy)
}
```

`Limit` is required. `MeVals2WithResponse` lists only the token owner's vals.

## SQLite

```go
var statement valtown.SqliteExecuteJSONBody_Statement
if err := statement.FromSqliteExecuteJSONBodyStatement0("SELECT 1 AS one"); err != nil {
	return err
}

res, err := client.SqliteExecuteWithResponse(ctx, valtown.SqliteExecuteJSONRequestBody{
	Statement: statement,
})
// res.JSON200 is *valtown.ResultSet: Columns []string, Rows [][]interface{}
```

Use `FromSqliteExecuteJSONBodyStatement1` for SQL with bound arguments, and pass values that way instead of formatting them into the SQL string.

## Errors

A transport failure comes back as `err`. An HTTP error status does not: check `res.StatusCode()` or `res.JSON200 == nil` (`JSON201` for creates), then read the typed error field such as `res.JSON400` or `res.JSON404`, or the raw `res.Body`.
