package main

import (
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.doichevkostia.dev/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)

	mux.HandleFunc("/ping", ping)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.makeHandler(app.home)))

	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.makeHandler(app.snippetView)))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.makeHandler(app.userSignup)))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.makeHandler(app.userSignupPost)))

	mux.Handle("GET /user/login", dynamic.ThenFunc(app.makeHandler(app.userLogin)))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.makeHandler(app.userLoginPost)))

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.makeHandler(app.snippetCreate)))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.makeHandler(app.snippetCreatePost)))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.makeHandler(app.userLogoutPost)))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
