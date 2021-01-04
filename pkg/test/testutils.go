package test

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/config"
	"path"
	"reflect"
	"unsafe"
)

func GetTestdataFolderPath() (string, error) {
	rootPath, err := config.GetRootProjectPath()
	if err != nil {
		return "", fmt.Errorf("unable to get root path: %s", err)
	}
	return path.Join(rootPath, "testdata"), nil
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
