package main

import (
	"context"
	"fmt"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"

	"time"

	"github.com/labstack/echo/v4"
)

// Create struct -> struct is like class in javascript
type Blog struct {
	Id int
	ProjectName string
	StartDate time.Time
	EndDate time.Time
	Duration string
	Author string
	Description string
	Technologies []string
	Javascript bool
	PHP bool
	Java bool
	ReactJS bool
	Image string
}

func main() {

	e := echo.New()

	connection.DatabaseConnect()

	e.Static("/public", "public")
    // e.GET("/", func(c echo.Context) error {
    //     return c.String(http.StatusOK, "Hello, Saya Ahmad")
    // })
    // e.Logger.Fatal(e.Start("localhost:700"))
	
	e.GET("/", index)
	e.GET("/blog", blog)
	e.GET("/testimonial", testimonial)
	e.GET("/contact", contact)
	e.GET("/detail-blog/:id", detailBlog)
	e.GET("/update-blog-form/:id", updateBlogForm)


	//POST routing
	e.POST("/blog", addBlog)
	e.POST("/delete-blog/:id", deleteBlog)
	e.POST("/update-blog/:id", updateBlog)




	e.Logger.Fatal(e.Start(":5000"))

}

	func getDuration(startDate, endDate time.Time) string {


	durationTime := int(endDate.Sub(startDate).Hours())
	durationDays := durationTime / 24
	durationWeeks := durationDays / 7
	durationMonths := durationWeeks / 4
	durationYears := durationMonths / 12

	var duration string
	
	if durationYears > 1 {
		duration = strconv.Itoa(durationYears) + " years"
	} else if durationYears > 0 {
		duration = strconv.Itoa(durationYears) + " year"
	} else {
		if durationMonths > 1 {
			duration = strconv.Itoa(durationMonths) + " months"
		} else if durationMonths > 0 {
			duration = strconv.Itoa(durationMonths) + " month"
		} else {
			if durationWeeks > 1 {
				duration = strconv.Itoa(durationWeeks) + " weeks"
			} else if durationWeeks > 0 {
				duration = strconv.Itoa(durationWeeks) + " week"
			} else {
				if durationDays > 1 {
					duration = strconv.Itoa(durationDays) + " days"
				} else {
					duration = strconv.Itoa(durationDays) + " day"
				}
			}
		}
	}

	return duration

}

