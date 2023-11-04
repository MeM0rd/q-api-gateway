package auth

type LoginModel struct {
	Email    string `json:"email"`
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
