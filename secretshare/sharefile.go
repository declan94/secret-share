package secretshare

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ShareFile create share parts of file with path "src"
func ShareFile(src string, dsts []string, k byte) error {
	fmt.Printf("Create sharing parts for [%s] \n", src)
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("Read file [%s] faild: %v", src, err)
	}
	finfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat file [%s] failed: %v", src, err)
	}
	mode := finfo.Mode()
	// must have write and read permission
	mode |= (6 << 6)
	parts, err := ShareBytes(data, byte(len(dsts)), k)
	if err != nil {
		return fmt.Errorf("Create sharing parts failed: %v", err)
	}
	for i, dst := range dsts {
		err = ioutil.WriteFile(dst, parts[i], mode)
		if err != nil {
			return fmt.Errorf("Write to part file [%s] failed: %v", dsts[i], err)
		}
	}
	return nil
}

// RecoverFile recover file by sharing parts
func RecoverFile(dst string, srcs []string) error {
	fmt.Printf("Recover [%s] \n", dst)
	parts := make([][]byte, len(srcs))
	for i, src := range srcs {
		data, err := ioutil.ReadFile(src)
		if err != nil {
			return fmt.Errorf("Read part file [%s] failed: %v", src, err)
		}
		parts[i] = data
	}
	finfo, err := os.Stat(srcs[0])
	if err != nil {
		return fmt.Errorf("Stat part file [%s] failed: %v", srcs[0], err)
	}
	mode := finfo.Mode()
	origin, err := RecoverBytes(parts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, origin, mode)
	if err != nil {
		return fmt.Errorf("Write to dst file [%s] failed: %v", dst, err)
	}
	return nil
}
