package hasher_test

import (
	"testing"

	"github.com/Josey34/goshop/pkg/hasher"
)

func TestBcryptHasher(t *testing.T) {
	h := hasher.NewBcryptHasher()

	t.Run("hash returns non-empty string", func(t *testing.T) {
		hashed, err := h.Hash("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if hashed == "" {
			t.Error("expected non-empty hash")
		}
		if hashed == "secret123" {
			t.Error("hash must differ from plaintext")
		}
	})

	t.Run("compare succeeds with correct password", func(t *testing.T) {
		hashed, _ := h.Hash("secret123")
		if err := h.Compare(hashed, "secret123"); err != nil {
			t.Errorf("unexpected err: %v", err)
		}
	})

	t.Run("compare fails with wrong password", func(t *testing.T) {
		hashed, _ := h.Hash("secret123")
		if err := h.Compare(hashed, "wrongpass"); err == nil {
			t.Error("expected error for wrong password")
		}
	})

	t.Run("two hashes of same password differ", func(t *testing.T) {
		h1, _ := h.Hash("secret123")
		h2, _ := h.Hash("secret123")
		if h1 == h2 {
			t.Error("expected unique hashes due to bcrypt salt")
		}
	})
}
