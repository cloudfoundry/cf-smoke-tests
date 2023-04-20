package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/env", env)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	response := fmt.Sprintf(`Healthy
	It just needed to be restarted!
	My application metadata: %v
	My port: %v
	My custom env variable: %v`, os.Getenv("VCAP_APPLICATION"), os.Getenv("PORT"), os.Getenv("CUSTOM_VAR"))
	fmt.Fprintln(res, response)
}

func env(res http.ResponseWriter, req *http.Request) {
	envVariables := make(map[string]string)
	for _, envKeyValue := range os.Environ() {
		keyValue := strings.Split(envKeyValue, "=")
		envVariables[keyValue[0]] = keyValue[1]
	}
	envJsonBytes, err := json.Marshal(envVariables)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(res, string(envJsonBytes))
}
