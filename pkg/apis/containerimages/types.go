package containerimages

import (
	"fmt"
	"github.com/heroku/docker-registry-client/registry"
	"log"
)

type ImageTags struct {
	Image string   `json:"image"`
	Tags  []string `json:"tags,omitempty"`
}

func (its *ImageTags) GetImageNames() []string {
	var items []string
	for _, it := range its.Tags {
		items = append(items, fmt.Sprintf("%s:%s", its.Image, it))
	}
	return items
}

func (its *ImageTags) emptyTags() bool {
	return its.Tags == nil || len(its.Tags) == 0
}

type ContainerImage struct {
	Org       string      `json:"org"`
	ImageTags []ImageTags `json:"image-tags"`
}

func (i *ContainerImage) GetImageNames() []string {
	var items []string
	for _, it := range i.ImageTags {
		for _, its := range it.GetImageNames() {
			items = append(items, fmt.Sprintf("%s/%s", i.Org, its))
		}
	}
	return items
}

func (i *ContainerImage) populateTags(r *registry.Registry) error {
	var err error
	for j, it := range i.ImageTags {
		if it.emptyTags() {
			i.ImageTags[j].Tags, err = r.Tags(fmt.Sprintf("%s/%s", i.Org, it.Image))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *ContainerImage) GetAllImageNames() []string {
	var images []string
	for _, it := range i.ImageTags {
		for _, it1 := range it.GetImageNames() {
			images = append(images, fmt.Sprintf("%s/%s", i.Org, it1))
		}
	}
	return images
}

type Registry struct {
	URL      string           `json:"url"`
	APIURL   string           `json:"api-url"`
	UserName string           `json:"user-name,omitempty"`
	Password string           `json:"password,omitempty"`
	Images   []ContainerImage `json:"images"`
}

func (r *Registry) populateEmptyTags() error {
	apiURL := r.APIURL
	if apiURL == "" {
		apiURL = r.URL
	}
	reg, err := registry.New(apiURL, r.UserName, r.Password)
	if err != nil {
		return err
	}
	for i, _ := range r.Images {
		err = r.Images[i].populateTags(reg)
		if err != nil {
			return fmt.Errorf("unable to populate tags %w", err)
		}
	}
	return nil
}

func (r *Registry) GetImages() ([]string, error) {
	var images []string
	err := r.populateEmptyTags()
	if err != nil {
		return images, err
	}
	log.Printf("%#v\n", r)
	for _, i := range r.Images {
		for _, j := range i.GetAllImageNames() {
			images = append(images, fmt.Sprintf("%s/%s", r.URL, j))
		}
	}
	return images, nil
}

type RegistryList struct {
	ContainerRegistries []Registry `json:"container-registries"`
}

func (rl *RegistryList) GetImages() ([]string, error) {
	var images []string
	for _, it := range rl.ContainerRegistries {
		i, err := it.GetImages()
		if err != nil {
			return images, err
		}
		images = append(images, i...)
	}
	return images, nil
}
