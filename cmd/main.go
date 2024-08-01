package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/teleology-io/foundation-cli/api"
	"github.com/urfave/cli/v2"
)

var client api.ApiClient

func main() {
	var variableName string
	var uniqueID string

	app := &cli.App{
		Name:     "foundation",
		Usage:    "CLI for the Foundation API",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Aliases: []string{"key"},
				Usage:   "The api-key to make requests with",
				EnvVars: []string{"FOUNDATION_API_KEY"},
				Action: func(cCtx *cli.Context, flag string) error {
					client = api.Create(flag)

					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "environment",
				Aliases: []string{"env"},
				Usage:   "Gets the environment",
				Action: func(cCtx *cli.Context) error {
					env, err := client.GetEnvironment()
					if err != nil {
						fmt.Println("environment command failed:", err.Error())
						return nil
					}

					for k, v := range env {
						fmt.Printf("%s=%v\n", k, v)
					}

					return nil
				},
			},
			{
				Name:    "configuration",
				Aliases: []string{"config"},
				Usage:   "Gets the configuration",
				Action: func(cCtx *cli.Context) error {
					config, err := client.GetConfiguration()
					if err != nil {
						fmt.Println("configuration command failed:", err.Error())
						return nil
					}

					fmt.Printf("%s\n", config)

					return nil
				},
			},
			{
				Name:  "variable",
				Usage: "Gets a specific variable",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Aliases:     []string{"n"},
						Usage:       "The name of the variable",
						Destination: &variableName,
					},
					&cli.StringFlag{
						Name:        "uid",
						Usage:       "The unique identifier of a user/entity",
						Destination: &uniqueID,
					},
				},
				Action: func(cCtx *cli.Context) error {
					result, err := client.GetVariable(variableName, uniqueID)
					if err != nil {
						fmt.Println("variable command failed:", err.Error())
						return nil
					}

					fmt.Printf("%v\n", result)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
