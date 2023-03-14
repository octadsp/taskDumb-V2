package main

import (
	"chapter2/connection"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template //mengirim request
}

type Project struct {
	ID              int
	Projectname     string
	StartDate       time.Time
	EndDate         time.Time
	Description     string
	Technology      []string
	Duration        string
	Image           string
	FormatStartDate string
	FormatEndDate   string
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

	// Routing              //to execute //localhost:5000/hello
	e.GET("/", home)                           //localhost:5000
	e.GET("/contact", contact)                 //localhost:5000/contact
	e.GET("/myProject", myproject)             //localhost:5000/myProject
	e.GET("/detailProject/:id", detailproject) //localhost:5000/detailProject/:id
	e.POST("/myProject", postproject)          //localhost:5000/myProject
	e.GET("/deleteProject/:id", deleteproject)
	// e.GET("/editProject/:id", editproject) //localhost:5000/deleteProject
	// e.POST("/editProject/:id", postEditProject) //localhost:5000/deleteProject

	fmt.Println("Server Berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

// Function Home
func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, description, technology FROM tb_project")

	var result []Project

	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Projectname, &each.StartDate, &each.EndDate, &each.Description, &each.Technology)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		// Logic ngitung seisih
		// >= 1 tahun = "x years"
		// < 1 tahun = "x months"
		// < 1 bulan = "x weeks"

		// contohSelisih = "bukanasd"

		layout := "2006-01-02"

		startDate := each.StartDate.Format(layout)
		endDate := each.EndDate.Format(layout)

		t1, _ := time.Parse(layout, endDate)
		t2, _ := time.Parse(layout, startDate)

		diff := t1.Sub(t2)

		days := int(diff.Hours() / 24)
		weeks := int(diff.Hours() / 24 / 7)
		months := int(diff.Hours() / 24 / 30)
		years := int(diff.Hours() / 24 / 365)

		var Duration string
		if years >= 1 {

			//Fungsi buat tahun lebih dari 1
			if years > 1 {
				Duration = strconv.Itoa(years) + " years"
			} else {
				//Fungsi buat tahun == 1
				Duration = strconv.Itoa(years) + " year"
			}
		} else if months < 12 && months > 0 {

			//Fungsi buat handling lebih dari 1 bulan
			if months > 1 {
				Duration = strconv.Itoa(months) + " months"
			} else {
				// Fungsi buat handling 1 bulan doang
				Duration = strconv.Itoa(months) + " month"
			}
		} else if weeks > 0 && weeks <= 4 {
			//Fungsi buat handling kurang dari 1 bulan
			if weeks > 1 {
				Duration = strconv.Itoa(weeks) + " weeks"
			} else {
				// Fungsi buat handling 1 week doang
				Duration = strconv.Itoa(weeks) + " week"
			}
		} else {
			//Fungsi buat handling kurang dari 7 hari
			if days > 1 && days < 7 {
				Duration = strconv.Itoa(days) + " days"

				//Fungsi buat handling 1 hari
			} else if days == 1 {
				Duration = strconv.Itoa(days) + " day"

				//Fungsi buat handling < 1 hari
			} else {
				Duration = "less than a day"
			}
		}

		var newEach = Project{
			ID:          each.ID,
			Projectname: each.Projectname,
			Description: each.Description,
			Technology:  each.Technology,
			Duration:    Duration,
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
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	description := c.FormValue("description")
	techValues := c.Request().Form["techIcon"]
	image := "project.png"

	technology := []string{}
	for _, val := range techValues {
		technology = append(technology, val)
	}

	fmt.Println(technology)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (project_name, description, image, start_date, end_date, technology) VALUES ($1, $2, $3, $4, $5, $6)", projectName, description, image, startDate, endDate, technology)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, description, image, technology, start_date, end_date FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.ID, &ProjectDetail.Projectname, &ProjectDetail.Description, &ProjectDetail.Image, &ProjectDetail.Technology, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("02 Jan 2006")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("02 Jan 2006")

	layout := "2006-01-02"

	startDate := ProjectDetail.StartDate.Format(layout)
	endDate := ProjectDetail.EndDate.Format(layout)

	t1, _ := time.Parse(layout, endDate)
	t2, _ := time.Parse(layout, startDate)

	diff := t1.Sub(t2)

	days := int(diff.Hours() / 24)
	weeks := int(diff.Hours() / 24 / 7)
	months := int(diff.Hours() / 24 / 30)
	years := int(diff.Hours() / 24 / 365)

	var Duration string
	if years >= 1 {

		//Fungsi buat tahun lebih dari 1
		if years > 1 {
			Duration = strconv.Itoa(years) + " years"
		} else {
			//Fungsi buat tahun == 1
			Duration = strconv.Itoa(years) + " year"
		}
	} else if months < 12 && months > 0 {

		//Fungsi buat handling lebih dari 1 bulan
		if months > 1 {
			Duration = strconv.Itoa(months) + " months"
		} else {
			// Fungsi buat handling 1 bulan doang
			Duration = strconv.Itoa(months) + " month"
		}
	} else if weeks > 0 && weeks <= 4 {
		//Fungsi buat handling kurang dari 1 bulan
		if weeks > 1 {
			Duration = strconv.Itoa(weeks) + " weeks"
		} else {
			// Fungsi buat handling 1 week doang
			Duration = strconv.Itoa(weeks) + " week"
		}
	} else {
		//Fungsi buat handling kurang dari 7 hari
		if days > 1 && days < 7 {
			Duration = strconv.Itoa(days) + " days"

			//Fungsi buat handling 1 hari
		} else if days == 1 {
			Duration = strconv.Itoa(days) + " day"

			//Fungsi buat handling < 1 hari
		} else {
			Duration = "less than a day"
		}
	}

	// technology := []string{}
	// for _, val := range ProjectDetail.Technology {
	// 	technology = append(technology, val)
	// }

	// ProjectDetail.Technology = technology
	ProjectDetail.Duration = Duration

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	return c.Render(http.StatusOK, "detailProject.html", data)
}

func deleteproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// func editproject(c echo.Context) error {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	var ProjectEdit = Project{}

// 	for i, data := range dataProject {
// 		if id == i {
// 			ProjectEdit = Project{
// 				Projectname: data.Projectname,
// 				Description: data.Description,
// 			}
// 		}
// 	}

// 	editProject := map[string]interface{}{
// 		"Project": ProjectEdit,
// 	}

// 	return c.Render(http.StatusOK, "editProject.html", editProject)
// }

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
