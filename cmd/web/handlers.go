package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
	"net/http"
	"snippetbox.doichevkostia.dev/internal/models"
	"snippetbox.doichevkostia.dev/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) error {
	snippets, err := app.snippets.Latest()
	if err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), "TEST", "TEST")

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.gohtml", data)
	return nil
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		return NewBadRequestError("invalid UUID", nil)
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
			return NewNotFoundError("No snippet with provided id", nil)
		} else {
			return err
		}
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.gohtml", data)
	return nil
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) error {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.gohtml", data)
	return nil
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) error {
	var formData snippetCreateForm
	err := app.decodePostForm(r, &formData)
	if err != nil {
		var decodeErrors form.DecodeErrors
		if errors.As(err, &decodeErrors) {
			return NewBadRequestError("invalid form", FormErrorsToFieldViolation(decodeErrors))
		} else {
			return err
		}
	}

	formData.CheckField(validator.NotBlank(formData.Title), "title", "This field can't be blank")
	formData.CheckField(validator.MaxChars(formData.Title, 100), "title", "This field cannot be more than 100 characters long")
	formData.CheckField(validator.NotBlank(formData.Content), "content", "This field can't be blank")
	formData.CheckField(validator.PermittedValue(formData.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !formData.Valid() {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, "create.gohtml", data)
		return nil
	}

	id, err := app.snippets.Insert(formData.Title, formData.Content, formData.Expires)

	if err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), "toast", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%s", id.String()), http.StatusSeeOther)
	return nil
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) error {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}

	app.render(w, r, http.StatusOK, "signup.gohtml", data)
	return nil
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) error {
	var formData userSignUpForm
	err := app.decodePostForm(r, &formData)

	if err != nil {
		var decodeErrors form.DecodeErrors
		if errors.As(err, &decodeErrors) {
			return NewBadRequestError("invalid form", FormErrorsToFieldViolation(decodeErrors))
		} else {
			return err
		}
	}

	formData.CheckField(validator.NotBlank(formData.Name), "name", "This field cannot be blank")
	formData.CheckField(validator.NotBlank(formData.Email), "email", "This field cannot be blank")
	formData.CheckField(validator.Matches(formData.Email, validator.EmailRX), "email", "This field must be a valid email address")
	formData.CheckField(validator.NotBlank(formData.Password), "password", "This field cannot be blank")
	formData.CheckField(validator.MinChars(formData.Password, 8), "password", "This field must be at least 8 characters long")

	if !formData.Valid() {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return nil
	}

	id, err := app.users.Insert(formData.Name, formData.Email, formData.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			formData.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = formData
			app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
			return nil
		} else {
			return err
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id.String())

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	return nil
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) error {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.gohtml", data)
	return nil
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) error {
	var formData userLoginForm

	err := app.decodePostForm(r, &formData)
	if err != nil {
		var decodeErrors form.DecodeErrors
		if errors.As(err, &decodeErrors) {
			return NewBadRequestError("invalid form", FormErrorsToFieldViolation(decodeErrors))
		} else {
			return err
		}
	}

	formData.CheckField(validator.NotBlank(formData.Email), "email", "This field cannot be blank")
	formData.CheckField(validator.Matches(formData.Email, validator.EmailRX), "email", "This field must be a valid email address")
	formData.CheckField(validator.NotBlank(formData.Password), "password", "This field cannot be blank")

	if !formData.Valid() {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		return nil
	}

	id, err := app.users.Authenticate(formData.Email, formData.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			formData.AddGeneralError("Invalid credentials")
			data := app.newTemplateData(r)
			data.Form = formData
			app.render(w, r, http.StatusUnauthorized, "login.gohtml", data)
			return nil
		} else {
			return err
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id.String())

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	return nil
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) error {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "toast", "Successful logout")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
