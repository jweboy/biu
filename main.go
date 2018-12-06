package main

import (
	"log"
	"os"

	"github.com/jweboy/biu/reader"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "biu"
	app.Usage = "fight the loneliness!"
	app.UsageText = "biu [global options] command [command options] [arguments...]"
	app.Author = "jweboy"
	app.Email = "jl940630@gmail.com"
	// app.Action = func(c *cli.Context) error {
	// 	fmt.Println("Hello friend!")
	// 	return nil
	// }
	app.Version = `
    ___       ___       ___       ___   
   /\__\     /\  \     /\__\     /\  \  
  |::L__L   _\:\  \   |::L__L   _\:\  \ 
  |:::\__\ /\/::\__\  |:::\__\ /\/::\__\
  /:;;/__/ \::/\/__/  /:;;/__/ \::/\/__/
  \/__/     \:\__\    \/__/     \:\__\  
             \/__/               \/__/   v0.0.1
`
	app.Commands = reader.Commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
