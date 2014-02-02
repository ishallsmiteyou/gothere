package handlers

import (
    "net/http"
    "database/sql"
    "gothere/templates"
    "gothere/models"
    "gothere/database"
    "gothere/utils"
    "gothere/password"
)

type oldValues struct {
    Flag        bool
    Firstname   string
    Lastname    string
    Email       string
}

func RegisterGet(w http.ResponseWriter) {
    /*
     * /register GET method handler.
     * Just render's the form.
     */

    templates.Render(w, "register", nil)
}

func RegisterPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    /*
     * /register POST method handler.
     * Validates the form,
     * check's if username is availible,
     * and then creates a user and redirects to
     * /login .
     */

    var user models.User
    user.Email = r.FormValue("email")
    user.Password = r.FormValue("password")
    user.Firstname = r.FormValue("firstname")
    user.Lastname = r.FormValue("lastname")

    repeat := r.FormValue("repeat")

    var old oldValues
    old.Flag = true
    old.Firstname = user.Firstname
    old.Lastname = user.Lastname
    old.Email = user.Email

    if ! utils.UserValidate(user, repeat) {
        templates.Render(w, "register", old)
    } else if database.GetPassword(db, user.Email) != "" {

        templates.Render(w, "register", old)
    } else {
        // Creates a user in the db.
        user.Password = password.NewPassword(user.Password)
        database.CreateUser(db, user)
        http.Redirect(w, r, "/login", http.StatusFound)
    }
}

