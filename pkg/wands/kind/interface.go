package kind

type KindRunnerInterface interface {
	UploadImageArchive(path, clusterName string) error
}
