package entity

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password`
}

type RegisterResponse struct {
	Status string        `json:"status"`
	User   UserPublicDTO `json:"user"`
}
