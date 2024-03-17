package mtauthserver

type Config struct {
	LdapServer     string
	LdapServer2    string
	LdapTLSVerify  bool
	LdapTLSVerify2 bool
	BindDN         string
	BindPassword   string
	SearchBase     string
	UserFilter     string
	GroupFilter    string
	SecretKey      string
	LogtoFile      bool
	LogFile        string
	Debug          bool
}

type UsernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type ResponseToken struct {
	Groups   []string `json:"groups"`
	Status   string   `json:"status"`
	Username string   `json:"username"`
	Exp      int      `json:"exp"`
	Token    string   `json:"token"`
}

type ResponseVerify struct {
	AppName  string   `json:"app_name"`
	Groups   []string `json:"groups"`
	Status   string   `json:"status"`
	Username string   `json:"username"`
	Exp      int      `json:"exp"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseAuthError struct {
	Message string `json:"msg"`
	Status  string `json:"status"`
}
