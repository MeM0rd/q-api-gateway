package auth

type RegisterModel struct {
	Email    string `json:"email"`
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
