package scape

import (
	"github.com/dmp42/scape-go/docker"
	/*
	"github.com/docker/engine-api/types"
	*/
	"os"
	"path"
	"path/filepath"
)

func Infer(selector * docker.Selector) {
	if selector.Path == "" {
		var err error
		selector.Path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	// If no url, try to outsmart from path
	if selector.URL == "" {
		selector.URL = path.Base(path.Dir(path.Dir(selector.Path))) + "/" + path.Base(path.Dir(selector.Path)) + "/" + path.Base(selector.Path)
	}

	if !filepath.IsAbs(selector.Path) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		selector.Path = path.Join(wd, selector.Path)
	}
}

// Prepare populates a selector with inferred values if some are missing, and return containers that matches the selector, if any
/*
func Prepare(selector docker.Selector) []types.Container {
	var containers []types.Container
	// If we have a name, try to find it
	if selector.Name != "" {
		containers = docker.Select(docker.Selector{Name: selector.Name}, false)
		if len(containers) > 0 {
			return containers
		}
	}

	// If no path, get pwd
	if selector.Path == "" {
		selector.Path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	// If no url, try to outsmart from path
	if selector.URL == "" {
		selector.URL = path.Base(path.Dir(path.Dir(selector.Path))) + "/" + path.Base(path.Dir(selector.Path)) + "/" + path.Base(selector.Path)
	}

	// Nothing found from name? Then lookup from the same project properties
	return docker.Select(docker.Selector{URL: selector.URL, Path: selector.Path}, false)
}
*/
