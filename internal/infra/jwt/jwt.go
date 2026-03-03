package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/victorbetoni/justore/app-manager/internal/config"
	jwtc "gopkg.in/dgrijalva/jwt-go.v3"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

func ExtractUserIdentifier(tokenStr, ip, userAgent string) (string, int, error) {

	token, err := jwtc.Parse(tokenStr, func(token *jwtc.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtc.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwtc.MapClaims); ok && token.Valid {

		isAdm, _ := strconv.Atoi(fmt.Sprintf("%v", claims[config.GetConfig().Jwt.ClaimKeys.AdminFlag]))
		tokenAgent, _ := claims[config.GetConfig().Jwt.ClaimKeys.UserAgent]
		tokenIp, _ := claims[config.GetConfig().Jwt.ClaimKeys.Ip]

		if config.GetConfig().Jwt.CheckIp && tokenIp != ip {
			return "", 0, errors.New("unauthorized")
		}

		if config.GetConfig().Jwt.CheckUserAgent && tokenAgent != tokenAgent {
			return "", 0, errors.New("unauthorized")
		}

		return fmt.Sprintf("%v", claims[config.GetConfig().Jwt.ClaimKeys.Identifier]), isAdm, nil
	}

	return "", 0, errors.New("unauthorized")
}

func DefineKeyPair(priv []byte, pub []byte) {
	privateKey, err := jwtc.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		log.Fatalf("Error while loading PRIVATE key: %s\n", err)
	}
	signKey = privateKey
	publicKey, err := jwtc.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		log.Fatalf("Error while loading PUBLIC key: %s\n", err)
	}
	verifyKey = publicKey
}
