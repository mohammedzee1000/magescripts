package kind

import (
	dockerparser "github.com/novln/docker-parser"
	"log"
	"os"
	"regexp"
	"scripts/pkg/apis/containerimages"
	"scripts/pkg/commands/dockerlike"
	"scripts/pkg/commands/kind"
)

type RegistryListWithExceptions struct {
	*containerimages.RegistryList
	ExceptionPatterns []string `json:"exception-patterns,omitempty"`
}

func UploadImagesJSON(rl *RegistryListWithExceptions, clusterName string, deleteAfterUpload bool) error {
	images, err := rl.GetImages()
	if err != nil {
		return err
	}
	for _, it := range images {
		log.Printf("working to upload image %s to kind cluster %s\n", it, clusterName)
		match := false
		for _, e := range rl.ExceptionPatterns {
			match, err = regexp.MatchString(e, it)
			if err != nil {
				return err
			}
			if match {
				break
			}
		}
		if match {
			log.Printf("skipping image %s due to exception match\n", it)
			continue
		}
		err = ValidateAndUploadImage(it, clusterName, deleteAfterUpload)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateAndUploadImage(image, clusterName string, deleteAfterUpload bool) error {
	var err error
	var path string
	_, err = dockerparser.Parse(image)
	if err != nil {
		return err
	}
	dr := dockerlike.NewDefault()
	err = dr.Pull(image)
	if err != nil {
		return err
	}
	path, err = dr.SaveImageTemp(image)
	if err != nil {
		return err
	}
	err = kind.NewDefault().UploadImageArchive(path, clusterName)
	if err != nil {
		return err
	}
	if deleteAfterUpload {
		err = dr.RemoveImage(image)
		if err != nil {
			return err
		}
	}
	if err = os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
