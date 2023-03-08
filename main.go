package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template //mengirim request
}

// Function Render = pada function ini mendeklarasikan alias t dengan template(dari package Template echo)
// Function Render memilik 4 parameter , w yaitu alias io.Writer atau mengirim data ke io.Writer
// w itu status, name itu nama file , data itu data interface ( data yg dikirimkan oleh html)
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Static("/public", "public")

	t := &Template{ // & menerima request
		templates: template.Must(template.ParseGlob("views/*.html")), //method Must
	}

	e.Renderer = t

	// Routing
	e.GET("/hello", helloWorld)       //to execute //localhost:5000/hello
	e.GET("/", home)                  //localhost:5000
	e.GET("/contact", contact)        //localhost:5000/contact
	e.GET("/myProject", myproject)    //localhost:5000/myProject
	e.POST("/myProject", postproject) //localhost:5000/myProject

	fmt.Println("Server Berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!") //http merupakan package dari golang yg akan mengirimkan StatusOK = 200
}

// Function Home
func home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// Function Contact
func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

// Function My Project
func myproject(c echo.Context) error {
	return c.Render(http.StatusOK, "myProject.html", nil)
}

func postproject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	description := c.FormValue("description")
	// tech1 := c.FormValue("tech1")
	// tech2 := c.FormValue("tech2")
	// tech3 := c.FormValue("tech3")
	// tech4 := c.FormValue("tech4")

	println("Project Name : " + projectName)
	println("Start Date : " + startDate)
	println("End Date : " + endDate)
	println("Description : " + description)
	// println("Technology : " + tech1)
	// println("Technology : " + tech2)
	// println("Technology : " + tech3)
	// println("Technology : " + tech4)

	return c.Redirect(http.StatusMovedPermanently, "/")
}
