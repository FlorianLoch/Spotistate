package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	constants "github.com/florianloch/cassette/internal"
	"github.com/florianloch/cassette/internal/handler"
	"github.com/florianloch/cassette/internal/middleware"
	"github.com/florianloch/cassette/internal/persistence"
	"github.com/florianloch/cassette/internal/util"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	redirectURL          *url.URL
	auth                 *spotify.Authenticator
	store                *sessions.CookieStore
	dao                  *persistence.PlayerStatesDAO
	isDevMode            bool
	webStaticContentPath = "/web/dist"
)

type m map[string]interface{}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	isDevMode = util.Env(constants.EnvENV, "") == "DEV"
	if isDevMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Running in DEV mode. Being more verbose. Set environment variable 'ENV' to 'DEV' to activate.")
	}

	var networkInterface = util.Env(constants.EnvNetworkInterface, constants.DefaultNetworkInterface)
	// We also have to check for "PORT" as that is how Heroku/Dokku etc. tells the app where to listen
	var port = util.Env(constants.EnvPort, os.Getenv("PORT"))
	var appURL = util.Env(constants.EnvAppURL, "http://"+networkInterface+":"+port+"/")

	var secret32Bytes, err = util.Make32ByteSecret(util.Env(constants.EnvSecret, ""))
	if err != nil {
		log.Fatal().Err(err).Msg("Could not generate secret. Aborting.")
	}

	var mongoDBURI = util.Env(constants.EnvMongoURI, "")
	if mongoDBURI == "" {
		log.Fatal().Msg("No URI for connecting to MongoDB given. Aborting.")
	}
	dao, err = persistence.Connect(mongoDBURI)
	if err != nil {
		log.Fatal().Err(err).Str("mongoDBURI", mongoDBURI).Msg("Failed connecting to MongoDB.")
	}

	store = sessions.NewCookieStore(secret32Bytes)

	redirectURL, err = url.Parse(appURL)
	if err != nil {
		log.Fatal().Err(err).Str("appURL", appURL).Msgf("'%s' variable is not set to a valid value.", constants.EnvAppURL)
	}
	redirectURL.Path = "/spotify-oauth-callback"
	tmp := spotify.NewAuthenticator(redirectURL.String(), spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	auth = &tmp

	gob.Register(&spotify.PrivateUser{})
	gob.Register(&oauth2.Token{})
	gob.Register(&m{})

	var clientID = util.Env(constants.EnvSpotifyClientID, "")
	var clientSecret = util.Env(constants.EnvSpotifyClientSecret, "")

	if clientID == "" || clientSecret == "" {
		log.Fatal().Msgf("Please make sure '%s' and '%s' are set. Aborting.", constants.EnvSpotifyClientID, constants.EnvSpotifyClientSecret)
	}

	auth.SetAuthInfo(clientID, clientSecret)

	var csrfMiddleware = csrf.Protect(
		secret32Bytes,
		csrf.RequestHeader(constants.CSRFHeaderName),
		csrf.CookieName(constants.CSRFCookieName),
		csrf.Secure(!isDevMode),
		csrf.MaxAge(60*60*24*365), // Cookie is valid for 1 year
		csrf.ErrorHandler(csrfErrorHandler{}),
	)

	var cwd, _ = os.Getwd()
	var staticAssetsPath = cwd + webStaticContentPath
	var spaHandler = handler.NewSpaHandler(staticAssetsPath, "index.html")
	log.Info().Msgf("Loading assets from: '%s'", staticAssetsPath)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		if isDevMode {
			r.Use(debugLogger)
		}
		r.Use(middleware.CreateConsentMiddleware(spaHandler))
		r.Use(middleware.CreateSpotifyAuthMiddleware(store, auth, redirectURL))
		r.Use(csrfMiddleware)
		r.Use(attachStore)

		// apiRouter.HandleFunc("/spotify-oauth-callback", func(w http.ResponseWriter, r *http.Request) {})

		r.Head("/csrfToken", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(constants.CSRFHeaderName, csrf.Token(r))
		})

		r.With(attachAuth).Get("/activeDevices", handler.ActiveDevicesHandler)
		r.With(attachAuth).With(attachDAO).Route("/playerStates", func(r chi.Router) {
			r.Post("/", handler.StorePostHandler)
			r.Get("/", handler.StoreGetHandler)
			r.With(attachSlot).Route("/{slot}", func(r chi.Router) {
				r.Put("/", handler.StorePostHandler)
				r.Delete("/", handler.StoreDeleteHandler)
				r.Post("/restore", handler.RestoreHandler)
			})
		})

		r.With(attachDAO).Route("/you", func(r chi.Router) {
			r.Get("/", handler.UserExportHandler)
			r.Delete("/", handler.UserDeleteHandler)
		})

	})

	// Provide the webapp following the SPA pattern: all non-API routes not being able
	// to be resolved within the assets directory will return the webapp entry point.
	r.NotFound(spaHandler.ServeHTTP)

	var interfacePort = networkInterface + ":" + port
	log.Info().Msgf("Webserver started on %s", interfacePort)

	err = http.ListenAndServe(interfacePort, r)
	log.Fatal().Err(err).Msg("Server terminated.")
}

func attachStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "store", store)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func attachAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "auth", auth)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func attachDAO(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "dao", dao)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func attachSlot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := checkSlotParameter(r)
		if err != nil {
			log.Debug().Err(err).Msg("Could not retrieve slot from request.")
			http.Error(w, fmt.Sprintf("Could not process request. Please make sure the given slot is valid: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "slot", slot)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkSlotParameter(r *http.Request) (int, error) {
	var slotStr = chi.URLParam(r, "slot")

	if slotStr == "" {
		return -1, errors.New("query parameter 'slot' not found")
	}

	var slot, err = strconv.Atoi(slotStr)
	if err != nil {
		return -1, errors.New("query parameter 'slot' is not a valid integer")
	}
	if slot < 0 {
		return -1, errors.New("query parameter 'slot' has to be >= 0")
	}

	return slot, nil
}

func debugLogger(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(constants.CSRFHeaderName)

		log.Debug().Str("csrfToken", token).Stringer("url", r.URL).Interface("cookies", r.Cookies()).Msg("")

		nxt.ServeHTTP(w, r)
	})
}

type csrfErrorHandler struct{}

func (csrfErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	failureReason := csrf.FailureReason(r)
	csrfToken := r.Header.Get(constants.CSRFHeaderName)
	csrfCookie, err := r.Cookie(constants.CSRFCookieName)
	csrfCookieContent := "ERROR: cookie not present."
	if err == nil {
		csrfCookieContent = csrfCookie.Value
	}

	http.Error(w,
		fmt.Sprintf("Failed verifying CSRF token. Supplied token: '%s'; Error: '%s'. Expect token to be contained in header '%s'. Cookie named '%s' contains '%s'",
			csrfToken, failureReason, constants.CSRFHeaderName, constants.CSRFCookieName, csrfCookieContent),
		http.StatusUnauthorized)
}
