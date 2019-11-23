package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

var flagProfile string

const (
	FieldKeyRequestID  = "REQUEST_ID"
	FieldKeyRequestURL = "REQUEST_URL"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	flag.StringVar(&flagProfile, "profile", "default", "Run profile")
}

func main() {
	flag.Parse()

	initLogger()
	initProfileVariables()
	startHTTPServing()
}

func initLogger() {
	if flagProfile == "default" {
		logger.SetFormatter(&logger.JSONFormatter{})
	}
}

func initProfileVariables() {
	logger.Info("Starting application with profile: ", flagProfile)

	if flagProfile != "default" {
		profileFileName := "profiles/" + flagProfile + ".env"
		if godotenv.Load(profileFileName) != nil {
			logger.Fatal("Error loading environment variables from: ", profileFileName)
		} else {
			logger.Info("Environment variables loaded from: ", profileFileName)
		}
	}
}

func port() string {
	port := "80"

	if flagProfile != "default" && os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return port
}

func startHTTPServing() {
	http.HandleFunc("/v1/users/", userByID)
	http.HandleFunc("/readiness", health)
	http.HandleFunc("/health", health)

	port := port()
	logger.Info("Starting server at port: ", port)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}

func health(w http.ResponseWriter, r *http.Request) {
}

func userByID(w http.ResponseWriter, r *http.Request) {
	log := GetLogger(r)
	if r.Method != http.MethodGet {
		log.Error("method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		log.Error("can't convert url path param for id, ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := User{
		ID:   id,
		Name: "some-user",
	}
	json, err := json.Marshal(result)
	if err != nil {
		log.Error("can't marshall user object to json, ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)

	log.Info("check success")
}

// GetLogger is function to get log information from request header and put them to logger
func GetLogger(r *http.Request) *logger.Entry {
	reqID := r.Header.Get(FieldKeyRequestID)
	if reqID == "" {
		reqID = uuid.New().String()
	}

	return logger.WithFields(logger.Fields{
		FieldKeyRequestID:  reqID,
		FieldKeyRequestURL: r.RequestURI,
	})
}
