package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.Static("/public", "public")
    // e.GET("/", func(c echo.Context) error {
    //     return c.String(http.StatusOK, "Hello, Saya Ahmad")
    // })
    // e.Logger.Fatal(e.Start("localhost:700"))
	
	e.GET("/", index)
	e.GET("/blog", blog)
	e.GET("/testimonial", testimonial)
	e.GET("/contact", contact)
	e.GET("/blog-detail/:id", blogDetail)

	e.POST("/blog", addBlog)



	e.Logger.Fatal(e.Start("localhost:5000"))

}

	func index (c echo.Context) error  {
		tmpl, err := template.ParseFiles("views/index.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return tmpl.Execute(c.Response(), nil)
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

	func blogDetail (c echo.Context) error {
		id := c.Param("id") 

		tmpl, err := template.ParseFiles("views/blog-detail.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		blogDetail := map[string]interface{}{ // interface -> tipe data apapun
			"Id":      id,
			"Title":   "Dumbways ID memang keren",
			"Content": "Dumbways ID adalah bootcamp terbaik sedunia seakhirat!",
		}
	
		return tmpl.Execute(c.Response(), blogDetail)
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
	
		fmt.Println("name :", projectName)
		fmt.Println("start :", startDate)
		fmt.Println("end :", endDate)
		fmt.Println("description: ", description)
		fmt.Println("Nilai dari checkbox :", javascript)
		fmt.Println("Nilai dari checkbox :", php)
		fmt.Println("Nilai dari checkbox :", java)
		fmt.Println("Nilai dari checkbox :", reactJS)
		fmt.Println(image)
	
		return c.Redirect(http.StatusMovedPermanently, "/blog")
	}