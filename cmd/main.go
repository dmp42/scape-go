package main

import (
	"github.com/codegangsta/cli"
	"github.com/dmp42/scape-go/docker"
	"github.com/dmp42/scape-go/scape"
	"os"
	//	"log"
	"fmt"
	"github.com/docker/engine-api/types"
)

func init() {

}

func main() {
	var debug bool
	var force bool

	app := cli.NewApp()

	app.Name = "Scape Go"
	app.Usage = "Use it, blame it!"
	app.Author = "dmp42"
	app.Email = "viapanda@gmail.com"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "language for the greeting",
			EnvVar:      "SCAPE_DEBUG",
			Destination: &debug,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "list",
			Aliases:     []string{"ls"},
			Usage:       "list",
			Description: "List all started scape projects",
			Action: func(c *cli.Context) {
				// Get all scape containers
				containers := docker.Select(docker.Selector{}, true)
				fmt.Println("______________________________________________________________________________________________________________________________")
				fmt.Printf("%-12s | %-20s | %-30s | %-50s\n", "ID", "Name", "Package", "Path")
				fmt.Println("______________________________________________________________________________________________________________________________")
				for _, c := range containers {
					fmt.Printf("%-12s | %-20s | %-30s | %-50s\n", c.ID[:12], c.Names[0][1:], c.Labels["com.dmp42.scape.url"], c.Labels["com.dmp42.scape.path"])
				}
			},
		},
		{
			Name:        "info",
			Aliases:     []string{},
			Usage:       "info [optional_name]",
			Description: "Returns info about (a) specific project(s) by name, defaulting to project(s) in pwd if no name is specified.",
			Action: func(c *cli.Context) {
				sel := docker.Selector{}
				if c.NArg() > 0 {
					sel.Name = c.Args()[0]
				}else{
					scape.Infer(&sel)
				}
				containers := docker.Select(sel, false)
				if len(containers) == 0 {
					fmt.Println("No matchin project found!")
					return
				}
				fmt.Println("______________________________________________________________________________________________________________________________")
				fmt.Printf("%-12s | %-20s | %-30s | %-50s\n", "ID", "Name", "Package", "Path")
				fmt.Println("______________________________________________________________________________________________________________________________")
				for _, c := range containers {
					fmt.Printf("%-12s | %-20s | %-30s | %-50s\n", c.ID[:12], c.Names[0][1:], c.Labels["com.dmp42.scape.url"], c.Labels["com.dmp42.scape.path"])
				}
			},
		},
		{
			Name:        "destroy",
			Aliases:     []string{},
			Usage:       "",
			Description: "Removes an existing scape project, defaulting to pwd if no name is specified.",
			Action: func(c *cli.Context) {
				sel := docker.Selector{}
				if c.NArg() > 0 {
					sel.Name = c.Args()[0]
				}else{
					scape.Infer(&sel)
				}
				containers := docker.Select(sel, false)
				for _, c := range containers {
					docker.Stop(c)
					docker.Remove(c)
				}
			},
		},
		{
			Name:        "init",
			Aliases:     []string{},
			Usage:       "scape init",
			Description: "Initializes a scape project. If no argument is specified, will default to whatever is in pwd, infer the package url from the path, and give it a random name. You can specify any of these explicitely (path, package url, project name).",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "path,p"},
				cli.StringFlag{Name: "url,u"},
				//				cli.StringFlag{Name: "tags,t"},
				cli.BoolFlag{Name: "force,f", Destination: &force},
			},
			Action: func(c *cli.Context) {
				sel := docker.Selector{Name: ""}
				// Optionally get the name from the command line
				var containers []types.Container
				if c.NArg() > 0 {
					sel.Name = c.Args()[0]
					// Do we have a container with that name ?
					containers = docker.Select(sel, false)
					if len(containers) > 0 {
						if !force {
							fmt.Println("A container already exists by that name! Use -f or --force to stop & remove it before creating this.")
							return
						}
						// Force remove
						for _, c := range containers {
							docker.Stop(c)
							docker.Remove(c)
						}
					}
				}
				// Complement the selector with local info if any
				sel.Path = c.String("path")
				sel.URL = c.String("url")
				// Make sure the selector is populated with default info
				scape.Infer(&sel)

				containers = docker.Select(sel, false)
				if len(containers) > 0 {
					if !force {
						fmt.Println("A container already exists for that project! Use -f or --force to stop & remove it before creating this.")
						return
					}
					// Force remove
					for _, c := range containers {
						docker.Stop(c)
						docker.Remove(c)
					}
				}
				// Tags
				//				tags := c.String("tags")

				// Run the container with the given selector
				docker.Run(sel)
			},
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() < 2 {
			fmt.Print("Say something darn it!")
			return
		}
		sel := docker.Selector{}
		sel.Name = c.Args()[0]
		docker.Exec(sel, c.Args()[1:])
	}

	app.Run(os.Args)
}
