package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/form3tech-oss/jwt-go"
)

type TokenExtractor func(r *http.Request) (string, error)
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err string)

type Options struct {
	ValidationKeyGetter jwt.Keyfunc
	Extractor           TokenExtractor
	ErrorHandler        ErrorHandler
	SigningMethod       jwt.SigningMethod
}

type JWTMiddleware struct {
	Options Options
}

type JWTContextKey struct{}

type AuthClaims struct {
	FacilityId int    `json:"fid,omitempty"`
	AccountId  int    `json:"aid,omitempty"`
	Category   int    `json:"cat,omitempty"`
	DeviceId   string `json:"did,omitempty"`
	jwt.StandardClaims
}

func defaultErrorHandler(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

func ExtractFromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Not Bearer token format")
	}
	return authHeaderParts[1], nil
}

func NewJWtMiddleware(options ...Options) *JWTMiddleware {
	var opts Options
	if len(options) == 0 {
		opts = Options{}
	} else {
		opts = options[0]
	}

	if opts.ErrorHandler == nil {
		opts.ErrorHandler = defaultErrorHandler
	}
	if opts.Extractor == nil {
		opts.Extractor = ExtractFromAuthHeader
	}
	return &JWTMiddleware{
		Options: opts,
	}
}

func (m *JWTMiddleware) HandleJWT(w http.ResponseWriter, r *http.Request) (AuthClaims, error) {
	authClaims := AuthClaims{}
	if r.Method == "OPTIONS" {
		return authClaims, nil
	}

	token, err := m.Options.Extractor(r)
	log.Printf("extract token :%s", token)

	if err != nil {
		m.Options.ErrorHandler(w, r, err.Error())
		return authClaims, fmt.Errorf("Error extract token: %v", err)
	}

	if token == "" {
		errorMsg := "token not found"
		m.Options.ErrorHandler(w, r, errorMsg)
		return authClaims, fmt.Errorf(errorMsg)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &authClaims, m.Options.ValidationKeyGetter)
	log.Printf("%+v", authClaims)

	if err != nil {
		m.Options.ErrorHandler(w, r, err.Error())
		return authClaims, fmt.Errorf("parse token error %v", err)
	}

	if m.Options.SigningMethod != nil &&
		m.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
		m.Options.ErrorHandler(w, r, errors.New("signing method error").Error())
		return authClaims, fmt.Errorf("signing method error")
	}

	if !parsedToken.Valid {
		m.Options.ErrorHandler(w, r, "token is invalid")
		return authClaims, errors.New("token is invalid")
	}

	return authClaims, nil
}

func (m *JWTMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authClaims, err := m.HandleJWT(w, r)
		if err != nil {
			m.Options.ErrorHandler(w, r, err.Error())
			return
		}
		if authClaims.AccountId != 0 {
			ctx := context.WithValue(r.Context(), JWTContextKey{}, authClaims)
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	})
}
