package controller

import (
	"fmt"
	"net/http"
	"time"

	// Asegúrate de importar el paquete donde están tus templates generados
	"github.com/alejandro-bustamante/sancho/server/ui/views"
)

type UserHandler struct {
	userService UserService // Tu interfaz inyectada
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 1. EL ESQUELETO BASE (Punto de entrada)
// Esta función decide qué renderizar basándose en si hay una sesión activa o no.
func (h *UserHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sancho_session")

	// Si no hay cookie o está vacía, mostramos la pantalla de login
	if err != nil || cookie.Value == "" {
		views.Login("").Render(r.Context(), w)
		return
	}

	// Si hay sesión, mostramos la interfaz principal
	views.Main(cookie.Value).Render(r.Context(), w)
}

// 2. AUTENTICACIÓN
func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// Con HTMX (hx-post), los datos vienen en el formulario, no en JSON
	if err := r.ParseForm(); err != nil {
		views.Login("Error al procesar el formulario").Render(r.Context(), w)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	authenticated, err := h.userService.AuthenticateUser(r.Context(), username, password)
	if err != nil || !authenticated {
		// Si hay error, volvemos a renderizar TODA la vista de login, pero inyectando el mensaje de error
		views.Login("Usuario o contraseña incorrectos").Render(r.Context(), w)
		return
	}

	// Login exitoso: creamos una cookie de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "sancho_session",
		Value:    username, // En un entorno real, aquí iría un token seguro (JWT) o un ID de sesión
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Redirigimos a la raíz ("/"). HTMX (gracias a hx-boost) interceptará esto
	// y cargará el HandleIndex transparentemente.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 3. CERRAR SESIÓN (Necesario para el ciclo completo)
func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sancho_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour), // Expiramos la cookie inmediatamente
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 4. REGISTRO (Refactorizado para Formularios)
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// Aquí deberías tener un template views.Register("Error")
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	_, err := h.userService.RegisterUser(r.Context(), username, password, email)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		// Renderizar vista de registro con error
		return
	}

	// Registro exitoso, redirigimos al login
	http.Redirect(w, r, "/login", http.StatusSeeOther) // O a la raíz "/"
}

// Funciones pendientes de refactorizar según las vistas que vayas creando
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<div>Placeholder: Usuario eliminado</div>"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<div>Placeholder: Usuario actualizado</div>"))
}
