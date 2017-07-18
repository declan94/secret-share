package main

import (
	"os"

	"github.com/declan94/secret-share/internal/tlog"
	"github.com/declan94/secret-share/secretshare"
)

func main() {
	args := ParseArgs()
	var err error
	if args.Recover {
		if args.Directory {
			err = secretshare.RecoverDirectory(args.Src, args.Parts)
		} else {
			err = secretshare.RecoverFile(args.Src, args.Parts)
		}

	} else {
		if args.Directory {
			err = secretshare.ShareDirectory(args.Src, args.Parts, byte(args.K))
		} else {
			err = secretshare.ShareFile(args.Src, args.Parts, byte(args.K))
		}
	}
	if err != nil {
		tlog.Fatal.Println(err)
		os.Exit(3)
	}
}
