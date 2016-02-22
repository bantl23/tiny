package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func main() {
	parse := true
	analyze := true
	code := true
	echo := false
	trace := false

	app := cli.NewApp()
	app.Name = "tiny"
	app.Usage = "tiny [flags] <filename>"
	app.Version = "1.0.0-alpha0"
	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:        "parse",
			Usage:       "Enable or disable code parsing",
			Destination: &parse,
		},
		cli.BoolTFlag{
			Name:        "analyze",
			Usage:       "Enable or disable code analysis",
			Destination: &analyze,
		},
		cli.BoolTFlag{
			Name:        "code",
			Usage:       "Enable or disable code generation",
			Destination: &code,
		},
		cli.BoolTFlag{
			Name:        "trace",
			Usage:       "True on code tracing",
			Destination: &trace,
		},
		cli.BoolTFlag{
			Name:        "echo",
			Usage:       "Print source code",
			Destination: &echo,
		},
	}
	app.Action = func(c *cli.Context) {

		if len(c.Args()) == 0 {
			fmt.Println("error: must supply filename to compile")
			os.Exit(1)
		}

		if analyze == false {
			code = false
		}
		if parse == false {
			analyze = false
			code = false
		}

		fmt.Printf("options: [parse=%t, analyze=%t, code=%t, echo=%t, trace=%t]\n",
			parse, analyze, code, echo, trace)
		for _, ifile := range c.Args() {
			if strings.HasSuffix(ifile, ".tny") == false {
				ifile = ifile + ".tny"
			}
			ofile := strings.TrimSuffix(ifile, ".tny") + ".tm"

			fmt.Println("compiling", ifile)
			fmt.Println("scan")
			fmt.Println("parse")
			fmt.Println("analyze")
			fmt.Println("codegen", ofile)
		}
	}
	app.Run(os.Args)
}
