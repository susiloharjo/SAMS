package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Check if the stored hash matches "user.1001"
	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"), []byte("user.1001"))
	fmt.Println("Password 'user.1001' check result:", err)
	
	// Check if the stored hash matches "password123"
	err = bcrypt.CompareHashAndPassword([]byte("$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"), []byte("password123"))
	fmt.Println("Password 'password123' check result:", err)
}
