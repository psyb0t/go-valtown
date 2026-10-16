# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `valtown_test`, and call the generated client against an `httptest` server.
- **TST2** Each test asserts the HTTP method, the path, the `Authorization: Bearer` header and the query or decoded JSON body it sends, then checks the decoded `JSON200`.
- **TST3** A spec patch that changes a generated type, such as the string enum fix, gets a test that uses the typed result.
- **TST4** `example_test.go` holds the README usage example verbatim. It talks to the live API, so it only has to compile.
- **TST5** `make test` runs `go test -race` in the dev image.
