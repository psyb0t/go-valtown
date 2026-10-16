// Package valtown is a Go client for the Val Town API, generated with
// oapi-codegen from the OpenAPI document Val Town publishes.
package valtown

//go:generate bash scripts/fetch-spec.sh https://api.val.town/openapi.json openapi.json
//go:generate bash scripts/fix-spec.sh fix-valtown-spec.jq openapi.json openapi.fixed.json
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.json
