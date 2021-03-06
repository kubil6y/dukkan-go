package data

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgconn"
)

func IsDuplicateRecord(err error) bool {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code == "23505"
		}
	}
	return false
}

// SLUGIFY SETUP BEGIN //////////////////////////////
var randSeed = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSeed.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSeed.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func Slugify(s string, n int) string {
	result := slug.Make(s + "-" + RandStringBytesMaskImprSrc(n))
	return result
}

// SLUGIFY SETUP END //////////////////////////////

func In(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
