# Agent

## Purpose
Hexg is a package that implements [Hexagonal Grids](https://www.redblobgames.com/grids/hexagons/) from Red Blob Games.

Package is written in Go and uses the stdlib.

We use "grid" instead of "map" because map is a reserved word in Go.

## Commands
- Run example: go run ./cmd/examples -output testdata/hexgrid.png
- Tests: `go test ./...`
- Run single test: `go test -v ./path/to/package -run TestName`
- Format code: `go fmt ./...`

## Code Style
- Standard Go formatting using `gofmt`
- Imports organized by stdlib first, then external packages
- Error handling: return errors to caller, log.Fatal only in main
- Function comments use Go standard format `// FunctionName does X`
- Variable naming follows camelCase
- File structure follows standard Go package conventions
