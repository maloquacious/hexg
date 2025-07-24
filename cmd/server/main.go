// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a web server for the hexg package
package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/maloquacious/hexg"
	"github.com/spf13/cobra"
)

//go:embed templates/*.gohtml
var templateFS embed.FS

var (
	host   string
	port   string
	server *http.Server
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

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Hexg web server",
	Long:  "A web server for the hexg hexagonal grid package",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexg: version %s\n", hexg.Version())
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  "Start the web server on the specified host and port",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind to")
	serveCmd.Flags().StringVar(&port, "port", "3000", "Port to bind to")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer() {
	mux := http.NewServeMux()

	// Web routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/corners", handleCorners)
	mux.HandleFunc("POST /neighbors", handleNeighbors)

	// API routes
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/version", handleVersion)
	mux.HandleFunc("/api/shutdown", handleShutdown)

	server = &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	fmt.Printf("Starting server on http://%s:%s\n", host, port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"version":"%s"}`, hexg.Version())
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"shutting down"}`)

	go func() {
		time.Sleep(100 * time.Millisecond) // Allow response to be sent
		fmt.Println("Shutdown requested via API")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: home\n", r.Method, r.URL)
	data := getPageData()

	t, err := template.ParseFS(templateFS, "templates/index.gohtml")
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

	t, err := template.ParseFS(templateFS, "templates/status.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCorners(w http.ResponseWriter, r *http.Request) {
	// Parse parameters
	sizeX, _ := strconv.ParseFloat(r.URL.Query().Get("sizeX"), 64)
	sizeY, _ := strconv.ParseFloat(r.URL.Query().Get("sizeY"), 64)
	originX, _ := strconv.ParseFloat(r.URL.Query().Get("originX"), 64)
	originY, _ := strconv.ParseFloat(r.URL.Query().Get("originY"), 64)

	// Set defaults if not provided
	if sizeX == 0 {
		sizeX = 1
	}
	if sizeY == 0 {
		sizeY = 1
	}

	corners := getHexagonCornersWithParams(sizeX, sizeY, originX, originY)

	data := struct{ Corners []CornerInfo }{Corners: corners}

	t, err := template.ParseFS(templateFS, "templates/corners.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleNeighbors(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request details for debugging
	log.Printf("[neighbors] Content-Type: %s", r.Header.Get("Content-Type"))
	
	// Parse form data (hx-include sends form data)
	if err := r.ParseForm(); err != nil {
		log.Printf("[neighbors] ParseForm error: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	
	log.Printf("[neighbors] Form values: %+v", r.Form)

	tnCoords := strings.TrimSpace(r.FormValue("tnCoords"))
	log.Printf("[neighbors] tn %q\n", tnCoords)
	if tnCoords == "" {
		log.Printf("[neighbors] tnCoords is empty")
		http.Error(w, "tnCoords parameter is required", http.StatusBadRequest)
		return
	}

	// Convert TribeNet coordinates to OffsetCoord
	offsetCoord, err := hexg.NewTribeNetCoord(tnCoords)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid TribeNet coordinates: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("[neighbors] tn %q: oc %+v\n", tnCoords, offsetCoord)

	// Convert to Hex (cube coordinates)
	centerHex := offsetCoord.ToCubeOdd()
	log.Printf("[neighbors] oc %+v: ch %+v\n", offsetCoord, centerHex)

	// Get neighbors
	neighbors := getNeighborsForHex(centerHex)

	data := struct {
		CenterHex string
		Neighbors []NeighborInfo
	}{
		CenterHex: fmt.Sprintf("%q", centerHex),
		Neighbors: neighbors,
	}

	t, err := template.ParseFS(templateFS, "templates/neighbors.gohtml")
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
	return getNeighborsForHex(h)
}

func getNeighborsForHex(h hexg.Hex) []NeighborInfo {
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
	return getHexagonCornersWithParams(1, 1, 0, 0)
}

func getHexagonCornersWithParams(sizeX, sizeY, originX, originY float64) []CornerInfo {
	l := hexg.NewLayoutFlat(hexg.NewPoint(sizeX, sizeY), hexg.NewPoint(originX, originY), false)
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
