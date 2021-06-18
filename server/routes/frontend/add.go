// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package frontend

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddToRouter adds a static file serve for all files configured in the frontend dir.
func AddToRouter(frontendDir string, router *gin.Engine) error {
	files, err := ioutil.ReadDir(frontendDir)
	if err != nil {
		return err
	}

	knownFiles := make([]string, len(files))
	for i, file := range files {
		routePath := path.Join("/", file.Name())
		knownFiles[i] = routePath
		if file.IsDir() {
			router.Static(routePath, filepath.Join(frontendDir, file.Name()))
			continue
		}

		router.GET(routePath, func(filename string) func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				ctx.File(filepath.Join(frontendDir, filename))
			}
		}(file.Name()))
	}
	// register index file
	router.GET("/", func(ctx *gin.Context) {
		ctx.File(filepath.Join(frontendDir, "index.html"))
	})

	// redirect all unknown path to the index file
	router.Use(func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.Next()
			return
		}
		if isKnownFile(knownFiles, ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		ctx.Next()
		if ctx.Writer.Status() == http.StatusNotFound {
			ctx.File(filepath.Join(frontendDir, "index.html"))
		}
	})
	return nil
}

func isKnownFile(knownFiles []string, prefix string) bool {
	for _, file := range knownFiles {
		if strings.HasPrefix(file, prefix) {
			return true
		}
	}
	return false
}
