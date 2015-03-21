package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/motemen/go-cli"
)

func main() {
	cli.Use(&cli.Command{
		Name:  "foo",
		Short: "description in one line",
		Long: `foo [-v] <arg>

Description in paragraphs, starting with a usage line.
Blah blah blah`,
		Action: func(flags *flag.FlagSet, args []string) error {
			verbose := flags.Bool("v", false, "set verbosity")
			flags.Parse(args)

			args = flags.Args()
			if len(args) < 1 {
				return cli.ErrUsage
			}

			if *verbose {
				log.Println("showing foo...")
			}

			fmt.Println("foo", args[0])

			if *verbose {
				log.Println("succeeded.")
			}

			return nil
		},
	})
	cli.Run(os.Args[1:])
}
