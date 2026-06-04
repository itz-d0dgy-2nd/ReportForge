package utilities

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

/*
Close → Closes the underlying fsnotify watcher
*/
func (_watcher *Watcher) Close() {
	_watcher.watcher.Close()
}

/*
WatchForChanges → Watches report directory for file changes and triggers callback with changed file path
*/
func (_watcher *Watcher) WatchForChanges(_reportPath string, _fileCache *FileCache, _onFileChange func(_changedFile string)) error {
	var (
		debounceTimer = time.NewTimer(0)
		pendingFiles  = make(map[string]struct{})
		debounceMutex sync.Mutex
	)

	if errDirectories := addDirectoriesRecursively(_watcher.watcher, _reportPath); errDirectories != nil {
		return errDirectories
	}

	Logger.Info("watching for file changes - press Ctrl+C to stop")

	<-debounceTimer.C

	for {
		select {
		case watcherEvent, errWatcherEvent := <-_watcher.watcher.Events:
			if !errWatcherEvent {
				return nil
			}

			if watcherEvent.Op&fsnotify.Create != 0 {
				if directoryInfo, errStat := os.Stat(watcherEvent.Name); errStat == nil && directoryInfo.IsDir() {
					Logger.Debug("new directory detected, adding to watch", "dir", watcherEvent.Name)
					_watcher.watcher.Add(watcherEvent.Name)
				}
			}

			if watcherEvent.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if ignoreFiles(watcherEvent.Name) || _watcher.paused.Load() {
					continue
				}

				debounceMutex.Lock()
				pendingFiles[watcherEvent.Name] = struct{}{}
				debounceTimer.Reset(DebounceInterval)
				debounceMutex.Unlock()
			}

		case <-debounceTimer.C:
			debounceMutex.Lock()
			filesToProcess := make(map[string]struct{}, len(pendingFiles))
			for path := range pendingFiles {
				filesToProcess[path] = struct{}{}
			}
			pendingFiles = make(map[string]struct{})
			debounceMutex.Unlock()

			if len(filesToProcess) == 0 {
				continue
			}

			if _watcher.skipNextTrigger.CompareAndSwap(true, false) {
				continue
			}

			var configChanged bool
			for path := range filesToProcess {
				Logger.Info("file changed", "file", filepath.Base(path))
				_fileCache.ReloadFile(path)

				if filepath.Ext(path) == ".yml" {
					Logger.Debug("CONFIG FILE CHANGED - Skipping subsequet regeneration triggered by engine modified files")
					configChanged = true
				}
			}

			_onFileChange("")
			if configChanged {
				_watcher.skipNextTrigger.Store(true)
				debounceTimer.Reset(2 * time.Second)
			}
		case errWatcherEvent, errWatcherEventOpen := <-_watcher.watcher.Errors:
			if !errWatcherEventOpen {
				return nil
			}
			Logger.Error("watcher error", "error", errWatcherEvent)
		}
	}
}

/*
addDirectoriesRecursively → Recursively adds all subdirectories to the watcher
THIS IS DUMB FSNOTIFY! Waiting for https://github.com/fsnotify/fsnotify/issues/18
*/
func addDirectoriesRecursively(_watcher *fsnotify.Watcher, _path string) error {
	return filepath.WalkDir(_path, func(path string, entry os.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil {
			Logger.Warn("error accessing path", "path", path, "error", errAnonymousFunction)
			return nil
		}

		if !entry.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		if errAdd := _watcher.Add(path); errAdd != nil {
			Logger.Warn("could not watch directory", "path", path, "error", errAdd)
		} else {
			Logger.Debug("watching directory", "path", path)
		}

		return nil
	})
}

/*
ignoreFiles → Returns true if file should be ignored by the watcher
*/
func ignoreFiles(_path string) bool {
	base := filepath.Base(_path)

	if strings.HasPrefix(base, ".") {
		return true
	}

	if strings.HasSuffix(base, "~") || strings.HasSuffix(base, ".swp") || strings.HasSuffix(base, ".tmp") {
		return true
	}

	if filepath.Ext(base) != ".md" && filepath.Ext(base) != ".yml" {
		return true
	}

	if strings.Contains(filepath.ToSlash(_path), "/"+ScreenshotsDirectory+"/") {
		return true
	}

	return false
}

/*
NewWatcher → Creates a new Watcher instance
*/
func NewWatcher() (*Watcher, error) {
	watcher, errWatcher := fsnotify.NewWatcher()
	if errWatcher != nil {
		return nil, NewExternalError("failed to create file watcher", errWatcher)
	}
	return &Watcher{
		watcher: watcher,
	}, nil
}
