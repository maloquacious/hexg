# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

`hexg` is a single-package Go library implementing hexagonal grids following the
[Red Blob Games Hexagonal Grids guide](https://www.redblobgames.com/grids/hexagons/)
and its [implementation notes](https://www.redblobgames.com/grids/hexagons/implementation.html).
The guide is the specification — when behavior is ambiguous, match the guide.

The package uses only the stdlib plus `github.com/maloquacious/semver` (for `version.go`).

## Commands

```sh
go test ./...                    # all tests
go test -v -run TestHex_Distance # single test (root package)
go fmt ./...
go vet ./...

# regenerate example images
go run ./cmd/examples -output testdata/pointy.png -orientation pointy
go run ./cmd/examples -output testdata/flat.png   -orientation flat
go run ./cmd/examples -output testdata/pointy-art.png -orientation pointy -art-sprite
```

## Architecture

Everything lives in the root package `hexg` (files are split by concern, not by
sub-package). `cmd/examples` is the only other package — a PNG grid renderer used to
eyeball layouts.

### Coordinate systems

Three representations, with conversions between them:

- **Cube/axial** — `Hex{q, r, s}` in `hex.go`. Fields are unexported and the invariant
  `q + r + s == 0` is maintained by construction: `NewHex(q, r)` derives `s`. Read via
  `Q()`, `R()`, `S()`, `QRS()`. Because `Hex` is a comparable value type it is used
  directly as a map key — `HexSet` is `map[Hex]struct{}`, the return type of all grid
  generators.
- **Offset** — `OffsetCoord{Col, Row}` in `offset.go`, for storing grids in 2-D arrays.
  Conversions come in q-variants (flat-top) and r-variants (pointy-top), each taking an
  `even bool` to select even- or odd-offset.
- **Fractional** — `FractionalHex{q, r, s}` (float) in `fractional.go`, the intermediate
  for interpolation and pixel conversion. `Round()` snaps to the nearest `Hex` while
  preserving the sum-zero invariant by discarding the coordinate with the largest
  rounding error.

`Point{X, Y}` in `screen.go` is the screen/pixel type.

### Layout is the hub

`Layout` (`layout.go`) is the only type that knows about orientation, and most
grid-shaped and screen-related work goes through it.

- `LayoutOffset` is a single enum — `OddR`, `EvenR`, `OddQ`, `EvenQ` — that encodes
  *both* the hex orientation (R = pointy-top/horizontal rows, Q = flat-top/vertical
  columns) *and* the odd/even offset parity. There is no separate "orientation" enum in
  the public API; the unexported `orientation` struct holding the forward/inverse
  matrices is selected inside `NewLayout`.
- The `Is*` predicates (`IsPointy`, `IsFlat`, `IsOddR`, `IsEven`, …) are how the rest of
  the code branches on layout. `Layout.CubeToOffset` / `OffsetToCube` dispatch to the
  right `Hex`/`OffsetCoord` method for the layout, so callers rarely pass `even` bools
  themselves.
- `Layout.RotateLeft`/`RotateRight` are *not* pass-throughs: for Q layouts they invoke
  the opposite `Hex` rotation, so that "left" means counter-clockwise on screen for both
  orientations. Use the `Layout` methods when the visual result matters; the bare
  `Hex.RotateLeft`/`RotateRight` are pure coordinate shuffles.
- Screen conversion is matrix-based: `HexToPixel` uses `f0..f3`, `PixelToFractionalHex`
  uses the inverse `b0..b3`, and `startAngle` (0.0 flat, 0.5 pointy) drives
  `HexCornerOffset`/`PolygonCorners`.
- Grid generators (`Hexagon`, `Rectangle`, `Parallelogram*`, `Triangle*`) return
  `HexSet`. Only `Rectangle` actually consults the layout — it is built in offset space
  via `OffsetToCube`, so it is correct for all four offset types; the others are
  layout-independent but are methods for API consistency. Ring/spiral generation lives on `Hex` instead and returns ordered
  `[]Hex`, since order is meaningful there.

### Serialization

`ConciseString()` — `"+q+r+s"`, e.g. `"+2-1-1"` — is the canonical wire format.
`json.go` implements `json.Marshaler`/`Unmarshaler` and also `driver.Valuer`/`sql.Scanner`
over that same string, all funneling through `scanConciseString`, which re-validates the
sum-zero invariant on parse. `String()` (`"(q, r, s)"`) is for humans only; don't parse it.

## Conventions

- Copyright header on every file: `// Copyright (c) 2025 Michael D Henderson. All rights reserved.`
- Tests use the external test package `hexg_test` and exercise only the public API. They
  are table-driven with an `id` field on each case, iterating an inline anonymous struct
  slice in the `for ... range` clause, and report with
  `if want, got := ...; want != got { t.Errorf("%d: ...", tc.id, ...) }`.
- "grid" is used instead of "map" throughout, since `map` is a Go keyword.
- Invalid `LayoutOffset` values and negative radii panic rather than returning errors.
  Direction arguments are mostly coerced into 0..5 (`Neighbor`, `DiagonalNeighbor`,
  `DiagonalDirectionVector`); the exception is `DirectionVector`, which indexes directly
  and panics. Keep that split.
- `version.go` holds the package version. **Bump it in every change that touches code**
  — patch for a bug fix, minor for new API, major for a break. Documentation-only
  changes (README, CLAUDE.md, AGENT.md, comments) do not get a bump. Release commit
  messages carry the version, e.g. `Release v1.0.0`.
- **Every version bump gets a git tag, pushed right after the merge lands.**
  Tag the commit on `main` that carries the bump (for a squash-merged PR, that is
  the squash commit), using a lightweight tag named `v<major>.<minor>.<patch>`:
  `git tag v1.1.0 <sha> && git push origin v1.1.0`. Without the tag the release is
  invisible to `go get`, and the bumped version ships silently inside a later one.
- Generated PNGs are gitignored (`*.png`), so `testdata/` images are not committed by
  default.
