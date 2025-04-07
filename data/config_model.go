package data

type Config struct {
	Secret string `json:"secret"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Port   string `json:"port"`
	DBName string `json:"dbname"`
	Host   string `json:"host"`
}
