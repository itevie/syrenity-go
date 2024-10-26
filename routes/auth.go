package routes

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"syrenity/server/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	router.POST("/get-token", func(c *gin.Context) {
		var auth struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindBodyWithJSON(&auth); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Invalid request body",
			})
			return
		}

		var user models.User
		if err := db.QueryRowx("SELECT * FROM users WHERE email = $1;", auth.Email).StructScan(&user); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, models.ErrorMessage{
				Message: "No user is associated with the provided email",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password)); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, models.ErrorMessage{
				Message: "Invalid email or password",
			})
			return
		}

		var token models.Token
		err := db.QueryRowx(
			"INSERT INTO tokens (token, account, identifier) VALUES ($1, $2, $3) RETURNING *;",
			models.GenerateToken(),
			user.ID,
			"/auth/get-token",
		).StructScan(&token)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to create token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token.Token,
		})
	})

	router.POST("/register", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Invalid request body",
			})
			return
		}

		fmt.Println(body, body.Username)

		var email = strings.ToLower(body.Email)
		var username = strings.ToLower(body.Username)
		var password = body.Password

		if _, err := mail.ParseAddress(email); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Invalid email address",
			})
			return
		}

		if len(username) < 3 || len(username) > 15 {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Username must be between 3-15 characters",
			})
			return
		}

		re := regexp.MustCompile("^[a-z0-9_-]+$")
		if !re.MatchString(username) {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Username must only contain alphanumeric characters or _, -",
			})
			return
		}

		if len(password) < 10 || len(password) > 100 {
			c.JSON(http.StatusBadRequest, models.ErrorMessage{
				Message: "Password must be between 10-100 characters",
			})
			return
		}

		var _user models.User
		userCheckErr := db.QueryRowx("SELECT * FROM users WHERE email = $1;", email).StructScan(&_user)
		if userCheckErr == nil {
			c.JSON(http.StatusConflict, models.ErrorMessage{
				Message: "Email already used",
			})
			return
		}

		var count int
		if err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE username = $1;", username); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to get amount of users with the username",
			})
			return
		}

		if count > 9995 {
			c.JSON(http.StatusConflict, models.ErrorMessage{
				Message: "Too many users have this username",
			})
			return
		}

		var discriminator = models.RandomDiscriminator()

		var user models.User
		err := db.QueryRowx("INSERT INTO users (username, password, email, discriminator) VALUES ($1, $2, $3, $4) RETURNING *", username, password, email, discriminator).StructScan(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, models.ErrorMessage{
				Message: "Failed to register user",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	})
}
