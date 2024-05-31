package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const apiKey = "4a3689591e7746a38fc120653242305"

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}


type TemperaturaResponse struct {
	Localidade string `json:"Localidade"`
	TemperaturaGraus float64 `json:"Temp_C"`
	TemperaturaFarenheit float64 `json:"Temp_F"`
	TemperaturaKelvin float64 `json:"Temp_K"`
}

// Struct para deserializar a resposta JSON
type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}


func HandleTemperatura(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]

	req, err := http.Get("http://viacep.com.br/ws/"+cep +"/json/")
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	defer req.Body.Close()

	var viaCep ViaCep
	if err := json.NewDecoder(req.Body).Decode(&viaCep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(viaCep.Localidade)
	location := viaCep.Localidade

	temp_c, err := getWeather(apiKey, location)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp_f := celsiusToFarenheit(temp_c)
	temp_k := celsiusToKelvin(temp_c)

	var temperaturaResponse TemperaturaResponse
	temperaturaResponse.Localidade = location
	temperaturaResponse.TemperaturaGraus = temp_c
	temperaturaResponse.TemperaturaFarenheit = temp_f
	temperaturaResponse.TemperaturaKelvin = temp_k

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperaturaResponse)

}

func getWeather(apiKey string, location string) (float64, error) {
	formattedLocation := url.QueryEscape(location)
	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, formattedLocation)
	
	resp, err := http.Get(urlWeather)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro ao chamar weather API com status: %s", resp.Status)
	}

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return 0, err
	}

	return weatherResp.Current.TempC, nil
}

func celsiusToFarenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}
func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/temperatura/{cep}", HandleTemperatura).Methods("GET")

	http.ListenAndServe(":8080", r)
}