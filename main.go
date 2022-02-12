package main

import (
	"accounts/scrapper"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func setHome(c echo.Context) error {
	return c.File("index.html")
}

func setMain(c echo.Context) error {
	defer os.Remove("jobs.csv")
	word := strings.ToLower(c.FormValue("word"))
	scrapper.Scrap(word)
	return c.Attachment("jobs.csv", "jobs.csv")
}

func main() {
	e := echo.New()
	e.GET("/", setHome)
	e.POST("/scrap", setMain)
	e.Logger.Fatal(e.Start(":1323"))
}
