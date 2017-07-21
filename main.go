package main

import (
	"fmt"
	"os"

	"github.com/declan94/secret-share/internal/exitcode"
	"github.com/declan94/secret-share/internal/tlog"
	"github.com/declan94/secret-share/secretshare"
)

// ReleaseVersion Release version string
const ReleaseVersion = "v0.1"

func main() {
	args := ParseArgs()
	if args.Version {
		printVersion()
		return
	}
	var err error
	if args.Recover {
		if args.Directory {
			err = secretshare.RecoverDirectory(args.Src, args.Parts)
		} else {
			err = secretshare.RecoverFile(args.Src, args.Parts)
		}

	} else {
		if args.Directory {
			err = secretshare.ShareDirectory(args.Src, args.Parts, byte(args.KNum))
		} else {
			err = secretshare.ShareFile(args.Src, args.Parts, byte(args.KNum))
		}
	}
	if err != nil {
		tlog.Fatal.Println(err)
		os.Exit(exitcode.Execution)
	}
}

func printVersion() {
	fmt.Printf("secret-share %s; share parts format version: %d.\n", ReleaseVersion, secretshare.CurVersion)
}
