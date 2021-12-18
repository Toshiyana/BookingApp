package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Toshiyana/BookingApp/internal/config"
	"github.com/Toshiyana/BookingApp/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	// use session in handlers
	session = scs.New()
	session.Lifetime = 24 * time.Hour // last for 24 hours
	// set cookie parameters because all session use cookies in one form.
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false // In test, set false because we use test server.

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}
