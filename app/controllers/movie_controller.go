package controllers

import (
	"revel-dynamodb-v1/app"

	"github.com/revel/revel"
)

type MovieController struct {
	*revel.Controller
}

func (c MovieController) GetAllMovies() revel.Result {
	movies, err := app.Service.GetAllMovies()
	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJSON(movies)
}
