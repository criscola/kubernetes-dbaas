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

// GetTestdataFolderPath is a utility function which tries to get the "testdata" folder which should contain data files
// used for testing.
func GetTestdataFolderPath() (string, error) {
	rootPath, err := GetRootProjectPath()
	if err != nil {
		return "", fmt.Errorf("unable to get root path: %s", err)
	}
	return path.Join(rootPath, "testdata"), nil
}

// GetUnexportedField is a "hacky" function which enables to read private fields from a struct.
func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
