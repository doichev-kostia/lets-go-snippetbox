package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.doichevkostia.dev/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	loglevel := flag.String("loglevel", "info", "Logger level")
	dsn := flag.String("dsn", "file:db.sqlite", "SQLite data source name")
	secure := flag.Bool("secure", true, "Use HTTPS server")
	cert := flag.String("cert", "./tls/cert.pem", "TLS certificate file")
	key := flag.String("key", "./tls/key.pem", "TLS key file")

	flag.Parse()

	if *secure == true {
		if _, err := os.Stat(*cert); os.IsNotExist(err) {
			log.Fatalf("Certificate file %s does not exist", *cert)
		}

		if _, err := os.Stat(*key); os.IsNotExist(err) {
			log.Fatalf("Key file %s does not exist", *key)
		}
	}

	var level slog.Level

	switch *loglevel {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db, PasswordCost: 12},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// some elliptic curves with assembly implementation. Idk what this is yet
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,

		IdleTimeout: time.Minute, // the keep-alive connection timeout
		// could use ReadHeaderTimeout and different read timeout for different handlers (body can have different size)
		ReadTimeout: 5 * time.Second,
		// writes made by a handler are buffered and written to the connection as one.
		// Idea is NOT to prevent the long-running handlers, but to prevent the data the handler returns  from taking too long to write
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", *addr)

	if *secure == true {
		err = srv.ListenAndServeTLS(*cert, *key)
	} else {
		err = srv.ListenAndServe()
	}
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
