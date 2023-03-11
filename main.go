package main

import (
	"chapter2/connection"
	"context"
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
	ID          int
	Projectname string
	// StartDate   time.Time
	// EndDate     string
	Description string
	Technology  []string
	// Image       string
	Selisih string
}

var dataProject = []Project{
	{
		ID:          0,
		Projectname: "Ini Project TESTING",
		// StartDate:   "2023-03-01",
		// EndDate:     "2023-03-10",
		Description: "Ini Project TESTING",
		// Technology:  map[string]string{"js": "/public/img/technology/js.png", "go": "/public/img/technology/go.png", "python": "/public/img/technology/python.png", "figma": "/public/img/technology/figma.png"},
		// Image:       "project1.png",
	},
	{
		ID:          1,
		Projectname: "Ini Project TESTING 2",
		// StartDate:   "2023-03-01",
		// EndDate:     "2023-03-10",
		Description: "Ini Project TESTING 2",
		// Technology:  map[string]string{"js": "/public/img/technology/js.png", "go": "/public/img/technology/go.png", "python": "/public/img/technology/python.png", "figma": "/public/img/technology/figma.png"},
		// Image:       "project1.png",
	},
}

// Function Render = pada function ini mendeklarasikan alias t dengan template(dari package Template echo)
// Function Render memilik 4 parameter , w yaitu alias io.Writer atau mengirim data ke io.Writer
// w itu status, name itu nama file , data itu data interface ( data yg dikirimkan oleh html)
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
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
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, description, technology FROM tb_project")

	var result []Project

	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Projectname, &each.Description, &each.Technology)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		// result = append(result, each)

		var contohSelisih = ""
		// Logic ngitung seisih
		// >= 1 tahun = "x years"
		// < 1 tahun = "x months"
		// < 1 bulan = "x weeks"

		contohSelisih = "bukanasd"

		var newEach = Project{
			ID:          each.ID,
			Projectname: each.Projectname,
			Description: each.Description,
			Technology:  each.Technology,
			Selisih:     contohSelisih,
		}
		result = append(result, newEach)
	}

	projects := map[string]interface{}{
		"Projects": result,
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
	// technology := c.

	var newProject = Project{
		ID:          0,
		Projectname: projectName,
		// StartDate:   startDate,
		// EndDate:     endDate,
		Description: description,
		// Technology:  technology,
		// Image:       "",
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
