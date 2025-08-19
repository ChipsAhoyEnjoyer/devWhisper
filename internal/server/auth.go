package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var hashCost = getHashCost()

const (
	issuer = "devWhisper"
)

func getHashCost() int {
	godotenv.Load(".env")
	costStr := os.Getenv("HASH_COST")
	if costStr == "" {
		log.Println("Warning: HASH_COST environment variable is not a valid integer, using default cost of 10")
		costStr = "10" // default cost
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		log.Printf("Warning: error using strconv.Atoi to convert hash cost to integer: %v\n", err)
		cost = 10
	}
	return cost
}

func hashPassword(password string) (hashedPW string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(b), err
}

func CheckPasswordHash(hash, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tk := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		},
	)
	tkString, err := tk.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tkString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tk, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(tk *jwt.Token) (any, error) {
			if tk.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("error incorrect signing method given: " + tk.Method.Alg())
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	tokenIssuer, err := tk.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if tokenIssuer != string(issuer) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	idStr, err := tk.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user id")
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if !strings.HasPrefix(bearer, "Bearer ") {
		return "", errors.New("no bearer token found")
	}
	return strings.TrimPrefix(bearer, "Bearer "), nil
}

func MakeRefreshToken() (string, error) {
	tk := make([]byte, 32)
	_, err := rand.Read(tk)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tk), nil
}

func validateUsername(username string) error {
	nameLen := len(username)
	if nameLen < 5 || nameLen > 15 {
		return errors.New("invalid username length; name must be between between 5-15 characters")
	}

	for _, r := range username {
		if !isAllowedRune(r) {
			return errors.New("invalid character in username: " + string(r))
		}
	}

	return nil
}

func validatePassword(password string) error {
	passLen := len(password)
	if passLen < 5 || passLen > 15 {
		return errors.New("invalid password length; password must be between 5-15 characters")
	}
	for _, r := range password {
		if !isAllowedRune(r) {
			return errors.New("invalid character in password: " + string(r))
		}
	}
	return nil

}
func isAllowedRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= '0' && r <= '9':
		return true
	case r == '_' || r == ' ':
		return true
	default:
		return false
	}
}
