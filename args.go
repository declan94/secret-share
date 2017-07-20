package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"path/filepath"

	"github.com/declan94/secret-share/internal/tlog"
)

var flagSet *flag.FlagSet

// CliArgs contains cli args value
type CliArgs struct {
	Src       string
	Parts     []string
	Recover   bool
	Directory bool
	KNum      int
	Version   bool
}

func printMyFlagSet(avoid map[string]bool) {
	flagSet.VisitAll(func(f *flag.Flag) {
		if avoid[f.Name] {
			return
		}
		s := fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.
		_, usage := flag.UnquoteUsage(f)
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		s += "\n    \t"
		s += strings.Replace(usage, "\n", "\n    \t", -1)
		fmt.Println(s)
	})
}

func usage() {
	fmt.Printf("Usage: %s [options] SRC/DST PART1 PART2 PART3 ...\n", path.Base(os.Args[0]))
	fmt.Printf("		There must be at least two sharing parts.")
	fmt.Printf("\noptions:\n")
	printMyFlagSet(map[string]bool{"debug": true})
	os.Exit(1)
}

// ParseArgs parse args from cli args
func ParseArgs() (args CliArgs) {

	var debug bool
	flagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.BoolVar(&args.Recover, "r", false, "Recover mode, recover dst file(directory) from sharing parts.")
	flagSet.IntVar(&args.KNum, "k", 0, "Sufficent count for recovering. Default all sharing parts are needed.")
	flagSet.BoolVar(&args.Version, "v", false, "Show version info.")
	flagSet.BoolVar(&debug, "debug", false, "Log out debug info.")

	flagSet.Usage = usage
	flagSet.Parse(os.Args[1:])

	tlog.Debug.Enabled = debug

	if args.Version {
		return args
	}

	if flagSet.NArg() < 3 {
		usage()
	}

	args.Src = flagSet.Arg(0)
	args.Parts = flagSet.Args()[1:]

	info, err := os.Stat(args.Src)
	if args.Recover {
		if err == nil {
			tlog.Fatal.Printf("[%s] Already exists. Change another place to recover.\n", args.Src)
			os.Exit(1)
		}
		if !os.IsNotExist(err) {
			tlog.Fatal.Printf("Stat error: %v\n", err)
			os.Exit(1)
		}
		infos := make([]os.FileInfo, len(args.Parts))
		for i, part := range args.Parts {
			info, err := os.Stat(part)
			if err != nil {
				tlog.Fatal.Printf("stat [%s] failed: %v", part, err)
				os.Exit(1)
			}
			if info.IsDir() {
				args.Directory = true
			}
			infos[i] = info
		}
		for _, info := range infos {
			if info.IsDir() != args.Directory {
				tlog.Fatal.Printf("Mixed file and directory parts")
				os.Exit(1)
			}
		}
	} else {
		if len(args.Parts) > 255 || len(args.Parts) < 2 {
			tlog.Fatal.Printf("the number of sharing parts should satisfy 2 <= n <= 255.")
			os.Exit(1)
		}
		if args.KNum == 0 {
			args.KNum = len(args.Parts)
		} else if args.KNum < 2 || args.KNum > len(args.Parts) {
			tlog.Fatal.Printf("k should satisfy 2 <= k <= n. (n is the count of total sharing parts)")
			os.Exit(1)
		}
		if err != nil {
			tlog.Fatal.Printf("Stat source file failed: %v", err)
			os.Exit(1)
		}
		if info.IsDir() {
			args.Directory = true
			for _, d := range args.Parts {
				os.MkdirAll(d, os.FileMode(0774))
				err := checkDirEmpty(d)
				if err != nil {
					tlog.Fatal.Println(err)
					os.Exit(1)
				}
			}
		} else {
			for _, p := range args.Parts {
				os.MkdirAll(filepath.Dir(p), os.FileMode(0774))
			}
		}
	}
	return args
}

// checkDirEmpty - check if "dir" exists and is an empty directory
func checkDirEmpty(dir string) error {
	err := checkDir(dir)
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}
	return fmt.Errorf("directory %s not empty", dir)
}

// checkDir - check if "dir" exists and is a directory
func checkDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	return nil
}
