package controllers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/andreasatle/react-go/go-auth/database"
	"github.com/andreasatle/react-go/go-auth/routes/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.StandardClaims
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Register(c *fiber.Ctx) error {
	data := map[string]string{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	res := database.DB.Create(&user)
	if res.Error != nil {
		return fmt.Errorf("Error inserting a new entry in DB: %v", res.Error)
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	data := map[string]string{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	user := models.User{}
	res := database.DB.Where("email = ?", data["email"]).First(&user)
	if res.Error != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect password!",
		})
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: expiresAt.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiresAt,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"jwt": token,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*Claims)
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "logout successful",
	})

}

func ForgotPassword(c *fiber.Ctx) error {
	data := map[string]string{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}
	database.DB.Create(&passwordReset)

	from := "admin@example.com"
	to := []string{data["email"]}
	url := "http://localhost:3000/reset/" + token
	message := []byte("Click on <a href \"" + url + "\">here</a> to reset your password!")

	err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "password reset successful",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	data := map[string]string{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	passwordReset := models.PasswordReset{}
	res := database.DB.Where("token = ?", data["token"]).Last(&passwordReset)
	if res.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	res = database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)
	if res.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Error updating database with new password!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "password reset successful",
	})
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
