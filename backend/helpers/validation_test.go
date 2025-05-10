package helpers

import "testing"

/** */
func TestValidPublicKey(t *testing.T) {
	input := "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e33d72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b"
	err := ValidateRawPublicKey(input)
	if err != nil {
		t.Errorf("Expected valid 128-char public key, got error: %v", err)
	}
}

func TestTooShortPublicKey(t *testing.T) {
	input := "aabbcc"
	err := ValidateRawPublicKey(input)
	if err == nil {
		t.Errorf("too short input..")
	}
}

func TestInvalidPublicKey(t *testing.T) {
	input := "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e32f72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b67h879jh89yui8ij"
	err := ValidateRawPublicKey(input)
	if err == nil {
		t.Errorf("too long input..")
	}
}
