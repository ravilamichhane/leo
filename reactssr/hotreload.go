package go_ssr

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/ravilmc/leo/reactssr/packages/utils"
)

type HotReload struct {
	engine           *Engine
	connectedClients map[string][]*websocket.Conn
}

// newHotReload creates a new HotReload instance
func newHotReload(engine *Engine) *HotReload {
	return &HotReload{
		engine:           engine,
		connectedClients: make(map[string][]*websocket.Conn),
	}
}

// Start starts the hot reload server and watcher
func (hr *HotReload) Start() {
	go hr.startServer()
	go hr.startWatcher()
}

// startServer starts the hot reload websocket server
func (hr *HotReload) startServer() {
	slog.Info(fmt.Sprintf("Hot reload websocket running on port %d", hr.engine.Config.HotReloadServerPort))
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("Failed to upgrade websocket")
			return
		}
		// Client should send routeID as first message
		_, routeID, err := ws.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message from websocket")
			return
		}
		err = ws.WriteMessage(1, []byte("Connected"))
		if err != nil {
			slog.Error("Failed to write message to websocket")
			return
		}
		// Add client to connectedClients
		hr.connectedClients[string(routeID)] = append(hr.connectedClients[string(routeID)], ws)
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", hr.engine.Config.HotReloadServerPort), nil)
	if err != nil {
		slog.Error("Hot reload server quit unexpectedly")
	}
}

// startWatcher starts the file watcher
func (hr *HotReload) startWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to start watcher")
		return
	}
	defer watcher.Close()
	// Walk through all files in the frontend directory and add them to the watcher
	if err = filepath.Walk(hr.engine.Config.FrontendDir, func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		slog.Error("Failed to add files in directory to watcher")
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			// Watch for file created, deleted, updated, or renamed events
			if event.Op.String() != "CHMOD" && !strings.Contains(event.Name, "gossr-temporary") {
				filePath := utils.GetFullFilePath(event.Name)
				slog.Info("Frontend Files Changed Reloading....")
				// Store the routes that need to be reloaded
				var routeIDS []string
				switch {
				case filePath == hr.engine.Config.LayoutFilePath: // If the layout file has been updated, reload all routes
					if err := hr.engine.BuildLayoutCSSFile(); err != nil {
						slog.Error("Failed to build global css file")
					}
					routeIDS = hr.engine.CacheManager.GetAllRouteIDS()
				case hr.layoutCSSFileUpdated(filePath): // If the global css file has been updated, rebuild it and reload all routes
					if err := hr.engine.BuildLayoutCSSFile(); err != nil {
						slog.Error("Failed to build global css file")
						continue
					}
					routeIDS = hr.engine.CacheManager.GetAllRouteIDS()
				case hr.needsTailwindRecompile(filePath): // If tailwind is enabled and a React file has been updated, rebuild the global css file and reload all routes
					if err := hr.engine.BuildLayoutCSSFile(); err != nil {
						slog.Error("Failed to build global css file")
						continue
					}
					fallthrough
				default:
					// Get all route ids that use that file or have it as a dependency
					routeIDS = hr.engine.CacheManager.GetRouteIDSWithFile(filePath)
				}
				// Find any parent files that import the file that was modified and delete their cached build
				parentFiles := hr.engine.CacheManager.GetParentFilesFromDependency(filePath)
				for _, parentFile := range parentFiles {
					hr.engine.CacheManager.RemoveServerBuild(parentFile)
					hr.engine.CacheManager.RemoveClientBuild(parentFile)
				}
				// Reload any routes that import the modified file
				go hr.broadcastFileUpdateToClients(routeIDS)

			}
		case err := <-watcher.Errors:
			slog.Error("Error watching files", slog.Any("error", err))
		}
	}
}

// layoutCSSFileUpdated checks if the layout css file has been updated
func (hr *HotReload) layoutCSSFileUpdated(filePath string) bool {
	return utils.GetFullFilePath(filePath) == hr.engine.Config.LayoutCSSFilePath
}

// needsTailwindRecompile checks if the file that was updated is a React file
func (hr *HotReload) needsTailwindRecompile(filePath string) bool {
	if hr.engine.Config.TailwindConfigPath == "" {
		return false
	}

	log.Println(filePath)
	fileTypes := []string{".tsx", ".ts", ".jsx", ".js"}
	for _, fileType := range fileTypes {
		if strings.HasSuffix(filePath, fileType) {
			return true
		}
	}
	return false
}

// broadcastFileUpdateToClients sends a message to all connected clients to reload the page
func (hr *HotReload) broadcastFileUpdateToClients(routeIDS []string) {
	// Iterate over each route ID
	for _, routeID := range routeIDS {
		// Find all clients listening for that route ID
		for i, ws := range hr.connectedClients[routeID] {
			// Send reload message to client
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client if browser is closed or page changed
				hr.connectedClients[routeID] = append(hr.connectedClients[routeID][:i], hr.connectedClients[routeID][i+1:]...)
			}
		}
	}
}
