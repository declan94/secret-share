package secretshare

import (
	"fmt"
	"io"
	"os"
)

// ShareFile create sharing parts of file
//  src: path to source directory
//  dsts: paths to out sharing parts
//  k: least count of sharing parts to recover origin data
func ShareFile(src string, dsts []string, k byte) error {
	fmt.Printf("Create sharing parts for [%s] \n", src)
	header := newShareFileHeader()
	header.knum = k
	hdBytes := header.searialize()
	// Open src file, check mode
	fd, err := os.Open(src)
	defer fd.Close()
	if err != nil {
		return fmt.Errorf("Open file [%s] faild: %v", src, err)
	}
	finfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat file [%s] failed: %v", src, err)
	}
	mode := finfo.Mode()
	// must have write and read permission
	mode |= (6 << 6)
	blkCnt := (finfo.Size()-1)/int64(header.blockSize) + 1
	pfds := make([]*os.File, len(dsts))
	for i, dst := range dsts {
		pfds[i], err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		defer pfds[i].Close()
		if err != nil {
			return fmt.Errorf("Open part file [%s] failed: %v", dst, err)
		}
		_, err = pfds[i].Write(hdBytes)
		if err != nil {
			return fmt.Errorf("Write part file [%s] header failed: %v", dst, err)
		}
	}
	// Read blocks and write parts
	block := make([]byte, header.blockSize)
	blkNo := 0
	var rdErr error
	fmt.Print("  0%")
	defer fmt.Print("\x08\x08\x08\x08")
	for rdErr == nil {
		fmt.Printf("\x08\x08\x08\x08%3d%%", int64(blkNo*100)/blkCnt)
		var n int
		block = block[:header.blockSize]
		n, rdErr = fd.Read(block)
		if rdErr != nil && rdErr != io.EOF {
			return fmt.Errorf("Read file [%s] failed: %v", src, rdErr)
		} else if n > 0 {
			block = block[:n]
			parts, err := ShareBytes(block, byte(len(dsts)), k)
			if err != nil {
				return fmt.Errorf("Create sharing parts for [%s] block#%d failed: %v", src, blkNo, err)
			}
			for i, pfd := range pfds {
				_, err = pfd.Write(parts[i])
				if err != nil {
					return fmt.Errorf("Write part file [%s] block#%d failed: %v", dsts[i], blkNo, err)
				}
			}
		}
		blkNo++
	}
	return nil
}

// RecoverFile recover file from sharing parts
//  srcs: paths of sharing parts
//  dst: path to output recovered file
func RecoverFile(dst string, srcs []string) error {
	fmt.Printf("Recover [%s] \n", dst)
	var err error
	var header shareFileHeader
	var tmpHeader shareFileHeader
	pfds := make([]*os.File, len(srcs))
	hdData := make([]byte, headerLen)
	for i, src := range srcs {
		pfds[i], err = os.Open(src)
		defer pfds[i].Close()
		if err != nil {
			return fmt.Errorf("Open part file [%s] failed: %v", src, err)
		}
		_, err = pfds[i].Read(hdData)
		if err != nil {
			return fmt.Errorf("Read part file [%s] header failed: %v", src, err)
		}
		tmpHeader.unsearialize(hdData)
		if i == 0 {
			header = tmpHeader
			if header.version != CurVersion {
				return fmt.Errorf("Part file [%s] version not matched. want: %d, get: %d", src, CurVersion, header.version)
			}
		} else if !header.equal(&tmpHeader) {
			return fmt.Errorf("Part file [%s] header not matched", src)
		}
	}
	if len(srcs) < int(header.knum) {
		return fmt.Errorf("Not enough parts, given: %d, need: %d", len(srcs), header.knum)
	}
	finfo, err := os.Stat(srcs[0])
	if err != nil {
		return fmt.Errorf("Stat part file [%s] failed: %v", srcs[0], err)
	}
	mode := finfo.Mode()
	blkCnt := (finfo.Size()-1)/int64(header.blockSize+shareBytesOverhead) + 1
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	defer fd.Close()
	if err != nil {
		return fmt.Errorf("Open dst file [%s] failed: %v", dst, err)
	}
	// Read blocks and Recover
	pblocks := make([][]byte, len(srcs))
	for i := range pblocks {
		pblocks[i] = make([]byte, header.blockSize+shareBytesOverhead)
	}
	blkNo := 0
	var rdErr error
	fmt.Print("  0%")
	defer fmt.Print("\x08\x08\x08\x08")
	for rdErr == nil {
		fmt.Printf("\x08\x08\x08\x08%3d%%", int64(blkNo*100)/blkCnt)
		var n0, n int
		for i, pfd := range pfds {
			pblocks[i] = pblocks[i][:header.blockSize+shareBytesOverhead]
			n, rdErr = pfd.Read(pblocks[i])
			if rdErr != nil && rdErr != io.EOF {
				return fmt.Errorf("Read part file [%s] failed: %v", srcs[i], err)
			}
			if i == 0 {
				n0 = n
			} else if n0 != n {
				return fmt.Errorf("Part file [%s] length not matched", srcs[i])
			}
			pblocks[i] = pblocks[i][:n]
		}
		if n0 > 0 {
			oriBlock, err := RecoverBytes(pblocks)
			if err != nil {
				return fmt.Errorf("Recover dst file [%s] block#%d failed: %v", dst, blkNo, err)
			}
			_, err = fd.Write(oriBlock)
			if err != nil {
				return fmt.Errorf("Write dst file [%s] block#%d failed: %v", dst, blkNo, err)
			}
		}
		blkNo++
	}
	return nil
}
