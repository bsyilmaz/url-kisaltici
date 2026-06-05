// Package shortcode rastgele ve tahmin edilmesi zor kisa kodlar uretir.
package shortcode

import (
	"crypto/rand"
	"math/big"
)

// alphabet kodlarda kullanilan karakter kumesi. Karisik olmamasi icin
// buyuk/kucuk harf ve rakamlar bir arada tutuluyor.
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const defaultLength = 6

// Generate verilen uzunlukta rastgele bir kod uretir. Uzunluk sifir veya
// negatif verilirse varsayilan uzunluk kullanilir. Tahmin edilebilirligi
// dusurmek icin crypto/rand uzerinden uretiyoruz.
func Generate(length int) (string, error) {
	if length <= 0 {
		length = defaultLength
	}

	result := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))
	for i := range result {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[idx.Int64()]
	}
	return string(result), nil
}
