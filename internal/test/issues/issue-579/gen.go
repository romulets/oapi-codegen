package issue579

//go:generate go run github.com/romulets/oapi-codegen/cmd/oapi-codegen --package=issue579 --generate=types,skip-prune --alias-types -o issue.gen.go spec.yaml
