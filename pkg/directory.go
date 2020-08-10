package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

func getWorkingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	return workingDir
}

var workingDir = getWorkingDir()
var assetDir = path.Join(workingDir, "asset")
var configDir = path.Join(workingDir, "config")
var dataDir = path.Join(workingDir, "data")
var scriptsDir = path.Join(workingDir, "scripts")

var DataPath = path.Join(dataDir, "data.json")
var ConfigPath = path.Join(configDir, "config.yaml")

var _ = assetDir
var _ = configDir
var _ = dataDir
var _ = scriptsDir
