package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	reset     = "\033[0m"
	blue      = "\033[34m"
	yellow    = "\033[33m"
	boldGreen = "\033[1;32m"
	clearLine = "\033[K"
	moveUp4   = "\033[4A"
)

func main() {
	secret := os.Getenv("TOTP_SECRET")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	dashes := strings.Repeat("-", 30)

	first := true

	for {
		now := time.Now()

		code, err := generateTOTP(secret, now, 6)
		if err != nil {
			panic(err)
		}

		if !first {
			fmt.Print(moveUp4)
		}

		n := now.Format("02/01/2006 15:04:05")

		fmt.Printf("\r%s%s%s%s\n"+
			"Time: %s%s%s%s\n"+
			"TOTP: %s%-2s%-2s%-2s%-2s%-2s%-2s%s%s\n"+
			"%s%s%s%s\n",
			blue,
			dashes,
			reset,
			clearLine,
			yellow,
			n,
			reset,
			clearLine,
			boldGreen,
			string(code[0]),
			string(code[1]),
			string(code[2]),
			string(code[3]),
			string(code[4]),
			string(code[5]),
			reset,
			clearLine,
			blue,
			dashes,
			reset,
			clearLine,
		)

		first = false

		<-ticker.C
	}
}

func generateTOTP(secret string, t time.Time, digits int) (string, error) {
	key, err := decodeKey(secret)
	if err != nil {
		return "", err
	}

	code := totp(key, t, digits)

	return fmt.Sprintf("%0*d", digits, code), nil
}

func decodeKey(secret string) ([]byte, error) {
	secret = strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, secret)

	secret = strings.ToUpper(secret)
	secret += strings.Repeat("=", -len(secret)&7)

	return base32.StdEncoding.DecodeString(secret)
}

func totp(key []byte, t time.Time, digits int) int {
	counter := uint64(t.UnixNano()) / 30e9

	return hotp(key, counter, digits)
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)

	if err := binary.Write(h, binary.BigEndian, counter); err != nil {
		panic(err)
	}

	sum := h.Sum(nil)

	offset := sum[len(sum)-1] & 0x0F

	value := binary.BigEndian.Uint32(sum[offset:]) & 0x7FFFFFFF

	divisor := uint32(1)

	for i := 0; i < digits && i < 8; i++ {
		divisor *= 10
	}

	return int(value % divisor)
}
