package auth

import (
	"encoding/json"
	"log/slog"
	"mg_vault/model"
	"mg_vault/storage"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil, jwt.WithAcceptableSkew(6000*time.Second))
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(hashedPassword))
	// _, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": 123})
	// slog.Debug("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func ProcessLoginRequest(w http.ResponseWriter, r *http.Request) {
	var loginRequestBody model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequestBody)
	if err != nil {
		slog.Error(loginRequestBody.Username)
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := storage.GetUserByUsername(loginRequestBody.Username)
	if err != nil {
		slog.Error(err.Error())
		slog.Error("Error during processing retriving user by username")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequestBody.Password))
	if err != nil {
		slog.Error(err.Error())
		slog.Error("Error during processing comparison of the password")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_, token, _ := TokenAuth.Encode(map[string]interface{}{"id": user.Id})

	cookie := http.Cookie{
		Name:  "jwt",
		Path:  "/",
		Value: token,
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/index")
}