func index(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// GET DATA FROM DATABASE
	dataBlog, errBlog := connection.Conn.Query(context.Background(), "SELECT id, project_name, description, image, start_date, end_date, technologies FROM public.tb_blog")

	if errBlog != nil {
		return c.JSON(http.StatusInternalServerError, errBlog.Error())
	}

	var resultBlogs []Blog
	for dataBlog.Next() {
		var each = Blog{}

		err := dataBlog.Scan(&each.Id, &each.ProjectName, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Technologies)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Duration
		each.Duration = getDuration(each.StartDate, each.EndDate)

		//CHECKBOX
		if checkValue(each.Technologies, "javascript") {
			each.Javascript = true
		}
		if checkValue(each.Technologies, "php") {
			each.PHP = true
		}
		if checkValue(each.Technologies, "java") {
			each.Java = true
		}
		if checkValue(each.Technologies, "reactJS") {
			each.ReactJS = true
		}

		resultBlogs = append(resultBlogs, each)
	}

	blog := map[string]interface{}{
		"Blogs": resultBlogs,
	}

	return tmpl.Execute(c.Response(), blog)
}

	
	func blog (c echo.Context) error  {
		tmpl, err := template.ParseFiles("views/blog.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return tmpl.Execute(c.Response(), nil)
	}

	func testimonial (c echo.Context) error  {
		tmpl, err := template.ParseFiles("views/testimonial.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return tmpl.Execute(c.Response(), nil)
	}

	func contact (c echo.Context) error  {
		tmpl, err := template.ParseFiles("views/contact.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return tmpl.Execute(c.Response(), nil)
	}

	func detailBlog (c echo.Context) error {
		id := c.Param("id")

		var blogDetail = Blog{}

		//query get 1 data
		idToInt, _ := strconv.Atoi(id)

		errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, description, image, start_date, end_date, technologies FROM public.tb_blog WHERE id = $1", idToInt).Scan(&blogDetail.Id, &blogDetail.ProjectName, &blogDetail.Description, &blogDetail.Image,  &blogDetail.StartDate,  &blogDetail.EndDate ,&blogDetail.Technologies)


		if errQuery != nil {
			return c.JSON(http.StatusInternalServerError, errQuery.Error())
		}
		
		// Duration
		blogDetail.Duration = getDuration(blogDetail.StartDate, blogDetail.EndDate)


		//CHECKBOX
		if checkValue(blogDetail.Technologies, "javascript") {
			blogDetail.Javascript = true
		}
		if checkValue(blogDetail.Technologies, "php") {
			blogDetail.PHP = true
		}
		if checkValue(blogDetail.Technologies, "java") {
			blogDetail.Java = true
		}
		if checkValue(blogDetail.Technologies, "reactJS") {
			blogDetail.ReactJS = true
		}


	data := map[string]interface{}{
		"Id" : id,
		"Blog": blogDetail,
		"startDateString" 	: blogDetail.StartDate.Format("2006-01-02"),
		"endDateString"		: blogDetail.EndDate.Format("2006-01-02"),
	}

	var tmpl, err = template.ParseFiles("views/detail-blog.html")
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)

	}



	func addBlog(c echo.Context) error {
		projectName := c.FormValue("input-projectName")
		startDate := c.FormValue("input-startDate")
		endDate := c.FormValue("input-endDate")
		description := c.FormValue("input-description")
		javascript := c.FormValue("input-javascript")
		php := c.FormValue("input-php")
		java := c.FormValue("input-java")
		react := c.FormValue("input-reactJS")
		technologies := []string{javascript,php,java,react}
		image := c.FormValue("input-image")
	
		

		insertDb, err := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_blog (project_name, description, image, start_date, end_date, technologies) VALUES ($1, $2, $3, $4, $5, $6)", projectName, description, image, startDate, endDate, technologies)
			
		fmt.Println("Row Affected : ", insertDb.RowsAffected()  )

		if err != nil {
			fmt.Println("There is something error guys")
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	func deleteBlog(c echo.Context) error {
		id := c.Param("id")
		idToInt, _ := strconv.Atoi(id)
		connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", idToInt)

		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	
	func updateBlogForm (c echo.Context) error {
		id := c.Param("id")

		var blogDetail = Blog{}

		//query get 1 data
		idToInt, _ := strconv.Atoi(id)

		errQuery := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, description, image, start_date, end_date, technologies FROM public.tb_blog WHERE id = $1", idToInt).Scan(&blogDetail.Id, &blogDetail.ProjectName, &blogDetail.Description, &blogDetail.Image,  &blogDetail.StartDate,  &blogDetail.EndDate ,&blogDetail.Technologies)


		if errQuery != nil {
			return c.JSON(http.StatusInternalServerError, errQuery.Error())
		}
		
		// Duration
		blogDetail.Duration = getDuration(blogDetail.StartDate, blogDetail.EndDate)

		//CHECKBOX
		if checkValue(blogDetail.Technologies, "javascript") {
			blogDetail.Javascript = true
		}
		if checkValue(blogDetail.Technologies, "php") {
			blogDetail.PHP = true
		}
		if checkValue(blogDetail.Technologies, "java") {
			blogDetail.Java = true
		}
		if checkValue(blogDetail.Technologies, "reactJS") {
			blogDetail.ReactJS = true
		}


	data := map[string]interface{}{
		"Id" : id,
		"Blog": blogDetail,
	}

	var tmpl, err = template.ParseFiles("views/update-blog.html")
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)

	}
	
	func updateBlog (c echo.Context) error {
		// Menangkap Id dari Query Params
		id, _:= strconv.Atoi(c.Param("id"))
		
		projectName := c.FormValue("input-projectName")
		startDate := c.FormValue("input-startDate")
		endDate := c.FormValue("input-endDate")
		description := c.FormValue("input-description")
		javascript := c.FormValue("input-javascript")
		php := c.FormValue("input-php")
		java := c.FormValue("input-java")
		react := c.FormValue("input-reactJS")
		technologies := []string{javascript,php,java,react}
		image := c.FormValue("input-image")
	
		

		insertDb, err := connection.Conn.Exec(context.Background(), "UPDATE public.tb_blog SET project_name=$1, description=$2, image=$3, start_date=$4, end_date=$5, technologies=$6 WHERE id=$7", projectName, description, image, startDate, endDate, technologies, id,)
			
		fmt.Println("Row Affected : ", insertDb.RowsAffected()  )

		if err != nil {
			fmt.Println("There is something error guys")
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		
		return c.Redirect(http.StatusMovedPermanently, "/")
	}


	func checkValue(slice []string, object string) bool {
		for _, data := range slice {
			if data == object {
				return true
			}
		}
		return false
	}