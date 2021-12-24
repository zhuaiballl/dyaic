package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  bspatch -loc LOCATION - Run bsdiff for files in LOCATION")
	fmt.Println("  commit -loc LOCATION - Commit changes in LOCATION")
	fmt.Println("  diff -loc LOCATION - Show changes in LOCATION")
	fmt.Println("  patch -loc LOCATION - Generate patch file for files in LOCATION")
	fmt.Println("  print -loc LOCATION - Show files in LOCATION")
	fmt.Println("  watch -loc LOCATION - Start watching LOCATION")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()
	bscommtCmd := flag.NewFlagSet("bscommit", flag.ExitOnError)
	bspatchCmd := flag.NewFlagSet("bspatch", flag.ExitOnError)
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	patchCmd := flag.NewFlagSet("patch", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)

	bscommtLocation := bscommtCmd.String("loc", "", "location to be committed")
	bspatchLocation := bspatchCmd.String("loc", "", "location of files we calc bspatch for")
	commitLocation := commitCmd.String("loc", "", "location to be committed")
	diffLocation := diffCmd.String("loc", "", "location where changes should be showed")
	patchLocation := patchCmd.String("loc", "", "location of files we calc patch for")
	printLocation := printCmd.String("loc", "", "location to be showed")
	watchLocation := watchCmd.String("loc", "", "location to be watched")

	switch os.Args[1] {
	case "bscommit":
		err := bscommtCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "bspatch":
		err := bspatchCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "commit":
		err := commitCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "diff":
		err := diffCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "patch":
		err := patchCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "watch":
		err := watchCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	}

	if bscommtCmd.Parsed() {
		cli.commit(*bscommtLocation, true)
	}

	if bspatchCmd.Parsed() {
		cli.patch(*bspatchLocation, true)
	}

	if commitCmd.Parsed() {
		cli.commit(*commitLocation, false)
	}

	if diffCmd.Parsed() {
		cli.printDiff(*diffLocation)
	}

	if patchCmd.Parsed() {
		cli.patch(*patchLocation, false)
	}

	if printCmd.Parsed() {
		cli.printFolder(*printLocation)
	}

	if watchCmd.Parsed() {
		cli.watch(*watchLocation)
	}
}
