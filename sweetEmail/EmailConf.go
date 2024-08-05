package sweetEmail

type EmailConf struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	EmailName string `json:"emailname"`
}
