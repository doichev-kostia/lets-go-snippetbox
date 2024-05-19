package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (app *application) makeHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		var apiError ApiError
		if errors.As(err, &apiError) {
			writeJSON(w, StatusCodeMap[apiError.Code], apiError)
		} else {
			internalError := NewApiError(ErrorInternal, fmt.Errorf("internal server error"), nil)
			writeJSON(w, StatusCodeMap[internalError.Code], internalError)
		}
		app.logger.Error("HTTP API error", "msg", err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	tmpl, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

// error can be DecodeErrors from form decoder
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodedErr *form.InvalidDecoderError
		if errors.As(err, &invalidDecodedErr) {
			panic(err)
		}
		return err
	}

	return nil
}

func FormErrorsToFieldViolation(errors form.DecodeErrors) []FieldViolation {
	violations := make([]FieldViolation, len(errors))
	for name, error := range errors {
		violations = append(violations, FieldViolation{
			Field:       name,
			Description: error.Error(),
		})
	}

	return violations
}
