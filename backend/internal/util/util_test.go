package util_test

import (
	"fmt"
	"mini-assets/internal/util"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hashedPassword, err := util.HashPassword("password")
	if err != nil {
		t.Fatalf("Error hashing password: %v\n", err)
	}

	fmt.Printf("HashPassword: %v\n", hashedPassword)
}
