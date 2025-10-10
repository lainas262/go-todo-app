package api

type RequestCreateUser struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RequestLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseLoginUser struct {
	AccessToken string `json:"token"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}
