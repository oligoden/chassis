package chassis

import (
	"math/rand"
)

const numLetterBytes = "0123456789"
const numalphaLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphaLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// RandString generates a random string
func RandString(n uint, s rand.Source) string {
	// solution from http://stackoverflow.com/a/31832326
	b := make([]byte, n)
	for i, cache, remain := int(n-1), s.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = s.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numalphaLetterBytes) {
			b[i] = numalphaLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandNumberString generates a random number string
func RandNumberString(n uint, s rand.Source) string {
	// solution from http://stackoverflow.com/a/31832326
	b := make([]byte, n)
	for i, cache, remain := int(n-1), s.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = s.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numLetterBytes) {
			b[i] = numLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandAlphaString generates a random alpha string
func RandAlphaString(n uint, s rand.Source) string {
	// solution from http://stackoverflow.com/a/31832326
	b := make([]byte, n)
	for i, cache, remain := int(n-1), s.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = s.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphaLetterBytes) {
			b[i] = alphaLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
