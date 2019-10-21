package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"test2/model"

	"github.com/sirupsen/logrus"
)

func randStr10() string {
	buff := make([]byte, 10)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:10]
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func createMessage(payload *model.Event) []byte {
	data, err := json.Marshal(payload)
	checkError(err)
	return data
}

func getLogger() logrus.FieldLogger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}
