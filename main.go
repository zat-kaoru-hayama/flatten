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
					installerRoot := args[1]
					flatDir := filepath.Join(installerRoot, "flat")
					msmDir := filepath.Join(installerRoot, "msm")

					matched := make(map[string]int)
					flatReferred := make(map[string]bool)

					filepath.Walk(approvedRoot, func(path string, f os.FileInfo, err error) error {
						if f.IsDir() {
							return nil
						}
						// test Flat
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
						flatReferred[name] = false
						return nil
					})
					filepath.Walk(msmDir, func(path string, f os.FileInfo, err error) error {
						if f.IsDir() {
							return nil
						}
						name := strings.ToUpper(f.Name())
						if _, ok := flatReferred[name]; !ok {
							fmt.Printf("[ERR] Not found in Flat: %s\n", path)
						} else {
							flatReferred[name] = true
						}
						return nil
					})
					for name, refer := range flatReferred {
						if !refer {
							fmt.Printf("[ERR] Not referred from MSM: %s\n", name)
						}
					}
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
