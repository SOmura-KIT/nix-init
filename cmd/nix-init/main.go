package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "nix-init"
	app.Usage = "Generate nix-shell init script"
	app.Commands = []*cli.Command{
		{
			Name:   "gen",
			Action: generateNixShell,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"o"},
					Usage:   "Place the output into `FILE`.",
					Value:   "shell.nix",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "Set the name of the nix-shell",
					Value:   "Template",
				},
				&cli.BoolFlag{
					Name:    "envrc",
					Aliases: []string{"e"},
					Usage:   "Generate .envrc file",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Overwrite existing files",
				},
				// TODO: implement this
				&cli.BoolFlag{
					Name:    "interactive",
					Aliases: []string{"i"},
					Usage:   "Ask for confirmation before overwriting",
				},
				&cli.BoolFlag{
					Name:    "pretend",
					Aliases: []string{"p"},
					Usage:   "Print the generated script to stdout",
				},
			},
		},
		{
			Name:   "list",
			Action: listTemplates,
		},
		{
			Name:   "config-file",
			Action: configFile,
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to the configuration file",
			Value:   fmt.Sprintf("%s/.config/nix-init/config.json", os.Getenv("HOME")),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// generate shell.nix file
func generateNixShell(c *cli.Context) error {
	templates, err := loadTemplates(c.String("config"))
	if err != nil {
		return err
	}

	// filter the templates to enable
	if c.NArg() == 0 {
		return fmt.Errorf("no template specified")
	}
	enables := []template{}
	for _, arg := range c.Args().Slice() {
		isFound := false
		for _, t := range templates {
			if t.Key == arg {
				enables = append(enables, t)
				isFound = true
				break
			}
		}
		if !isFound {
			return fmt.Errorf("template %s not found", arg)
		}
	}

	text := makeText(enables, c.String("name"))

	if c.Bool("pretend") {
		fmt.Println(text)
		return nil
	}

	// write the nix-shell script to the file
	path := c.String("path")
	// check if the file already exists
	if _, err := os.Stat(path); err == nil {
		if !c.Bool("force") {
			return fmt.Errorf("file %s already exists", path)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	// generate the .envrc file
	if c.Bool("envrc") {
		if _, err := os.Stat(".envrc"); err == nil {
			fmt.Println(".envrc already exists")
			fmt.Println("Do $ echo \"use_nix\" >> .envrc to enable nix-shell")
		}

		envrc, err := os.Create(".envrc")
		if err != nil {
			return err
		}
		defer envrc.Close()

		_, err = envrc.WriteString("use_nix")
		if err != nil {
			return err
		}
	}

	return nil
}

// generate the nix-shell script from the templates
func makeText(templates []template, name string) string {
	lines := []string{
		"{ pkgs ? import <nixpkgs> {} }:",
		"",
		"pkgs.mkShell {",
		"  name = \"" + name + "\";",
		"  buildInputs = with pkgs; [",
	}

	for _, t := range templates {
		lines = append(lines, "    # "+t.Key)
		for _, p := range t.Pkgs {
			lines = append(lines, "    "+p)
		}
	}

	lines = append(lines, "  ];")
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

// load the configuration file and list the available templates
func listTemplates(c *cli.Context) error {
	templates, err := loadTemplates(c.String("config"))
	if err != nil {
		return err
	}

	for _, t := range templates {
		fmt.Printf("%s:\n", t.Key)
		for _, p := range t.Pkgs {
			fmt.Printf("  - %s\n", p)
		}
	}

	return nil
}

type template struct {
	Key  string   `json:"key"`
	Pkgs []string `json:"pkgs"`
}

// load the configuration file
func loadTemplates(pathToConfig string) ([]template, error) {
	// open the configuration file
	config, err := os.Open(pathToConfig)
	if err != nil {
		return nil, err
	}
	defer config.Close()

	// decode the configuration file
	var templates []template
	decoder := json.NewDecoder(config)
	err = decoder.Decode(&templates)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// print the path to the configuration file
func configFile(c *cli.Context) error {
	fmt.Println("~/.config/nix-init/config.json")
	return nil
}
