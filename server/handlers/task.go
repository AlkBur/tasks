package handlers

import (
	"errors"
	"net/http"
	"strings"
	"tasks/db"
	"tasks/lib"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func CreateTask(s TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		var taskRequest db.Task

		err := json.NewDecoder(r.Body).Decode(&taskRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into lib.JSON during registration. %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&taskRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during Note struct validation %v", err)
			lib.JSON(w, lib.Msg{"error": "wrongly formatted or missing Note parameter"}, http.StatusBadRequest)
			return
		}

		retID, err := s.CreateTask(ctx, taskRequest.Title, taskRequest.User, taskRequest.Text)
		switch {
		case errors.Is(err, db.ErrTaskAlreadyExists):
			l.Error().Err(err).Msgf("Task creation failed, a task with that ID already exists")
			lib.JSON(w, lib.Msg{"error": "a Task with that id already exists! ID must be unique."}, http.StatusForbidden)
			return
		default:
			l.Info().Msgf("Task with ID %v has been created for user: %s", retID, taskRequest.User)
			lib.JSON(w, lib.Msg{"success": "note creation successful!"}, http.StatusCreated)
		}
	}
}

func GetAllTasksFromUser(s TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		username := r.URL.Query().Get("username")
		if username == "" {
			l.Error().Msgf("error fetching username, the request parameter is empty. %s", username)
			lib.JSON(w, lib.Msg{"error": "user not in request params"}, http.StatusBadRequest)
			return
		}

		notes, err := s.GetAllTasksFromUser(ctx, username)
		switch {
		case errors.Is(err, db.ErrTaskNotFound):
			l.Info().Msgf("Requested user has no Task!. %s", username)
		default:
			l.Info().Msgf("Retriving user task for %s was successful!", username)
			lib.JSON(w, notes, http.StatusOK)
		}
	}
}

func DeleteTask(s TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Info().Msgf("Could not convert ID to UUID.")
			lib.JSON(w, lib.Msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		id, err := s.DeleteTask(ctx, reqUUID)
		switch {
		case errors.Is(err, db.ErrNoRows):
			l.Info().Msg("User has no tasks to delete from!")
		default:
			l.Info().Msgf("Deleting task %v was successful!", id)
			lib.JSON(w, lib.Msg{"success": "task deleted"}, http.StatusOK)
			return
		}
	}
}

func UpdateTask(s TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := lib.SetupHandler(w, r.Context())
		defer cancel()

		var isTextValid bool = true

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Error().Err(err).Msgf("Could not convert ID to UUID.")
			lib.JSON(w, lib.Msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		updateRequest := struct {
			ID    uuid.UUID `json:"id"`
			Title string    `json:"title" validate:"required,min=4"`
			Text  string    `json:"text"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into lib.JSON during registration. %v", err)
			lib.JSON(w, lib.Msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}
		validate := validator.New()
		err = validate.Struct(&updateRequest)
		if err != nil {
			l.Error().Err(err).Msgf("title must be more than 4 characters long!")
			lib.JSON(w, lib.Msg{"error": "title of a note must be more than 4 characters long!"}, http.StatusBadRequest)
			return
		}

		if updateRequest.Text == "" {
			isTextValid = false
		}

		id, err := s.UpdateTask(ctx, reqUUID, updateRequest.Title, updateRequest.Text, isTextValid)
		switch {
		case err != nil:
			l.Info().Err(err).Msgf("Could not update Note %v", reqUUID)
			lib.JSON(w, lib.Msg{"error": "could not update note"}, http.StatusInternalServerError)
			return
		default:
			l.Info().Msgf("Updating note %v was successful!", id)
			lib.JSON(w, lib.Msg{"success": "note deleted"}, http.StatusOK)
			return
		}
	}
}
