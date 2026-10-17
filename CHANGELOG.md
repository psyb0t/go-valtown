# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.0, 2026-10-16

- Added `valtown`, a REST client generated with oapi-codegen from Val Town's OpenAPI document at api.val.town (Val Town API, spec version 1, OpenAPI 3.1.0). It covers all 35 paths: vals, branches, files, environment variables and variable groups, blobs, SQLite, email, telemetry, users, orgs and aliases.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec from api.val.town, gives untyped string enums `type: string` so they get typed constants, and regenerates the client. The spec is not committed.
- Added the `go-valtown` agent skill. Install it with `claude plugin install go-valtown@psyb0t` (Claude Code), `codex plugin add go-valtown@psyb0t` (Codex) or `openclaw skills install @psyb0t/go-valtown` (ClawHub).
