package handlers

import (
	"net/http"

	"github.com/Toshiyana/BookingApp/pkg/config"
	"github.com/Toshiyana/BookingApp/pkg/models"
	"github.com/Toshiyana/BookingApp/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHanlders set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// All of the handlers access to the repository using receivers(m *Repository)
// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr // store visitor's ip addres
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTamplate(w, "home.page.tmpl", &models.TemplateData{})
	// fmt.Fprintf(w, "This is the home page")
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, world."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTamplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

	// sum := addValues(2, 2)
	// _, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page and 2 + 2 is %d", sum))
}

// func addValues(x, y int) int {
// 	return x + y
// }

// // Divide is the divide page handler
// func Divide(w http.ResponseWriter, r *http.Request) {
// 	dividend := 100.0
// 	divisor := 0.0
// 	f, err := divideValues(dividend, divisor)
// 	if err != nil {
// 		fmt.Fprintf(w, "cannot divide by zero")
// 		return
// 	}

// 	fmt.Fprintf(w, fmt.Sprintf("%f divide by %f is %f", dividend, divisor, f))
// }

// func divideValues(x, y float64) (float64, error) {
// 	if y == 0 {
// 		err := errors.New("cannot divide by zero")
// 		return 0, err
// 	}
// 	result := x / y
// 	return result, nil
// }
