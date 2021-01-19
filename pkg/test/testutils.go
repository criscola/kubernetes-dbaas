package test

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"unsafe"
)

// GetRootProjectPath is a utility function which gets the root project path based on the location of the main.go file.
// Do not move main.go in a place other than the root project path.
func GetRootProjectPath() (string, error) {
	var p string
	// Walk up the FS tree until a path containing main.go is found
	currentPath, _ := os.Getwd()
	mainFile := "main.go"

	for {
		p = path.Join(currentPath, mainFile)
		// Look if the current path contains main.go
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			// path of main.go found
			break
		}

		// Walk up the FS tree
		i := strings.LastIndex(currentPath, "/")
		currentPath = currentPath[:i]

		if currentPath == "" {
			return "", fmt.Errorf("path not found")
		}
	}
	return currentPath, nil
}

func GetTestdataFolderPath() (string, error) {
	rootPath, err := GetRootProjectPath()
	if err != nil {
		return "", fmt.Errorf("unable to get root path: %s", err)
	}
	return path.Join(rootPath, "testdata"), nil
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
