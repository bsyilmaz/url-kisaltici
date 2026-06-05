package shortcode

import (
	"strings"
	"testing"
)

func TestGenerateUzunluk(t *testing.T) {
	kod, err := Generate(8)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if len(kod) != 8 {
		t.Errorf("uzunluk 8 olmaliydi, %d geldi", len(kod))
	}
}

func TestGenerateVarsayilanUzunluk(t *testing.T) {
	// Sifir verilince varsayilan uzunluk (6) kullanilmali.
	kod, err := Generate(0)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if len(kod) != defaultLength {
		t.Errorf("varsayilan uzunluk %d olmaliydi, %d geldi", defaultLength, len(kod))
	}
}

func TestGenerateNegatifUzunluk(t *testing.T) {
	kod, err := Generate(-5)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if len(kod) != defaultLength {
		t.Errorf("negatif uzunlukta varsayilan kullanilmaliydi, %d geldi", len(kod))
	}
}

func TestGenerateSadeceGecerliKarakter(t *testing.T) {
	kod, err := Generate(30)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	for _, c := range kod {
		if !strings.ContainsRune(alphabet, c) {
			t.Errorf("alfabe disinda karakter uretildi: %q", c)
		}
	}
}

func TestGenerateBenzersizlik(t *testing.T) {
	// Ayni kodun art arda uretilme ihtimali pratikte sifir olmali.
	gorulen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		kod, err := Generate(7)
		if err != nil {
			t.Fatalf("beklenmeyen hata: %v", err)
		}
		if gorulen[kod] {
			t.Errorf("ayni kod iki kez uretildi: %s", kod)
		}
		gorulen[kod] = true
	}
}
