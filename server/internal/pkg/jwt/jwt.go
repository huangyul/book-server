package jwt

import "github.com/golang-jwt/jwt/v5"

type JWT interface {
	generateToken(claims jwt.Claims) (string, string, error)
}

var _ JWT = (*jwtService)(nil)

type jwtService struct {
}

// generateToken implements JWT.
func (j *jwtService) generateToken(claims jwt.Claims) (string, string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if err != nil {
		return "", "", err
	}
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := token.SignedString([]byte("secret"))
}
