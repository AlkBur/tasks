package handlers

import (
	"errors"
	"io"
	"net/http"
	"tasks/db"
	"tasks/lib"
	"tasks/server/auth"
	"tasks/service"
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterUser(s TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		var req db.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			lib.JSON(w, lib.Msg{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPw, err := lib.Hash(req.Password)
		if err != nil {
			if errors.Is(err, lib.ErrTooShort) {
				l.Error().Err(err).Msgf("The given password is too short%v", err)
				lib.JSON(w, lib.Msg{"error": "password is too short"}, http.StatusBadRequest)
				return
			}
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error during password hashing"}, http.StatusInternalServerError)
			return
		}

		uname, err := s.RegisterUser(ctx, &db.User{
			Username: req.Username,
			Password: hashedPw,
			Email:    req.Email,
		})

		switch {
		case errors.Is(err, service.ErrAlreadyExists):
			l.Error().Err(err).Msgf("registration failed, username or email already in use for user %s", req.Username)
			lib.JSON(w, lib.Msg{"error": "username or email already in use"}, http.StatusForbidden)
			return
		case errors.Is(err, service.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during User registration! %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error during user registration"}, http.StatusInternalServerError)
			return
		default:
			lib.JSON(w, lib.Msg{"success": "User registration successful!"}, http.StatusCreated)
			l.Info().Msgf("User registration for %s was successful!", uname)
		}
	}
}

func LoginUser(s TaskService, t auth.TokenManager, tokenDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		req := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		user, err := s.GetUser(ctx, req.Username)
		switch {
		case errors.Is(err, service.ErrNotFound):
			l.Error().Err(err).Msgf("user: %s is not found", req.Username)
			lib.JSON(w, lib.Msg{"error": "user is not found"}, http.StatusForbidden)
			return
		case errors.Is(err, service.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during user lookup! %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error during user lookup!"}, http.StatusInternalServerError)
			return
		}

		err = lib.Validate(user.Password, req.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", req.Username)
			lib.JSON(w, lib.Msg{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := t.CreateToken(req.Username, tokenDuration)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			lib.JSON(w, lib.Msg{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		lib.SetCookie(w, "paseto", token, payload.ExpiresAt)
		lib.JSON(w, lib.Msg{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", req.Username)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, _, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		uname, err := io.ReadAll(r.Body)
		if err != nil {
			lib.JSON(w, lib.Msg{"error": "couldn't decode request body"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
		})
		lib.JSON(w, lib.Msg{"success": "user successfully logged out"}, http.StatusOK)
		l.Info().Msgf("User logout for %s was successful!", string(uname))
	}
}
