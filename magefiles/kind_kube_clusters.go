package main

import (
	"encoding/json"
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"io/ioutil"
	"scripts/kind"
	"strings"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

type Kind mg.Namespace

//UploadImagesJSON downloads provided images in json file and uploads to kind cluster
func (Kind) UploadImagesJSON(registryListFile, clusterName string, deleteAfterUpload bool) error {
	f, err := ioutil.ReadFile(registryListFile)
	if err != nil {
		return err
	}
	var rl *kind.RegistryListWithExceptions
	err = json.Unmarshal(f, &rl)
	if err != nil {
		return err
	}
	return kind.UploadImagesJSON(rl, clusterName, deleteAfterUpload)
}

//UploadImages uploads images from a file containing image list (1 per line) to specified kind cluster
func (Kind) UploadImages(imageListFile, clusterName string, deleteAfterUpload bool) error {
	f, err := ioutil.ReadFile(imageListFile)
	if err != nil {
		return err
	}
	lines := strings.Split(string(f), "\n")
	if len(lines) < 1 {
		return fmt.Errorf("file does not container any images")
	}
	for _, l := range lines {
		err = kind.ValidateAndUploadImage(l, clusterName, deleteAfterUpload)
		if err != nil {
			return err
		}
	}
	return nil
}
