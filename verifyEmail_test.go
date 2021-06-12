package main

import (
	"testing"
)

const (
	SampleEmail = "sample@mail.com"
)

func TestGenerateAndVerify(t *testing.T) {
	t.Log("generating Request")
	key, err := GenerateEmailVerificationRequest(SampleEmail)
	if err != nil {
		t.Error("Generation Request failed")
	}

	result, err := VerifyEmail(key)

	if err != nil || !result {
		t.Error("Verify Failed")
	}
}
