package handlers

import (
    "log"
    "net/http"
    "gothere/models"
    "gothere/database"
    "gothere/password"
    "gothere/templates"
)

func RegisterGet(w http.ResponseWriter) {
    // /register GET method handler.
    // Just render's the form.

    templates.Render(w, "register", nil)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
    // /register POST method handler.
    // Validates the form,
    // check's if username is availible,
    // and then creates a user and redirects to
    // /login .

    db := database.GetConnection()

    var user models.User
    // Model out of form data.
    user.Email = r.FormValue("email")
    user.Password = r.FormValue("password")
    user.Firstname = r.FormValue("firstname")
    user.Lastname = r.FormValue("lastname")

    repeat := r.FormValue("repeat")

    var old models.RegisterContext
    // Model for return form.
    // In case there the data wasn't valid
    old.Firstname = user.Firstname
    old.Lastname = user.Lastname
    old.Email = user.Email
    old.Flag =  user.UserValidate(repeat)

    if  old.Flag != "" {
        templates.Render(w, "register", old)
        return
    }

    pass, _ := database.GetPassword(db, user.Email)
    // Checks if user exists.
    if pass != "" {
        old.Flag = "Vartotojas su šiuo el. pašto adresu jau egzistuoja."
        templates.Render(w, "register", old)
        return
    }

    user.Password = password.NewPassword(user.Password)
    database.CreateUser(db, &user)
    // Creates a user in the db.
    http.Redirect(w, r, "/login", http.StatusFound)
    log.Printf("USER CREATED (%s)\n", user.Email)
}

