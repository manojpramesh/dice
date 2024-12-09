package geo

import "strings"

const geoAlphabet = "0123456789bcdefghjkmnpqrstuvwxyz"

func EncodeHash(hash uint64) string {
	var sb strings.Builder
	for i := 0; i < 11; i++ {
		var idx int
		if i == 10 {
			idx = 0
		} else {
			idx = int((hash >> (52 - ((i + 1) * 5))) & 0x1f)
		}
		sb.WriteByte(geoAlphabet[idx])
	}
	return sb.String()
}
