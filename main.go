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
	StartDate string
	EndDate string
	Duration string
	Author string
	Description string
	Javascript bool
	PHP bool
	Java bool
	ReactJS bool
	Image string
}

// Dummy data
var dataBlogs = []Blog {
	{
		ProjectName: "Dumbways ",
		Duration: "4 Bulan 10 Hari",
		Author: "Ahmad Syahran Zidane",
		Description: "Halo Guys",
		Javascript: true,
		PHP: true,
		Java: true,
		ReactJS: true,
		Image: "result1.jpg",
	},
	{
		ProjectName: "Dumbways ",
		Duration: "4 Bulan 10 Hari",
		Author: "Ahmad Syahran Zidane",
		Description: "Halo Guys",
		Javascript: true,
		PHP: true,
		Java: true,
		ReactJS:true,
		Image: "result2.jpg",
	},
	{
		ProjectName: "Dumbways ",
		Duration: "4 Bulan 10 Hari",
		Description: "Halo Guys",
		Javascript: true,
		PHP: true,
		Java: true,
		ReactJS:true,
		Image: "result3.jpg",
	},

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

	//POST routing
	e.POST("/blog", addBlog)
	e.POST("/delete-blog/:id", deleteBlog)



	e.Logger.Fatal(e.Start(":5000"))

}

	func getDuration(startDate string, endDate string) string {
	startTime, _ := time.Parse("2006-01-02", startDate) 
	endTime, _ := time.Parse("2006-01-02", endDate) 

	durationTime := int(endTime.Sub(startTime).Hours())
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

	func index (c echo.Context) error  {		
		tmpl, err := template.ParseFiles("views/index.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// GET DATA FROM DATABASE
		dataBlog, errBlog := connection.Conn.Query(context.Background(), "SELECT project_name, author, description, javascript, php, java, react, image, duration FROM tb_blog")

		if errBlog != nil {
			return c.JSON(http.StatusInternalServerError, errBlog.Error())
		}

		var resultBlogs []Blog
		for dataBlog.Next() {
			var each = Blog{}

			err := dataBlog.Scan(&each.ProjectName, &each.Author, &each.Description, &each.Javascript, &each.PHP, &each.Java, &each.ReactJS, &each.Image, &each.Duration)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			each.Author = "Syahran"


			
			resultBlogs = append(resultBlogs, each)
		}

		blog :=  map[string]interface{} {
			"Blogs" : resultBlogs,
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
		id, _ := strconv.Atoi(c.Param("id"))


		var blogDetail = Blog{}

	for i, data := range dataBlogs {
		if id == i {
			blogDetail = Blog{
				ProjectName:    data.ProjectName,
				StartDate:  	data.StartDate,
				EndDate:    	data.EndDate,
				Duration:   	data.Duration,
				Description: 	data.Description,
				Javascript:     data.Javascript,
				PHP:    		data.PHP,
				Java:     		data.Java,
				ReactJS: 		data.ReactJS,
				Image: 			data.Image,
			}
		}
	}

	data := map[string]interface{}{
		"Blog":   blogDetail,
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
		javascript := c.FormValue("input-JavaScript")
		php := c.FormValue("input-php")
		java := c.FormValue("input-java")
		reactJS := c.FormValue("input-reactJS")
		image := c.FormValue("input-image")
		
	
		newBlog := Blog {
			ProjectName: projectName,
			Duration: getDuration(startDate, endDate),
			Author: "Ahmad Syahran Zidane",
			Description: description,
			Javascript: (javascript == "javascript"),
			PHP: (php == "php"),
			Java: (java == "java"),
			ReactJS: (reactJS == "reactJS") ,
			Image: image,
		}


		dataBlogs = append(dataBlogs, newBlog)

		fmt.Println(dataBlogs)
	
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	func deleteBlog(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		dataBlogs = append(dataBlogs[:id], dataBlogs[id+1:]...)
		return c.Redirect(http.StatusMovedPermanently, "/")
	}