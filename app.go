package main

import (
	"fmt"
	"os"
	"time"
    "log"
    "net/http"
    "io/ioutil"

	"github.com/urfave/cli"
)

func handleApi(w http.ResponseWriter, r *http.Request){
    configStr := r.PostFormValue("config")
    insomniaFile, _, err := r.FormFile("insomnia")
    if err != nil {
        log.Fatal(err)
        fmt.Fprintf(w, "INVALID_FILE")
    }
    insomniaBody, err := ioutil.ReadAll(insomniaFile)
    if err != nil {
        log.Fatal(err)
    }

    swagger := &Swagger{}
    outputYaml := swagger.GenerateBuffer(string(insomniaBody), configStr, "yaml")
    fmt.Fprintf(w, outputYaml)
}

func main() {
	app := cli.NewApp()
	app.Name = "swaggonmia"
	app.Usage = "Insomnia to Swagger converter"
	app.Version = "1.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nick Wallace",
			Email: "nwallace@fyberstudios.com",
		},
	}
	app.Commands = []cli.Command{
        {
            Name:   "api",
            Usage:  "Run as API Server",
            Action: func(c *cli.Context) error {
                http.HandleFunc("/", handleApi)
                http.ListenAndServe(":8080", nil)
                return nil
            },
        },
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "config, c",
					Usage: "Load configuration from `FILE`",
				},
				&cli.StringFlag{
					Name:  "insomnia, i",
					Usage: "Insomnia JSON `FILE`",
				},
				&cli.StringFlag{
					Name:  "output, o",
					Value: "yaml",
					Usage: "Output json|yaml",
				},
			},
			Usage: "Generate Swagger documentation",
			Action: func(c *cli.Context) error {
				var insomniaFile = c.String("insomnia")
				var configFile = c.String("config")
				var outputFormat = c.String("output")

				if insomniaFile == "" || configFile == "" {
					cli.ShowCompletions(c)
				}

				if outputFormat == "" {
					outputFormat = "json"
				}

				swagger := &Swagger{}
				swagger.Generate(insomniaFile, configFile, outputFormat)

				return nil
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Wrong command %q !", command)
	}
	app.Run(os.Args)
}
