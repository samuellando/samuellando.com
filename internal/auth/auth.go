package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const MAX_TOKEN_AGE = time.Hour

type header struct {
	Alg string
	Typ string
}

type payload struct {
	ValidTo time.Time
}

func getPrivateKey() *rsa.PrivateKey {
	b64Key := os.Getenv("RSA_KEY")
	decoded, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(decoded)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if validCredentials(req) {
		updateJWT(w)
		http.Redirect(w, req, "/", 303)
		return
	} else {
		fmt.Fprintf(w, "Invalid login dredentials")
		http.Redirect(w, req, "/", 401)
		return
	}
}

func Deauthenticate(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", 303)
}

func IsAuthenticated(req *http.Request, ws ...http.ResponseWriter) bool {
	if cookie, err := req.Cookie("session"); err == nil && validJWT(cookie.Value) {
		if len(ws) > 0 {
			updateJWT(ws[0])
		}
		return true
	} else {
		return false
	}
}

func updateJWT(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    createJWT(),
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(MAX_TOKEN_AGE),
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func createJWT() string {
	header := header{Alg: "RS256", Typ: "JWT"}
	payload := payload{ValidTo: time.Now().Add(MAX_TOKEN_AGE)}
	headerB, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}
	payloadB, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	headerB64 := base64.RawURLEncoding.EncodeToString(headerB)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadB)
	hash := sha256.Sum256([]byte(headerB64 + "." + payloadB64))
	sig, err := rsa.SignPKCS1v15(rand.Reader, getPrivateKey(), crypto.SHA256, hash[:])
	if err != nil {
		panic(err)
	}
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)
	token := headerB64 + "." + payloadB64 + "." + sigB64
	return token
}

func validJWT(jwt string) bool {
	public := getPrivateKey().Public().(*rsa.PublicKey)
	// Split the JWT
	parts := strings.Split(jwt, ".")
	// Calculate the hash of the first two parts
	hash := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	// Base64 decode all the parts
	headerB, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		log.Println(err)
		return false
	}
	payloadB, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println(err)
		return false
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		log.Println(err)
		return false
	}
	// Verify the hash to the signature
	if err := rsa.VerifyPKCS1v15(public, crypto.SHA256, hash[:], sig); err != nil {
		log.Println(err)
		return false
	}
	// Verify the header
	header := header{}
	err = json.Unmarshal(headerB, &header)
	if err != nil {
		log.Println(err)
		return false
	}
	if header.Alg != "RS256" || header.Typ != "JWT" {
		log.Println("Invalid header values")
		return false
	}
	// Verify the payload
	payload := payload{}
	err = json.Unmarshal(payloadB, &payload)
	if err != nil {
		log.Println(err)
		return false
	}
	if payload.ValidTo.Unix() < time.Now().Unix() {
		log.Println("Token is expired")
		return false
	}
	return true
}

func validCredentials(req *http.Request) bool {
	reqUser := req.PostFormValue("user")
	reqPassword := req.PostFormValue("password")
	htpasswd := os.Getenv("ADMIN_HTPASSWD")
	adminUser := strings.Split(htpasswd, ":")[0]
	adminPassword := strings.Split(htpasswd, ":")[1]
	if reqUser == adminUser {
		if err := bcrypt.CompareHashAndPassword([]byte(adminPassword), []byte(reqPassword)); err == nil {
			return true
		}
	}
	return false
}
