# Agent

## Purpose
Hexg is a package that implements [Hexagonal Grids](https://www.redblobgames.com/grids/hexagons/) from Red Blob Games.

Package is written in Go and uses the stdlib.

We use "grid" instead of "map" because map is a reserved word in Go.

## Commands
- Build server: `go build -o ./dist/server ./cmd/server`
- Run server: `cd testdata && ../dist/server serve --host localhost --port 3033`
- Version: `./dist/server version`
- Health check: `curl http://localhost:3033/api/health`
- Shutdown server: `curl http://localhost:3033/api/shutdown`
- Tests: `go test ./...`
- Run single test: `go test -v ./path/to/package -run TestName`
- Format code: `go fmt ./...`

## Back End
- REST API server written in Go
- Includes /api/shutdown route to gracefully stop the server
- Uses Go templates

## Front End
- Built with HTMX, AlpineJS, Missing.css

## Code Style
- Standard Go formatting using `gofmt`
- Imports organized by stdlib first, then external packages
- Error handling: return errors to caller, log.Fatal only in main
- Function comments use Go standard format `// FunctionName does X`
- Variable naming follows camelCase
- File structure follows standard Go package conventions

## Architecture Notes
- Backend serves HTMX endpoints

## Bash Scripts
- Always use `${VARIABLE}` with curly braces for all variables
- Always quote variable references: "${VARIABLE}"
- Use `set -e` for early exit on errors
- Include descriptive echo statements with emoji for visual feedback
- Test endpoints in sequence with explicit validation
- Exit with error code on test failures
- Use curl with proper headers and jq for parsing responses
