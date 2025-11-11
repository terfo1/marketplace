package middleware

import (
	"PROJECTTEST/internal/helpers"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"strings"
	"time"
)

var JWTSecret []byte

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userID bson.ObjectID) (string, error) {
	claims := Claims{
		UserID: userID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return Token.SignedString(JWTSecret)
}

func GetUserIDFromToken(r *http.Request) (bson.ObjectID, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return bson.NilObjectID, errors.New("missing auth token")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 {
		return bson.NilObjectID, errors.New("invalid auth header")
	}
	tkn, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return bson.NilObjectID, err
	}
	claims, ok := tkn.Claims.(*Claims)
	if !ok || !tkn.Valid {
		return bson.NilObjectID, errors.New("invalid token")
	}
	id, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return bson.NilObjectID, err
	}
	return id, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUserIDFromToken(r)
		if err != nil {
			helpers.RespondError(w, http.StatusUnauthorized, "unauthorized: "+err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}
