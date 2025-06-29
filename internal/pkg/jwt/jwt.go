package jwt

import (
	"net-http-boilerplate/internal/config"
	"net-http-boilerplate/internal/pkg/encrypt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	config config.JWT
}

func NewJWT(cfg config.JWT) *JWT {
	return &JWT{config: cfg}
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (j *JWT) GenerateToken(id string, email string) (string, string, error) {
	encryptedID, err := encrypt.EncryptData(id)
	if err != nil {
		return "", "", nil
	}

	encryptedEmail, err := encrypt.EncryptData(email)
	if err != nil {
		return "", "", nil
	}

	encryptedIssuer, err := encrypt.EncryptData(j.config.Issuer)
	if err != nil {
		return "", "", nil
	}

	claims := Claims{
		ID:    encryptedID,
		Email: encryptedEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    encryptedIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.RegisteredClaims{
		Issuer:    encryptedIssuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Refresh token valid for 7 days
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	// Extract encrypted strings
	encryptedID := claimsMap["id"].(string)
	encryptedEmail := claimsMap["email"].(string)

	// Decrypt
	decryptedID, err := encrypt.DecryptData(encryptedID)
	if err != nil {
		return nil, err
	}

	decryptedEmail, err := encrypt.DecryptData(encryptedEmail)
	if err != nil {
		return nil, err
	}

	// Convert numeric claims
	iat := int64(claimsMap["iat"].(float64))
	exp := int64(claimsMap["exp"].(float64))

	return &Claims{
		ID:    decryptedID,
		Email: decryptedEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claimsMap["iss"].(string),
			IssuedAt:  jwt.NewNumericDate(time.Unix(iat, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
		},
	}, nil
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", jwt.ErrTokenExpired
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}
