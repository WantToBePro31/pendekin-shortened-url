package handler

import (
	"net/http"
	"pendekin/dto"
	"pendekin/helper"
	"pendekin/service"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandlerApi interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type userHandlerApi struct {
	userService    service.UserService
	sessionService service.SessionService
	redisService   service.RedisService
}

func InitUserHandlerApi(userService service.UserService, sessionService service.SessionService, redisService service.RedisService) *userHandlerApi {
	return &userHandlerApi{userService, sessionService, redisService}
}

func (uha *userHandlerApi) Register(c *gin.Context) {
	var register dto.UserRegister
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your requested data is invalid!"})
		return
	}

	user, _ := uha.userService.CheckUserByUsernameEmail(register.Username, register.Email)
	if user.ID != 0 {
		if register.Username == user.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken!"})
			return
		}
		if register.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken!"})
			return
		}
	}

	hashed, err := helper.HashPassword(register.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password!"})
		return
	}

	register.Password = hashed

	if err := uha.userService.StoreUser(&register); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "New user data is successfully stored!"})
}

func (uha *userHandlerApi) Login(c *gin.Context) {
	var login dto.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your requested data is invalid!"})
		return
	}

	if login.Username == "" && login.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your requested data is invalid!"})
		return
	}

	user, _ := uha.userService.CheckUserByUsernameEmailPassword(login.Username, login.Email, login.Password)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	jwt, err := helper.GenerateJWT(user.ID, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT!"})
		return
	}

	session := dto.Session{
		JWT:            jwt,
		ExpirationTime: expirationTime,
		UserId:         user.ID,
	}

	if err := uha.sessionService.StoreSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new login session"})
		return
	}

	c.SetCookie("auth_cookie", session.JWT, 18000, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "New login session is successfully stored!"})
}

func (uha *userHandlerApi) Logout(c *gin.Context) {
	cookie, _ := c.Request.Cookie("auth_cookie")
	session, err := uha.sessionService.GetSession(cookie.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get login session"})
		return
	}

	if err := uha.sessionService.RemoveSession(session.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove login session"})
		return
	}

	if err = uha.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear User's shortened URL cache!"})
		return
	}

	c.SetCookie("auth_cookie", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login session is successfully removed!"})
}
