package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gitlab.com/bantl23/python/gen"
	"gitlab.com/bantl23/python/symtbl"
	"gitlab.com/bantl23/python/syntree"
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
		for _, ifilename := range c.Args() {
			if strings.HasSuffix(ifilename, ".tny") == false {
				ifilename = ifilename + ".tny"
			}
			ofilename := strings.TrimSuffix(ifilename, ".tny") + ".tm"

			fmt.Println("compiling", ifilename)
			ifile, _ := os.Open(ifilename)
			fmt.Println("parsing")
			yyParse(NewLexer(ifile))
			syntree.Print(root, 0)
			table := make(symtbl.SymTbl)
			table.BuildTable(root)
			fmt.Println("=====")
			table.PrintTable()
			fmt.Println("analyze")
			table.CheckTable(root)
			fmt.Println("code generation")
			gen := new(gen.Gen)
			gen.Generate(root, &table, ofilename)
		}
	}
	app.Run(os.Args)
}
