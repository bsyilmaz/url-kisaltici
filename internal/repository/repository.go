package repository

import (
	"context"
	"errors"
)

// ErrNotFound istenen kod veritabaninda bulunamadiginda doner.
var ErrNotFound = errors.New("kayit bulunamadi")

// Repository kisaltma kayitlarinin saklandigi katmani temsil eder.
// Arkasinda Redis olabilir, baska bir sey olabilir; ust katmanlar bunu bilmez.
type Repository interface {
	Save(ctx context.Context, code, originalURL string) error
	Find(ctx context.Context, code string) (string, error)
	Exists(ctx context.Context, code string) (bool, error)
}
