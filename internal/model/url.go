package model

// URL kisaltilmis bir kod ile asil adres arasindaki eslesmeyi temsil eder.
type URL struct {
	Code        string `json:"code"`
	OriginalURL string `json:"original_url"`
}
