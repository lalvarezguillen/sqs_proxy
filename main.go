package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "The path to the json file that contains the configuration to run this proxy",
		},
	}
	app.Action = func(c *cli.Context) error {
		conf := c.String("config")
		if conf == "" {
			return fmt.Errorf("--config is a required argument")
		}
		fmt.Println("Configuration file:", conf)

		proxyConf, err := loadConfig(conf)
		if err != nil {
			return err
		}
		p := newProxy(proxyConf)
		p.Start()
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
