package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type configData struct {
	ApiKey string `json:"ApiKey"`
}

type weatherResponse struct {
	Name string `json:"name"`

	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`

	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Rain struct {
		H float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		Id      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`

	Main struct {
		Kelvin    float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
}

func main()  {
http.HandleFunc("/weather/",weather)
http.ListenAndServe(":9090",nil)
}

func weather(w http.ResponseWriter , r *http.Request) {
	city := strings.SplitN(r.URL.Path,"/",3)[2]
	data , err := query(city)
	if err != nil {
		http.Error(w ,err.Error(),http.StatusInternalServerError)
		return
	}
	//convert kelvin to centigrade
	data.Main.Kelvin -= 273.15
	w.Header().Set("Content-type","application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)

}

func query(city string) (weatherResponse , error)  {
	//load config file
	apiConfig , err := loadConfig("config.json")
	if err != nil {
		return weatherResponse{}, err
	}
	//call webservice request
	resp , err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" +city+ "&APPID="+apiConfig.ApiKey)
	if err != nil {
		return weatherResponse{}, err
	}

	defer resp.Body.Close()

	var r weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&r);err != nil{
		return weatherResponse{}, err
	}
	return r , nil
}
func loadConfig(filename string) (configData , error)  {
	bytes , err := os.ReadFile(filename)
	if err != nil {
		return configData{}, err
	}
	var c configData
	err = json.Unmarshal(bytes , &c)
	if err != nil {
		return configData{}, err
	}

	return c ,nil
}
