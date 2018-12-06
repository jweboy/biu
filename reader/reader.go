package reader

import (
	"github.com/jweboy/biu/commands/renren"
	"github.com/urfave/cli"
)

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "get",
			Usage:  "Get movies",
			Action: renren.GetUSMovie,
		},
	}
}
