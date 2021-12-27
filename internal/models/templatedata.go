package models

import "github.com/Toshiyana/BookingApp/internal/forms"

// TemplateData holds data sent from holders to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string // security token
	Flash           string // success message
	Warning         string // warning message
	Error           string // error message
	Form            *forms.Form
	IsAuthenticated int
}
