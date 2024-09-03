package wallet

import "github.com/mr-tron/base58"
import "log"

// Base58Encode Base58 encoding was designed to encode Bitcoin addresses. It has the following characteristics:
// its alphabet avoids similar looking letters.
// it does not use non-alphanumeric characters.
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))

	if err != nil {
		log.Panic(err)
	}

	return decode
}
