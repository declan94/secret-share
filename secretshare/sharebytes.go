package secretshare

import (
	"bytes"
	"crypto/md5"
	"errors"

	"github.com/codahale/sss"
)

const shareBytesOverhead = md5.Size + 1

// ShareBytes create secret-share parts for bytes array. The part id and md5 hash of the original data will be append after the data before creating sharing parts. Use RecoverBytes to recover and verify origin data from sharing parts
// data: data to create shares
// n: count of sharing parts
// k: least count of sharing parts to recover origin data
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

// RecoverBytes recover bytes array from sharing parts
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
