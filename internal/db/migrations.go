package db

import (
	"path/filepath"
	"runtime"
	"strings"
)

func GetPathToMigrationsDir() string {
	//nolint:dogsled
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	// basePath should be
	// /go/src/github.com/musicmash/musicmash/internal/db

	slittedPath := strings.Split(basePath, "/")
	slittedPath = slittedPath[:len(slittedPath)-2]
	// slittedPath should be
	// /go/src/github.com/musicmash/musicmash/internal

	slittedPath = append(slittedPath, "migrations")
	return strings.Join(slittedPath, "/")
}
