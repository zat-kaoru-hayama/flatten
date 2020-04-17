package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "status",
				Usage: "view approved and flat",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) < 2 {
						return errors.New("two few arguments")
					}
					approvedRoot := args[0]
					flatDir := args[1]

					matched := make(map[string]int)

					filepath.Walk(approvedRoot, func(path string, f os.FileInfo, err error) error {
						if f.IsDir() {
							return nil
						}
						basename := filepath.Base(path)
						stubpath := filepath.Join(flatDir, basename)
						_, err = os.Stat(stubpath)
						if err != nil {
							if os.IsNotExist(err) {
								fmt.Printf("[ERR] Not found in FLAT: %s\n", basename)
							} else {
								fmt.Fprintln(os.Stderr, err.Error())

							}
						} else {
							matched[strings.ToUpper(basename)]++
						}
						return nil
					})
					filepath.Walk(flatDir, func(path string, f os.FileInfo, err error) error {
						if f.IsDir() {
							return nil
						}
						name := strings.ToUpper(f.Name())
						if count, ok := matched[name]; !ok || count <= 0 {
							fmt.Printf("[ERR] Not found in Approved: %s\n", name)
						} else {
							matched[name]--
						}
						return nil
					})
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
