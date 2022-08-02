package kind

import (
	"github.com/magefile/mage/sh"
	"os"
)

type KindRunner struct {
	executor string
}

func New(executor string) *KindRunner {
	return &KindRunner{executor: executor}
}

func NewDefault() *KindRunner {
	return New("kind")
}

func (k *KindRunner) UploadImageArchive(path, clusterName string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return sh.Run(k.executor, "load", "image-archive", path, "--name", clusterName)
}
