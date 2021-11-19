package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ITA-Dnipro/Dp-218_Go/db"
	"github.com/ITA-Dnipro/Dp-218_Go/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var userIDKey = "userID"

func users(router chi.Router) {
	router.Get("/", getAllUsers)
	router.Post("/", createUser)
	router.Route("/{"+userIDKey+"}", func(r chi.Router) {
		r.Use(UserContext)
		r.Get("/", getUser)
		r.Put("/", updateUser)
		r.Delete("/", deleteUser)
	})
}

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, userIDKey)
		if userId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid user ID")))
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddUser(user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dbInstance.GetAllUsers()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, users); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	user, err := dbInstance.GetUserById(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	err := dbInstance.DeleteUser(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	render.Render(w, r, StatusOK)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	userData := models.User{}
	if err := render.Bind(r, &userData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	userData, err := dbInstance.UpdateUser(userId, userData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &userData); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
