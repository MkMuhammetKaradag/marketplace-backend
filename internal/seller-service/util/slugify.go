// util/slugify.go
package util

import (
	"regexp"
	"strings"
)

// Slugify, metni URL dostu bir slug'a çevirir (küçük harf, boşlukları tire ile değiştirme vb.)
func Slugify(s string) string {
	s = strings.ToLower(s)
	// Türkçe karakterleri çevirme (Örn: ş->s, ı->i, ç->c, ğ->g, ü->u, ö->o)
	s = strings.ReplaceAll(s, "ş", "s")
	s = strings.ReplaceAll(s, "ı", "i")
	s = strings.ReplaceAll(s, "ç", "c")
	s = strings.ReplaceAll(s, "ğ", "g")
	s = strings.ReplaceAll(s, "ü", "u")
	s = strings.ReplaceAll(s, "ö", "o")

	// Alfabetik olmayan karakterleri tire ile değiştirme
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")

	// Baş ve sondaki tireleri temizleme
	s = strings.Trim(s, "-")

	return s
}
