package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
	"os"
)

type ClimaResponse struct {
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
	Kelvin     float64 `json:"temp_k"`
}

type CepResponse struct {
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

var validate = validator.New()

func main() {
	http.HandleFunc("/clima", weatherHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if err := validate.Var(cep, "len=8,numeric"); err != nil {
		http.Error(w, "CEP INCORRETO", http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocationFromCep(cep)
	if err != nil {
		http.Error(w, "NÃO FOI POSSÍVEL LOCALIZAR CEP", http.StatusNotFound)
		return
	}

	fmt.Printf("testando aqui e verificando: %+v", location)
	weather, err := getWeather(location.Localidade, location.Uf)
	if err != nil {
		http.Error(w, "Falha ao tentar obter temperatura", http.StatusInternalServerError)
		return
	}

	response := ClimaResponse{
		Celsius:    weather.Celsius,
		Fahrenheit: weather.Celsius*1.8 + 32,
		Kelvin:     weather.Celsius + 273.15,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getLocationFromCep(cep string) (CepResponse, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return CepResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CepResponse{}, fmt.Errorf("invalid cep")
	}

	var location CepResponse
	err = json.NewDecoder(resp.Body).Decode(&location)
	return location, err
}

func getWeather(city string, uf string) (ClimaResponse, error) {
	cityParam := city
	stateParam := uf
	query := url.QueryEscape(cityParam + "," + stateParam)

	apiKey := ""
	urlNew := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, query)

	println(urlNew)
	resp, err := http.Get(urlNew)
	if err != nil {
		return ClimaResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ClimaResponse{}, fmt.Errorf("falha ao tentar pegar o clima")
	}

	var weather ClimaResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	return weather, err
}
