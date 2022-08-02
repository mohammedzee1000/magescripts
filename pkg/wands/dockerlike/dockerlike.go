package dockerlike

import (
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"os"
)

type DockerRunner struct {
	runner string
}

func New(runner string) *DockerRunner {
	return &DockerRunner{runner: runner}
}

func NewDefault() *DockerRunner {
	return New("docker")
}

func (dr *DockerRunner) Pull(image string) error {
	return sh.Run(dr.runner, "pull", image)
}

func (dr *DockerRunner) SaveImage(image string, filePath string) error {
	o, err := sh.Output(dr.runner, "save", image)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, []byte(o), os.ModePerm)
}

func (dr *DockerRunner) SaveImageTemp(image string) (string, error) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	err = tmp.Close()
	if err != nil {
		return "", err
	}
	err = os.Remove(path)
	if err != nil {
		return "", err
	}
	return path, dr.SaveImage(image, path)
}

func (dr *DockerRunner) RemoveImage(image string) error {
	return sh.Run(dr.runner, "rmi", "-f", image)
}
