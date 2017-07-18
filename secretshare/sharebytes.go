package secretshare

import (
	"bytes"
	"crypto/md5"
	"errors"

	"github.com/codahale/sss"
)

// ShareBytesOverhead is the length overhead of sharing parts compared to the original data.
// i.e. the length of each sharing part = the length of original data + ShareBytesOverhead
const ShareBytesOverhead = md5.Size + 1

// ShareBytes create sharing parts for bytes array.
//
// First append md5 hash after `data` for later validation when recovering,
// then create `n` sharing parts, at least `k` parts are needed for recovering.
// Use RecoverBytes to recover and verify origin data from sharing parts
func ShareBytes(data []byte, n, k byte) ([][]byte, error) {
	hash := md5.Sum(data)
	hashedData := append(data, hash[:]...)
	splits, err := sss.Split(n, k, hashedData)
	if err != nil {
		return nil, err
	}
	results := make([][]byte, n)
	i := 0
	for id, share := range splits {
		results[i] = append(share, id)
		i++
	}
	return results, nil
}

// RecoverBytes recover bytes array from sharing parts.
// Validate using the pre-appended hash in ShareBytes
func RecoverBytes(shares [][]byte) ([]byte, error) {
	shareMap := make(map[byte][]byte)
	for _, share := range shares {
		id := share[len(share)-1]
		shareMap[id] = share[:len(share)-1]
	}
	hashedData := sss.Combine(shareMap)
	hash := hashedData[len(hashedData)-md5.Size:]
	data := hashedData[:len(hashedData)-md5.Size]
	curHash := md5.Sum(data)
	if !bytes.Equal(curHash[:], hash) {
		return nil, errors.New("Decrypted data check failed! Sharing part file broken or not sufficent count of parts")
	}
	return data, nil
}
