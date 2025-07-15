package jwt_util

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// GetInt64Claim retrieves an int64 claim from the JWT claims.
// It checks if the claim exists and is of type float64, then converts it to int64.
func GetInt64Claim(claims jwt.MapClaims, key string) (int64, error) {
	if val, ok := claims[key]; ok {
		if f, ok := val.(float64); ok {
			return int64(f), nil
		}
		return 0, fmt.Errorf("claim %s is not a float64", key)
	}
	return 0, fmt.Errorf("claim %s not found", key)
}

// GetStringClaim retrieves a string claim from the JWT claims.
// It checks if the claim exists and is of type string.
func GetStringSliceClaim(claims jwt.MapClaims, key string) []string {
	if val, ok := claims[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			strSlice := make([]string, len(slice))
			for i, v := range slice {
				if str, ok := v.(string); ok {
					strSlice[i] = str
				}
			}
			return strSlice
		}
	}
	return nil
}
