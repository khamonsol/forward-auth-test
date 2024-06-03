package util

import "github.com/golang-jwt/jwt/v4"

func GetStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		return val.(string)
	}
	return ""
}

func GetInt64Claim(claims jwt.MapClaims, key string) int64 {
	if val, ok := claims[key]; ok {
		return int64(val.(float64))
	}
	return 0
}

func GetStringSliceClaim(claims jwt.MapClaims, key string) []string {
	if val, ok := claims[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, len(slice))
			for i, v := range slice {
				result[i] = v.(string)
			}
			return result
		}
	}
	return []string{}
}
