package handler

import (
	"encoding/json"
	"net/http"
	"privateCabin/entity"
	"privateCabin/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// @Summary GetUser user
// @Description Авторизация с проверкой логина и пароля
// @Tags Авторизация
// @Param input body entity.LoginRequest true "Login / Password"
// @Accept json
// @Success 200 {object} entity.LoginSuccessResponse
// @Failure 400 {object} entity.ErrorResponse
// @Router /login [post]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.ErrorResponse{
			Status: "error",
			Error:  "invalid request",
		})
		return
	}

	user, err := h.service.GetUser(req.Login, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(entity.ErrorResponse{
			Status: "error",
			Error:  "invalid login or password",
		})
		return
	}

	token := "mock-access-token"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.LoginSuccessResponse{
		Status: "sucess",
		User:   *user,
		Token:  token,
	})

}

// @Summary CreateUser user
// @Description Регистрация нового пользователя с логином и паролем
// @Tags Регистрация
// @Accept json
// @Produce json
// @Param input body entity.RegisterUser true "Login / Password"
// @Success 200 {object} entity.RegisterResponse
// @Failure 400 {object} entity.ErrorResponse
// @Router /register [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req entity.RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.ErrorResponse{
			Status: "error",
			Error:  "invalid request",
		})
		return
	}

	user, err := h.service.CreateUser(req.Login, req.Password)
	if err != nil {

		if err.Error() == "login_already_exists" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Status: "error",
				Error:  "login already exists",
			})
			return
		}
		////OTHER ERROR
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	// УСПЕШНЫЙ ВАРИАНТ — ЭТО ВАЖНО
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entity.RegisterResponse{
		Status: "success",
		User:   *user,
	})
}
