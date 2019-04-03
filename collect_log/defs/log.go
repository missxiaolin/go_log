package defs

const (
	CONFIG_PATH = "./config/"
)

type Log struct {
	Url string `json:"url"`
	UrlId int `json:"url_id"`
	Type int `json:"type"`
	Ip string
}
