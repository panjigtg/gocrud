package helper

import (
	"net/url"
	"regexp"
)

// regex: hanya huruf, angka, spasi, garis bawah, strip
var safeChars = regexp.MustCompile(`[^a-zA-Z0-9\s-_]`)

// CleanText membersihkan string dari karakter berbahaya
func CleanText(input string) string {
	unescaped, _ := url.PathUnescape(input)
	return safeChars.ReplaceAllString(unescaped, "")
}
