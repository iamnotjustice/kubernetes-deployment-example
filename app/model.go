package main

//IPGeoData is the response we get from ipstack API
type IPGeoData struct {
	IPGeoDataErr
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentName string  `json:"continent_name"`
	CountryName   string  `json:"country_name"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag      string `json:"country_flag"`
		CountryFlagEmoji string `json:"country_flag_emoji"`
		CallingCode      string `json:"calling_code"`
	} `json:"location"`
}

type IPGeoDataErr struct {
	Success bool `json:"success"`
	Error   *struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}
