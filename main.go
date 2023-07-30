package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"text/template"

	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

// USER STRUCT
type User struct {
	Id int
	Name string
	Email string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name string
}

var userData = SessionData{}


func main() {

	e := echo.New()

	connection.DatabaseConnect()

	e.Static("/public", "public")
	e.Static("/uploads", "uploads")
    // e.GET("/", func(c echo.Context) error {
    //     return c.String(http.StatusOK, "Hello, Saya Ahmad")
    // })
    // e.Logger.Fatal(e.Start("localhost:700"))

	// To Use sessions using echo
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))
	
	e.GET("/", index)
	e.GET("/blog", blog)
	e.GET("/testimonial", testimonial)
	e.GET("/contact", contact)
	e.GET("/detail-blog/:id", detailBlog)
	e.GET("/update-blog-form/:id", updateBlogForm)


	//POST routing
	e.POST("/blog", middleware.UploadFile(addBlog))
	e.POST("/delete-blog/:id", deleteBlog)
	e.POST("/update-blog/:id", updateBlog)

	// REGISTER
	e.GET("/form-register", formRegister)
	e.POST("/register", register)

	//LOGIN
	e.GET("/form-login", formLogin)
	e.POST("/login", login)

	// LOG OUT
	e.POST("/logout", logout)


	e.Logger.Fatal(e.Start("Localhost:5000"))

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

		// GET DATA SESSION FOR LOGIN
	sess, _ := session.Get("session", c)
	
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	// SESSIONS CONDITION
	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	// GET DATA FROM DATABASE
	dataBlog, errBlog := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, project_name, description, image, start_date, end_date, technologies, tb_user.name AS author FROM public.tb_blog JOIN tb_user ON tb_blog.author_id = tb_user.id ORDER BY tb_blog.id DESC")

	if errBlog != nil {
		return c.JSON(http.StatusInternalServerError, errBlog.Error())
	}

	var resultBlogs []Blog
	for dataBlog.Next() {
		var each = Blog{}

		err := dataBlog.Scan(&each.Id, &each.ProjectName, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Technologies, &each.Author)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

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
		"FlashStatus" : sess.Values["status"],
		"FlashMessage" : sess.Values["message"],
		"DataSession" : userData,
	}



	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
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

		errQuery := connection.Conn.QueryRow(context.Background(), `
		SELECT tb_blog.id, project_name, description, image, start_date, end_date, technologies, tb_user.name AS author
		FROM public.tb_blog
		JOIN tb_user ON tb_blog.author_id = tb_user.id
		WHERE tb_blog.id = $1
	`, idToInt).Scan(&blogDetail.Id, &blogDetail.ProjectName, &blogDetail.Description, &blogDetail.Image, &blogDetail.StartDate, &blogDetail.EndDate, &blogDetail.Technologies, &blogDetail.Author)
	

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
		// image := c.FormValue("input-image")

		sess, _ := session.Get("session", c)
		author := sess.Values["id"].(int)

		image := c.Get("dataFile").(string)
	
		

		insertDb, err := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_blog (project_name, description, image, start_date, end_date, technologies, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", projectName, description, image, startDate, endDate, technologies, author)
			
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

	func formRegister (c echo.Context) error {
		var tmpl, err = template.ParseFiles("views/form-register.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message " : err.Error()})
		}

		return tmpl.Execute(c.Response(), nil)
	}

	func register(c echo.Context) error  {
		// to make sure request body is form data format, not JSON , XML etc
		err := c.Request(). ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		name := c.FormValue("input-name")
		email := c.FormValue("input-email")
		passsword := c.FormValue("input-password")

		passswordHash, _ := bcrypt.GenerateFromPassword([]byte(passsword), 10)

		_ , err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passswordHash)
		if err != nil {
			redirectWithMessage(c, "Register failed, please try again", false, "/form-register")
		}

		fmt.Println(err)
		
		return redirectWithMessage(c, "Register success !", true, "/form-login")
}

	func formLogin (c echo.Context) error {

		sess, _ := session.Get("session", c)
	
		flash := map[string]interface{}{
			"FlashStatus" : sess.Values["status"],
			"FlashMessage" : sess.Values["message"],
		}

		delete(sess.Values, "message")
		delete(sess.Values, "status")
		sess.Save(c.Request(), c.Response())
		
		var tmpl, err = template.ParseFiles("views/form-login.html")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message " : err.Error()})
		}

		return tmpl.Execute(c.Response(), flash)
	}

	func login (c echo.Context) error  {
		err := c.Request().ParseForm()

		if err != nil {
			log.Fatal(err)
		}

		email := c.FormValue("input-email")
		passsword := c.FormValue("input-password")

		user := User{}
		err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email,&user.Password)

		if err != nil {
			return redirectWithMessage(c, "Email Incorrect !", false, "form-login")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passsword))
		if err != nil {
			return redirectWithMessage(c, "Password Incorrected", false, "/form-login")
		}

		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = 10800 // Batas maksimal sesi adalah 3 jam
		sess.Values["message"] = "Login Success!"
		sess.Values["status"] = true
		sess.Values["name"] = user.Name
		sess.Values["email"] = user.Email
		sess.Values["id"] = user.Id
		sess.Values["isLogin"] = true
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	func redirectWithMessage (c echo.Context, message string, status bool, path string ) error {
		sess, _ := session.Get("session", c)
		sess.Values["message"] = message
		sess.Values["status"] = status
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusMovedPermanently, path)
	}

	func logout (c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusMovedPermanently, "/")
	}