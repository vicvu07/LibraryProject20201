package model

import jwt "github.com/dgrijalva/jwt-go"

// Claim jwt claim
type Claim struct {
	Username   string
	Group      []string
	Permission [][]string
	jwt.StandardClaims
}
