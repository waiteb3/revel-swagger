package common

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/revel/revel"
)

// The module path

// Get the directory of the Swagger-UI assets
var ModulePath string
var SwaggerAssetsDir string

func init() {
	_, ModulePath, _, _ = runtime.Caller(1)
	ModulePath = path.Dir(ModulePath)
	SwaggerAssetsDir = filepath.Join(ModulePath, "swagger-ui-master", "dist")
	revel.INFO.Println(ModulePath, SwaggerAssetsDir)
}

var warnedOnce = false

func UnzipSwaggerAssets() {
	dir := ModulePath
	fstat, err := os.Lstat(filepath.Join(dir, "swagger-ui-master"))
	if err == nil && fstat.IsDir() {
		if !warnedOnce {
			warnedOnce = true
			revel.INFO.Println("Swagger-UI Assets may have already been unzipped at " + dir +
				" and will be left alone.")
		}
		return
	}

	zipr, err := zip.OpenReader(filepath.Join(dir, "master.zip"))
	if err != nil {
		revel.ERROR.Println("Unable to find Swagger-UI Assets master.zip.", filepath.Join(dir, "master.zip"), err)
		return
	}

	for _, file := range zipr.File {
		if file.FileInfo().IsDir() {
			err := os.Mkdir(filepath.Join(dir, file.Name), file.Mode())
			if err != nil {
				revel.ERROR.Println("Error making directory for unzipping Swagger-UI Assets.", err)
				return
			}
			continue
		}

		reader, err := file.Open()
		if err != nil {
			reader.Close()
			revel.ERROR.Println("Unable to open file inside Swagger-UI Assets zip.", err)
			return
		}

		newfile, err := os.OpenFile(filepath.Join(dir, file.Name), os.O_WRONLY|os.O_CREATE, file.Mode())
		if err != nil {
			revel.ERROR.Println("Unable to make file for unzipping Swagger-UI Assets.", err)
			return
		}

		_, err = io.Copy(newfile, reader)
		if err != nil {
			revel.ERROR.Println("Unable to copy contents into file while unzipping Swagger-UI Assets.", err)
			return
		}

		if err := reader.Close(); err != nil {
			revel.ERROR.Println("Unable to close reader while unzipping Swagger-UI Assets.", err)
			return
		}
	}
}
