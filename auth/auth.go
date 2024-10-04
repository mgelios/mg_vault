package auth

import (
	"fmt"
	"log"
	"mg_vault/model"
	"mg_vault/storage"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil, jwt.WithAcceptableSkew(6000*time.Second))

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": 123})

	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func GetTokenForUser(user model.User) string {
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"id": user.Id})

	return tokenString
}

func ProcessLoginRequest(loginRequestBody model.LoginRequest) string {
	user, err := storage.GetUserByUsername(loginRequestBody.Username)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error during processing retriving user by login")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequestBody.Password))

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error during processing comparison of the password")
		return ""
	}

	return GetTokenForUser(user)
}
