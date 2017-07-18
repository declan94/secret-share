package secretshare

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"os"

	"github.com/declan94/secret-share/internal/tlog"
)

// ShareDirectory create sharing parts for a directory
//  src: path to source directory
//  dsts: paths to out sharing parts
//  k: least count of sharing parts to recover origin data
func ShareDirectory(src string, dsts []string, k byte) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("Read dir [%s] failed: %v", src, err)
	}
	for _, finfo := range files {
		subDsts := make([]string, len(dsts))
		for i, d := range dsts {
			subDsts[i] = filepath.Join(d, finfo.Name())
		}
		subSrc := filepath.Join(src, finfo.Name())
		if finfo.IsDir() {
			for _, d := range subDsts {
				os.Mkdir(d, finfo.Mode()|(6<<6))
			}
			err := ShareDirectory(subSrc, subDsts, k)
			if err != nil {
				tlog.Warn.Printf("Create share for subdir [%s] failed: %v", subSrc, err)
			}
		} else {
			err := ShareFile(subSrc, subDsts, k)
			if err != nil {
				tlog.Warn.Printf("Create share for subfile [%s] failed: %v", subSrc, err)
			}
		}
	}
	return nil
}

// RecoverDirectory recover directory from sharing parts
//  srcs: paths of sharing parts
//  dst: path to output recovered directory
func RecoverDirectory(dst string, srcs []string) error {
	src := srcs[0]
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("Read dir [%s] failed: %v", src, err)
	}
	for _, finfo := range files {
		subSrcs := make([]string, len(srcs))
		for i, d := range srcs {
			subSrcs[i] = filepath.Join(d, finfo.Name())
		}
		subDst := filepath.Join(src, finfo.Name())
		if finfo.IsDir() {
			os.Mkdir(subDst, finfo.Mode())
			err := RecoverDirectory(subDst, subSrcs)
			if err != nil {
				tlog.Warn.Printf("Recover subdir [%s] failed: %v", subDst, err)
			}
		} else {
			err := RecoverFile(subDst, subSrcs)
			if err != nil {
				tlog.Warn.Printf("Recover subfile [%s] failed: %v", subDst, err)
			}
		}
	}
	return nil
}
