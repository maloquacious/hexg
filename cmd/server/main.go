// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a web server for the hexg package
package main

import (
	"fmt"
	"github.com/maloquacious/hexg"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type PageData struct {
	Version     string
	GitStatus   string
	GitDirty    bool
	CurrentTime string
	Neighbors   []NeighborInfo
	Corners     []CornerInfo
}

type NeighborInfo struct {
	Direction int
	Hex       string
}

type CornerInfo struct {
	Corner int
	X      float64
	Y      float64
}

func main() {
	fmt.Printf("hexg: version %s\n", hexg.Version())

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/status", handleStatus)

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := getPageData()
	
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hexg Package Info</title>
    <link rel="stylesheet" href="https://unpkg.com/missing.css@1.1.1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://unpkg.com/alpinejs@3.13.5/dist/cdn.min.js" defer></script>
</head>
<body>
    <main>
        <h1>Hexg Package Information</h1>
        
        <div x-data="{ refreshing: false }">
            <section>
                <h2>Package Version</h2>
                <p><strong>{{.Version}}</strong></p>
            </section>

            <section>
                <h2>Git Status</h2>
                <div id="git-status">
                    <p>Status: <span class="{{if .GitDirty}}status-dirty{{else}}status-clean{{end}}">
                        {{if .GitDirty}}Dirty{{else}}Clean{{end}}
                    </span></p>
                    {{if .GitStatus}}<pre>{{.GitStatus}}</pre>{{end}}
                </div>
                <button 
                    hx-get="/status" 
                    hx-target="#git-status"
                    x-on:click="refreshing = true"
                    x-on:htmx:after-request="refreshing = false"
                    x-bind:disabled="refreshing"
                >
                    <span x-show="!refreshing">Refresh Git Status</span>
                    <span x-show="refreshing">Refreshing...</span>
                </button>
            </section>

            <section>
                <h2>Origin Cell Neighbors</h2>
                <p>Neighbors of hex (0,0,0):</p>
                <table>
                    <thead>
                        <tr>
                            <th>Direction</th>
                            <th>Coordinates</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Neighbors}}
                        <tr>
                            <td>{{.Direction}}</td>
                            <td>{{.Hex}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </section>

            <section>
                <h2>Hexagon Corners</h2>
                <div style="display: flex; gap: 2rem; align-items: flex-start;">
                    <div>
                        <p>Corner locations for flat layout:</p>
                        <table>
                            <thead>
                                <tr>
                                    <th>Corner</th>
                                    <th>Point (x, y)</th>
                                </tr>
                            </thead>
                            <tbody>
                                {{range .Corners}}
                                <tr>
                                    <td>{{.Corner}}</td>
                                    <td>({{.X}}, {{.Y}})</td>
                                </tr>
                                {{end}}
                            </tbody>
                        </table>
                    </div>
                    <div>
                        <p>Hexagon visualization:</p>
                        <svg width="200" height="200" viewBox="-3 -3 6 6" style="border: 1px solid #ccc; background: #f9f9f9;">
                            <polygon 
                                points="{{range $i, $corner := .Corners}}{{if $i}}, {{end}}{{$corner.X}},{{$corner.Y}}{{end}}"
                                fill="none" 
                                stroke="#333" 
                                stroke-width="0.1"
                            />
                            {{range .Corners}}
                            <circle cx="{{.X}}" cy="{{.Y}}" r="0.1" fill="#666"/>
                            <text x="{{.X}}" y="{{.Y}}" dy="-0.2" text-anchor="middle" font-size="0.3" fill="#333">{{.Corner}}</text>
                            {{end}}
                        </svg>
                    </div>
                </div>
            </section>
        </div>
    </main>

    <footer>
        <p>Current time (UTC): {{.CurrentTime}}</p>
    </footer>

    <style>
        .status-clean { color: green; font-weight: bold; }
        .status-dirty { color: orange; font-weight: bold; }
        pre { background: #f5f5f5; padding: 1rem; border-radius: 4px; overflow-x: auto; }
        footer { margin-top: 2rem; padding-top: 1rem; border-top: 1px solid #eee; text-align: center; color: #666; }
    </style>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	data := getPageData()
	tmpl := `<p>Status: <span class="{{if .GitDirty}}status-dirty{{else}}status-clean{{end}}">
        {{if .GitDirty}}Dirty{{else}}Clean{{end}}
    </span></p>
    {{if .GitStatus}}<pre>{{.GitStatus}}</pre>{{end}}`

	t, err := template.New("status").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPageData() PageData {
	version := hexg.Version().String()
	gitStatus, gitDirty := getGitStatus()
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
	neighbors := getOriginNeighbors()
	corners := getHexagonCorners()

	return PageData{
		Version:     version,
		GitStatus:   gitStatus,
		GitDirty:    gitDirty,
		CurrentTime: currentTime,
		Neighbors:   neighbors,
		Corners:     corners,
	}
}

func getOriginNeighbors() []NeighborInfo {
	h := hexg.NewHex(0, 0, 0)
	var neighbors []NeighborInfo

	for _, direction := range []int{0, 1, 2, 3, 4, 5} {
		n := h.Neighbor(direction)
		neighbors = append(neighbors, NeighborInfo{
			Direction: direction,
			Hex:       fmt.Sprintf("%q", n),
		})
	}

	return neighbors
}

func getHexagonCorners() []CornerInfo {
	l := hexg.NewLayoutFlat(hexg.NewPoint(1, 1), hexg.NewPoint(0, 0), false)
	var corners []CornerInfo
	
	for corner, point := range l.PolygonCorners() {
		corners = append(corners, CornerInfo{
			Corner: corner,
			X:      point.X,
			Y:      point.Y,
		})
	}
	
	return corners
}

func getGitStatus() (string, bool) {
	var statusInfo []string
	var isDirty bool

	// Get build-time VCS information from runtime/debug
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if setting.Value != "" {
					statusInfo = append(statusInfo, fmt.Sprintf("Build commit: %s", setting.Value[:min(8, len(setting.Value))]))
				}
			case "vcs.modified":
				if setting.Value == "true" {
					statusInfo = append(statusInfo, "Working tree was dirty at build time")
					isDirty = true
				} else {
					statusInfo = append(statusInfo, "Working tree was clean at build time")
				}
			case "vcs.time":
				if setting.Value != "" {
					if buildTime, err := time.Parse(time.RFC3339, setting.Value); err == nil {
						statusInfo = append(statusInfo, fmt.Sprintf("Build time: %s", buildTime.UTC().Format("2006-01-02 15:04:05 UTC")))
					} else {
						statusInfo = append(statusInfo, fmt.Sprintf("Build time: %s", setting.Value))
					}
				}
			}
		}
	}

	if len(statusInfo) == 0 {
		statusInfo = append(statusInfo, "No VCS information available")
	}

	return strings.Join(statusInfo, "\n"), isDirty
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


