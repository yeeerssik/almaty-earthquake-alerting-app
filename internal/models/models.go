package models

type Response struct {
	Type     string `json:"type"`
	Metadata struct {
		Generated int64  `json:"generated"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		Status    int    `json:"status"`
		API       string `json:"api"`
		Count     int    `json:"count"`
	} `json:"metadata"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Mag     float64 `json:"mag"`
			Place   string  `json:"place"`
			Time    int64   `json:"time"`
			Updated int64   `json:"updated"`
			Tz      any     `json:"tz"`
			URL     string  `json:"url"`
			Detail  string  `json:"detail"`
			Felt    any     `json:"felt"`
			Cdi     any     `json:"cdi"`
			Mmi     any     `json:"mmi"`
			Alert   any     `json:"alert"`
			Status  string  `json:"status"`
			Tsunami int     `json:"tsunami"`
			Sig     int     `json:"sig"`
			Net     string  `json:"net"`
			Code    string  `json:"code"`
			Ids     string  `json:"ids"`
			Sources string  `json:"sources"`
			Types   string  `json:"types"`
			Nst     int     `json:"nst"`
			Dmin    float64 `json:"dmin"`
			Rms     float64 `json:"rms"`
			Gap     int     `json:"gap"`
			MagType string  `json:"magType"`
			Type    string  `json:"type"`
			Title   string  `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		ID string `json:"id"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}
