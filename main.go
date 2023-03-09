package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template //mengirim request
}

type Project struct {
	Projectname string
	Duration    string
	Description string
}

var dataProject = []Project{
	{
		Projectname: "Ini Project TESTING",
		Duration:    "5 Bulan",
		Description: "Ini Project TESTING",
	},
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
	e.GET("/hello", helloWorld)                //to execute //localhost:5000/hello
	e.GET("/", home)                           //localhost:5000
	e.GET("/contact", contact)                 //localhost:5000/contact
	e.GET("/myProject", myproject)             //localhost:5000/myProject
	e.GET("/detailProject/:id", detailproject) //localhost:5000/detailProject/:id
	e.POST("/myProject", postproject)          //localhost:5000/myProject
	e.GET("/deleteProject/:id", deleteproject)
	e.GET("/editProject/:id", editproject) //localhost:5000/deleteProject
	// e.POST("/editProject/:id", postEditProject) //localhost:5000/deleteProject

	fmt.Println("Server Berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!") //http merupakan package dari golang yg akan mengirimkan StatusOK = 200
}

// Function Home
func home(c echo.Context) error {
	projects := map[string]interface{}{
		"Projects": dataProject,
	}

	return c.Render(http.StatusOK, "index.html", projects)
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
	// startDate := c.FormValue("startDate")
	// endDate := c.FormValue("endDate")
	description := c.FormValue("description")

	var newProject = Project{
		Projectname: projectName,
		Duration:    "9 Bulan Ges",
		Description: description,
	}

	dataProject = append(dataProject, newProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Project{
				Projectname: data.Projectname,
				Duration:    data.Duration,
				Description: data.Description,
			}
		}
	}

	//data yang akan dikirimkan ke html menggunakan map interface
	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
	}

	return c.Render(http.StatusOK, "detailProject.html", detailProject)
}

func deleteproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectEdit = Project{}

	for i, data := range dataProject {
		if id == i {
			ProjectEdit = Project{
				Projectname: data.Projectname,
				Duration:    data.Duration,
				Description: data.Description,
			}
		}
	}

	editProject := map[string]interface{}{
		"Project": ProjectEdit,
	}

	return c.Render(http.StatusOK, "editProject.html", editProject)
}

// func postEditProject(c echo.Context) error {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	projectNameEdit := c.FormValue("projectNameEdit")
// 	descriptionEdit := c.FormValue("descriptionEdit")

// 	newEditedProject := Project{
// 		Projectname: projectNameEdit,
// 		Duration:    "9 Bulan Ges",
// 		Description: descriptionEdit,
// 	}

// 	dataProject[id] = append(dataProject, newEditedProject )

// 	return c.Redirect(http.StatusMovedPermanently, "/")
// }