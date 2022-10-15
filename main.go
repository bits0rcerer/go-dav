package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/gorilla/handlers"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	RootDirKey   = "GODAV_ROOT"
	URLPrefixKey = "GODAV_PREFIX"
	NoAuthKey    = "GODAV_NO_AUTH"
	PortKey      = "PORT"

	UserKeyPrefix = "GODAV_USER_"
)

func parseUsers() map[string][]byte {
	userMap := make(map[string][]byte)

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, UserKeyPrefix) {
			continue
		}

		splits := strings.SplitN(env, "=", 2)
		if len(splits) != 2 {
			continue
		}
		key := splits[0]
		password := splits[1]

		user := strings.TrimPrefix(key, UserKeyPrefix)
		if user == "" {
			log.Panicln("empty user is not allowed ->", env)
		}

		var err error
		userMap[user], err = hex.DecodeString(password)
		if err != nil {
			log.Panicln("invalid user pass hash ->", env)
		}
	}

	return userMap
}

func compareBytes(a, b []byte) bool {
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func auth(noAuth bool, userMap map[string][]byte, onSuccess http.HandlerFunc) http.HandlerFunc {
	if noAuth {
		log.Println("[!] Running WITHOUT authorization. Make sure you know what you are doing.")
		return onSuccess
	}

	if len(userMap) == 0 {
		log.Println("[!] Running WITHOUT any allowed users.")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", "Basic")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		authPair, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		authSplits := strings.SplitN(string(authPair), ":", 2)
		if len(authSplits) != 2 {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		user := authSplits[0]
		password := authSplits[1]

		hash := sha256.Sum256([]byte(password))
		if !compareBytes(userMap[user], hash[:]) {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		onSuccess(w, r)
	}
}

type RecoveryLogger struct{}

func (receiver RecoveryLogger) Println(errs ...interface{}) {
	log.Println("[!] handler panicked on:")
	for i, err := range errs {
		log.Printf("\t#%03d\t%v\n", i, err)
	}
}

func main() {
	rootPath := os.Getenv(RootDirKey)
	if rootPath == "" {
		log.Panicf("[!] $%s is not set. Consider $%s=/data when using Docker\n", RootDirKey, RootDirKey)
	}

	davHandler := webdav.Handler{
		Prefix:     os.Getenv(URLPrefixKey),
		FileSystem: webdav.Dir(rootPath),
		LockSystem: webdav.NewMemLS(),
	}

	noAuthEnv := strings.ToLower(os.Getenv(NoAuthKey))
	noAuth := noAuthEnv == "true" || noAuthEnv == "yes" || noAuthEnv == "1"

	mux := http.NewServeMux()
	mux.Handle("/",
		handlers.RecoveryHandler(handlers.RecoveryLogger(RecoveryLogger{}))(
			handlers.ProxyHeaders(
				handlers.CombinedLoggingHandler(log.Writer(), auth(noAuth, parseUsers(), davHandler.ServeHTTP)),
			),
		),
	)

	port := os.Getenv(PortKey)
	if port == "" {
		port = "8080"
	}

	log.Println("[*] Listening on", ":"+port)
	log.Panicln(http.ListenAndServe(":"+port, mux))
}
