package tasmotaupdater

type tasmota struct {
	Ip      string `json:"ip"`
	Name    string `json:"dn"`
	Version string `json:"sw"`
}

var tasmotas = make(map[string]tasmota)
