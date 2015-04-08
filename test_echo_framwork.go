package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"github.com/thoas/stats"
	"net/http"
	"fmt"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users map[string]user

func init() {
	users = map[string]user {
		"1":user{
			ID: "1",
			Name:"Wreck-It Ralph",
		}
	}
}

func createUser(c *echo.Context) {
	u := new(user)
	if err := c.Bind(u); err == nil {
		users[u.ID] = *u;
		if err := c.JSON(http.StatusCreated, u); err == nil {
			fmt.Println("Hello")
		}
		return
	}
}

func getUsers(c *echo.Context) {
	c.JSON(http.StatusOK, users)
}

func getUser(c *echo.Context) {
	c.JSON(http.StatusOK, users[c.P(0)])
}

func main() {
	e := echo.New()

	e.Use(mw.Logger)
	e.Use(cors.Default().Handler)


}
