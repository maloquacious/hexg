# Agent

## Purpose
Hexg is a package that implements [Hexagonal Grids](https://www.redblobgames.com/grids/hexagons/) from Red Blob Games.

Package is written in Go and uses the stdlib.

We use "grid" instead of "map" because map is a reserved word in Go.

## Commands
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

## V1 Release Plan

### Phase 5: Cube Folder Migration (COMPLETE - v0.15.0)
Migrated unique functionality from cube/ to root:
- [x] DiagonalDirectionVector(), DiagonalNeighbor() → hex.go
- [x] ReflectQ/R/S() → hex.go, fractional.go
- [x] Scale(), Ring(), Spiral() → hex.go
- [x] FractionalHex.Scale() → fractional.go
- [x] Deleted cube/ folder (stub files had no implementation)

### Phase 6: Implement cmd/examples (COMPLETE - v0.15.1)
- [x] Image generator with flags for orientation, size, origin
- [x] Art-sprite flag for art-friendly sizing (w/2, h/√3 or w/√3, h/2)
- [x] Render hex coordinates (q,r,s) on each hex
- [x] Output to PNG file

### Phase 7: Final
- [x] Run `go test ./...` and ensure all pass
- [x] Run `go fmt ./...`
- [ ] Tag v1.0.0