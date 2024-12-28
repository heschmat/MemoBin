package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/heschmat/MemoBin/internal/models"
	"github.com/heschmat/MemoBin/internal/validator"
)

// N.B. struct fields must be exported (i.e., start with a capital letter)
// in order to be ready by the html/template package when rendering the template.
// This does NOT apply to maps, though.
type memoCreateForm struct {
	// Include struct tags: tell the decoder how to map HTML form values into struct fields.
	// e.g., name "title" in the form matches with field "Title" in the struct.
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validator.Validator `form:"-"`
}

// Hold the form data for user auth:
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	memos, err := app.memos.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Memos = memos
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}


func (app *application) memoView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi((r.PathValue("id")))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	memo, err := app.memos.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// If no matching record is found, return a 404 Not Found response.
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Memo = memo
	data.Flash = flash // pass the `flash` message to the template
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) memoCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new `memoCreateForm` instance and pass it to the template.
	// We also could set default values for the fields.
	data.Form = memoCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) memoCreatePost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes.
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	// if the limit above is hit, `MaxBytesreader` sets a flag on *http.ResponseWriter*
	// which instructs the server to close the underlying TCP connection.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// // Get the `expires` value from the *form* as normal.
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// form := memoCreateForm {
	// 	Title: r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// }

	var form memoCreateForm
	// The follwoing will essentially fill the struct with the relevant values fro the HTML form.
	// N.B. Type conversions are handled automatically too.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		// If there's a problem, return a 400 Bad Request response to the client.
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 chars long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank.")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "Permitted values: 1, 7, 365")

	// If there are any validation errors,
	// re-display the `create.tmpl.html` template, passing the `memoCreateForm` instance
	// as dynamic data in the Form field.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// To indicate that there was a validation error, we use the HTTP status code 442.
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.memos.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Add as string key:value - "flash":"Memo created successfully" to the session data.
	// if - and only if - the *memo* is created successfully.
	app.sessionManager.Put(r.Context(), "flash", "Memo created successfully.")

	// Redirect the user to the relevant page for the memo.
	http.Redirect(w, r, fmt.Sprintf("/memo/view/%d", id), http.StatusSeeOther)
}

// ============================================================================== #
// User Authentication
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-valued instance of *userSignupForm* struct.
	var form userSignupForm

	// Parse the form data.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Validate the form.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")

	// In case of error, redisplay the signup form along with a 422 status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	// Try to create a *new user* record in the database.
	// If *email* already exists, add an error message to the form & re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already registered.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}
	// If successful. Confirm with a *flash message* that their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Signup successful. Please log in.")

	// Add redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a form for logging in a user...")

}

// ============================================================================== #

// A helper method to decode form data:
// `dst`: target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination,
		// the `Decode()` method will return an error with the type *form.InvalidDecoderError.
		var invalidDecoderErr *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}
		// For all other errors, we return them as normal:
		return err
	}
	return nil
}
