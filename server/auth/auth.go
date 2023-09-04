package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID int64
}

func CreateToken(key string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": claims.ID,
		})

	s, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return s, nil
}
