package cache

import (
	"sync"

	"github.com/ravilmc/leo/goreact/internal/reactbuilder"
)

type Manager struct {
	serverBuilds             *serverBuilds
	clientBuilds             *clientBuilds
	routeIDToParentFile      *routeIDToParentFile
	parentFileToDependencies *parentFileToDependencies
}

func NewManager() *Manager {
	return &Manager{
		serverBuilds: &serverBuilds{
			builds: make(map[string]reactbuilder.BuildResult),
			lock:   sync.RWMutex{},
		},
		clientBuilds: &clientBuilds{
			builds: make(map[string]reactbuilder.BuildResult),
			lock:   sync.RWMutex{},
		},
		routeIDToParentFile: &routeIDToParentFile{
			reactFiles: make(map[string]string),
			lock:       sync.RWMutex{},
		},
		parentFileToDependencies: &parentFileToDependencies{
			dependencies: make(map[string][]string),
			lock:         sync.RWMutex{},
		},
	}
}

type serverBuilds struct {
	builds map[string]reactbuilder.BuildResult
	lock   sync.RWMutex
}

func (cm *Manager) GetServerBuild(filePath string) (reactbuilder.BuildResult, bool) {
	cm.serverBuilds.lock.RLock()
	defer cm.serverBuilds.lock.RUnlock()
	build, ok := cm.serverBuilds.builds[filePath]
	return build, ok
}

func (cm *Manager) SetServerBuild(filePath string, build reactbuilder.BuildResult) {
	cm.serverBuilds.lock.Lock()
	defer cm.serverBuilds.lock.Unlock()
	cm.serverBuilds.builds[filePath] = build
}

func (cm *Manager) RemoveServerBuild(filePath string) {
	cm.serverBuilds.lock.Lock()
	defer cm.serverBuilds.lock.Unlock()
	if _, ok := cm.serverBuilds.builds[filePath]; !ok {
		return
	}
	delete(cm.serverBuilds.builds, filePath)
}

type clientBuilds struct {
	builds map[string]reactbuilder.BuildResult
	lock   sync.RWMutex
}

func (cm *Manager) GetClientBuild(filePath string) (reactbuilder.BuildResult, bool) {
	cm.clientBuilds.lock.RLock()
	defer cm.clientBuilds.lock.RUnlock()
	build, ok := cm.clientBuilds.builds[filePath]
	return build, ok
}

func (cm *Manager) SetClientBuild(filePath string, build reactbuilder.BuildResult) {
	cm.clientBuilds.lock.Lock()
	defer cm.clientBuilds.lock.Unlock()
	cm.clientBuilds.builds[filePath] = build
}

func (cm *Manager) RemoveClientBuild(filePath string) {
	cm.clientBuilds.lock.Lock()
	defer cm.clientBuilds.lock.Unlock()
	if _, ok := cm.clientBuilds.builds[filePath]; !ok {
		return
	}
	delete(cm.clientBuilds.builds, filePath)
}

type routeIDToParentFile struct {
	reactFiles map[string]string
	lock       sync.RWMutex
}

func (cm *Manager) SetParentFile(routeID, filePath string) {
	cm.routeIDToParentFile.lock.Lock()
	defer cm.routeIDToParentFile.lock.Unlock()
	cm.routeIDToParentFile.reactFiles[routeID] = filePath
}

func (cm *Manager) GetRouteIDSForParentFile(filePath string) []string {
	cm.routeIDToParentFile.lock.RLock()
	defer cm.routeIDToParentFile.lock.RUnlock()
	var routes []string
	for route, file := range cm.routeIDToParentFile.reactFiles {
		if file == filePath {
			routes = append(routes, route)
		}
	}
	return routes
}

func (cm *Manager) GetAllRouteIDS() []string {
	cm.routeIDToParentFile.lock.RLock()
	defer cm.routeIDToParentFile.lock.RUnlock()
	routes := make([]string, 0, len(cm.routeIDToParentFile.reactFiles))
	for route := range cm.routeIDToParentFile.reactFiles {
		routes = append(routes, route)
	}
	return routes
}

func (cm *Manager) GetRouteIDSWithFile(filePath string) []string {
	reactFilesWithDependency := cm.GetParentFilesFromDependency(filePath)
	if len(reactFilesWithDependency) == 0 {
		reactFilesWithDependency = []string{filePath}
	}
	var routeIDS []string
	for _, reactFile := range reactFilesWithDependency {
		routeIDS = append(routeIDS, cm.GetRouteIDSForParentFile(reactFile)...)
	}
	return routeIDS
}

type parentFileToDependencies struct {
	dependencies map[string][]string
	lock         sync.RWMutex
}

func (cm *Manager) SetParentFileDependencies(filePath string, dependencies []string) {
	cm.parentFileToDependencies.lock.Lock()
	defer cm.parentFileToDependencies.lock.Unlock()
	cm.parentFileToDependencies.dependencies[filePath] = dependencies
}

func (cm *Manager) GetParentFilesFromDependency(dependencyPath string) []string {
	cm.parentFileToDependencies.lock.RLock()
	defer cm.parentFileToDependencies.lock.RUnlock()
	var parentFilePaths []string
	for parentFilePath, dependencies := range cm.parentFileToDependencies.dependencies {
		for _, dependency := range dependencies {
			if dependency == dependencyPath {
				parentFilePaths = append(parentFilePaths, parentFilePath)
			}
		}
	}
	return parentFilePaths
}
