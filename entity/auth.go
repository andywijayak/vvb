package entity

//Структура успешного ответа
type LoginSuccessResponse struct {
	Status string        `json:"status"`
	User   UserPublicDTO `json:"user"`
	Token  string
}

//Публичная версия пользователя без секретных данных
type UserPublicDTO struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

//Общая структура под все ошибки
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

//Структура входящих данных, если POST
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
