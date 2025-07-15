package jwt_util

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// LoadPublicKey loads the public key from the specified path in the environment variable.
// It returns the parsed RSA public key or an error if the file cannot be read or parsed.
func LoadPublicKey() (*rsa.PublicKey, error) {
	jwtPublicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if jwtPublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH environment variable is not set")
	}

	keyData, err := os.ReadFile(jwtPublicKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// LoadPrivateKey loads the private key from the specified path in the environment variable.
// It returns the parsed RSA private key or an error if the file cannot be read or parsed.
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	jwtPrivateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if jwtPrivateKeyPath == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY_PATH environment variable is not set")
	}

	keyData, err := os.ReadFile(jwtPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}
