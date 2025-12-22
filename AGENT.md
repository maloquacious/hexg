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

### Phase 2: Clean Up Root
- [ ] Review and organize exports
- [ ] Ensure consistent naming conventions
- [ ] Remove unused code or dead paths

### Phase 3: Documentation
- [ ] Add Go doc comments to all exported functions in hex.go
- [ ] Add Go doc comments to all exported functions in conversions.go
- [ ] Add Go doc comments to all exported functions in fractional.go
- [ ] Add Go doc comments to all exported functions in layout.go
- [ ] Add Go doc comments to all exported functions in math.go
- [ ] Add Go doc comments to all exported functions in offset.go
- [ ] Add Go doc comments to all exported functions in screen.go
- [ ] Add Go doc comments to all exported functions in json.go

### Phase 4: README
- [ ] Create README.md with:
  - Package overview and purpose
  - Installation instructions
  - Basic usage examples
  - Link to Red Blob Games hexagon guide
  - Acknowledgements section (Red Blob Games, Amp)
  - License info

### Phase 5: Cube Folder Migration
Migrate all unique functionality from cube/ to root:
- [ ] coords.go
- [ ] distances.go (+ tests)
- [ ] field_of_view.go
- [ ] geometry.go
- [ ] hex_to_pixel.go
- [ ] lines.go
- [ ] math.go
- [ ] movement.go
- [ ] neighbors.go (+ tests)
- [ ] paths.go
- [ ] pixel_to_hex.go
- [ ] references.go
- [ ] reflection.go
- [ ] rings.go
- [ ] rotation.go
- [ ] rounding.go
- [ ] storage.go
- [ ] strings.go
- [ ] wrapping.go

### Phase 6: Final
- [ ] Run `go test ./...` and ensure all pass
- [ ] Run `go fmt ./...`
- [ ] Remove cube/ folder after migration complete
- [ ] Tag v1.0.0