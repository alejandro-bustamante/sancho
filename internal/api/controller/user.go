package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type AuthenticateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	userService UserService // Inyectamos la interfaz
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "<div class='error'>Solicitud inválida</div>", http.StatusBadRequest)
		return
	}

	_, err := h.userService.RegisterUser(r.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		http.Error(w, "<div class='error'>No se pudo registrar el usuario</div>", http.StatusBadRequest)
		fmt.Println("error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("<div class='success'>Usuario registrado exitosamente. (Redirigir a login)</div>"))
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var req AuthenticateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "<div class='error'>Solicitud inválida</div>", http.StatusBadRequest)
		return
	}

	authenticated, err := h.userService.AuthenticateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "<div class='error'>Error de autenticación</div>", http.StatusUnauthorized)
		return
	}

	if !authenticated {
		http.Error(w, "<div class='error'>Usuario o contraseña incorrectos</div>", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div class='success'>Sesión iniciada correctamente (Aquí puedes devolver un header HX-Redirect para ir al Dashboard)</div>"))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Placeholder: Usuario eliminado</div>"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Placeholder: Usuario actualizado</div>"))
}
