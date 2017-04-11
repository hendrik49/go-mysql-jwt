package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"	
	"github.com/appleboy/gin-jwt"
)
var mySigningKey = []byte("didok49")

type Person struct {
	ID         int
	FirstName string
	LastName  string
}

func personHandler(c *gin.Context) {

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/godb")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	var (
		person  Person
		persons []Person
	)
	rows, err := db.Query("select id, firstName, lastName from person;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&person.ID, &person.FirstName, &person.LastName)
		persons = append(persons, person)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": persons,
		"count":  len(persons),
	})
}

func personbyID(c *gin.Context) {
	
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/godb")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	var (
		person Person
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select id, firstName, lastName from person where id = ?;", id)
	err = row.Scan(&person.ID, &person.FirstName, &person.LastName)
	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": person,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
	}

	func savePerson(c *gin.Context) {
	
		var buffer bytes.Buffer
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/godb")
		if err != nil {
			fmt.Print(err.Error())
		}
		defer db.Close()
		// make sure connection is available
		err = db.Ping()
		if err != nil {
			fmt.Print(err.Error())
		}

		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		stmt, err := db.Prepare("insert into person (firstName, lastName) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(firstName, lastName)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(firstName)
		buffer.WriteString(" ")
		buffer.WriteString(lastName)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully", name),
		})
	}

	func updatePerson(c *gin.Context) {

		var buffer bytes.Buffer
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/godb")
		if err != nil {
			fmt.Print(err.Error())
		}
		defer db.Close()
		// make sure connection is available
		err = db.Ping()
		if err != nil {
			fmt.Print(err.Error())
		}

		id := c.Query("id")
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		stmt, err := db.Prepare("update person set firstName= ?, lastName= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(firstName, lastName, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(firstName)
		buffer.WriteString(" ")
		buffer.WriteString(lastName)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	}

	func deletePerson(c *gin.Context) {
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/godb")
		if err != nil {
			fmt.Print(err.Error())
		}
		defer db.Close()
		// make sure connection is available
		err = db.Ping()
		if err != nil {
			fmt.Print(err.Error())
		}

	id := c.Query("id")
	stmt, err := db.Prepare("delete from person where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})
}


func main() {
	router := gin.Default()	
	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("didok49"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	router.POST("/login", authMiddleware.LoginHandler)

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// GET a person detail
		auth.GET("/person/:id",personbyID) 
		// GET all persons
		auth.GET("/persons",personHandler) 
		// POST new person details
		auth.POST("/person",savePerson) 
		// PUT - update a person details
		auth.PUT("/person", updatePerson)
		// Delete resources
		auth.DELETE("/person", deletePerson)
	}
	http.ListenAndServe(":3000", router)
}


