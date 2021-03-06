package handlers

import (
	"fmt"
	"github.com/nymoral/gothere/cookies"
	"github.com/nymoral/gothere/database"
	"github.com/nymoral/gothere/models"
	"github.com/nymoral/gothere/templates"
	"github.com/nymoral/gothere/utils"
	"log"
	"net/http"
)

func GuessesGet(w http.ResponseWriter, r *http.Request) {
	// /guesses handler for GET method request.
	// Renders a page only for users with valid sessionid cookie.
	// All the rest are redirected to /login .

	db := database.GetConnection()

	sessionid := cookies.GetCookieVal(r, "sessionid")
	username := cookies.UsernameFromCookie(sessionid)
	pk, is_admin := database.GetPkAdmin(db, username)

	if username == "" || pk == -1 {
		// Gorilla failed to decode it.
		// Or encoded username does not exist in db.
		http.Redirect(w, r, "/login/", http.StatusFound)
	} else if is_admin {
		http.Redirect(w, r, "/admin/", http.StatusFound)
	} else {
		// Fetches users guesses from the db and gets data for
		// result submit dropbox.

		var F models.GuessContext
		F.OpenGames = database.GamesList(db, "open")
		F.Guesses = database.UsersGuesses(db, pk)
		F.Error = false
		templates.Render(w, "guesses", F)
	}

}

func GuessesPost(w http.ResponseWriter, r *http.Request) {
	// /guesses POST method.
	// Checks if user trying to submit is in valid.

	db := database.GetConnection()

	sessionid := cookies.GetCookieVal(r, "sessionid")
	username := cookies.UsernameFromCookie(sessionid)

	var guess models.Guess
	guess.Userpk, _ = database.GetPkAdmin(db, username)

	var F models.GuessContext
	F.OpenGames = database.GamesList(db, "open")
	F.Guesses = database.UsersGuesses(db, guess.Userpk)
	F.Error = false

	if username == "" || guess.Userpk < 0 {
		// Gorilla failed to decode it.
		http.Redirect(w, r, "/login/", http.StatusFound)
	} else {
		var nr int
		var err error
		//Extract data from request and check if form is valid.
		if utils.ExtractResult(r.FormValue("result_2"), &guess.Result1, &guess.Result2) ||
			utils.ExtractResult(r.FormValue("result_1"), &guess.Result1, &guess.Result2) {
			if guess.Result1 < 0 || guess.Result2 < 0 {
				F.Error = true
			}
		} else {
			nr, err = fmt.Sscanf(r.FormValue("result_1"), "%d", &guess.Result1)
			if nr != 1 || err != nil || guess.Result1 < 0 {
				F.Error = true
			}
			nr, err = fmt.Sscanf(r.FormValue("result_2"), "%d", &guess.Result2)
			if nr != 1 || err != nil || guess.Result2 < 0 {
				F.Error = true
			}
		}
		nr, err = fmt.Sscanf(r.FormValue("game-id"), "%d", &guess.Gamepk)
		if nr != 1 || err != nil {
			F.Error = true
		}
		if F.Error {
			templates.Render(w, "guesses", F)
		} else {
			// Submit a guess.
			database.GiveResult(db, &guess)
			http.Redirect(w, r, "/guesses/", http.StatusFound)
			log.Printf("GUESS BY (%d). GAME (%d)\n", guess.Userpk, guess.Gamepk)
		}
	}

}
