package dockerlike

type DockerRunnerInterface interface {
	Pull(image string) error
	SaveImage(image, filepath string) error
	SaveImageTemp(image string) (string, error)
	RemoveImage(image string) error
}
