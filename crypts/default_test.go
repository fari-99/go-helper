package crypts

import "testing"

func TestDecryptShortInputDoesNotPanic(t *testing.T) {
	for _, input := range []string{"", "a", "YWJj"} {
		if _, err := NewEncryptionBase().SetPassphrase("key").Decrypt([]byte(input)); err == nil {
			t.Errorf("Decrypt(%q) expected error", input)
		}
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	enc, err := NewEncryptionBase().SetPassphrase("key").Encrypt([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	dec, err := NewEncryptionBase().SetPassphrase("key").Decrypt(enc)
	if err != nil || string(dec) != "hello" {
		t.Fatalf("got %q, %v", dec, err)
	}
}
