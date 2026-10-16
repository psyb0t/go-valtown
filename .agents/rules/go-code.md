# Go code

Refines `~/.claude/rule-details/working-with-generated-code.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` is generated. Never edit it by hand. Change the spec patch or `config.yaml` and run `make generate`.
- **GO2** Every fix to the upstream spec lives in `scripts/fix-valtown-spec.jq`, run by `scripts/fix-spec.sh` from `gen.go`. Each rule in the filter carries a comment that says what oapi-codegen does wrong without it. Patches stay generic: match on schema shape, never on one schema name from today's spec.
- **GO3** The generation steps are declared in `gen.go`. A new generated package gets its own `gen.go`, not a Makefile line.
- **GO4** `config.yaml` keeps the `number` to `float64` mapping.
- **GO5** `make lint` passes with zero issues. Lines are at most 80 columns, the `.golangci.yml` limit; split calls one argument per line rather than adding `//nolint`.
- **GO6** Hand-written Go in this repo is tests and examples only. The package carries no hand-written helpers on top of the generated code.
