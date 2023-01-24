package handler

import (
	"net/http"
	"pendekin/helper"
	"pendekin/dto"
	"pendekin/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandlerWeb interface {
	Register(c *gin.Context)
	RegisterProcess(c *gin.Context)
	Login(c *gin.Context)
	LoginProcess(c *gin.Context)
	Logout(c *gin.Context)
}

type userHandlerWeb struct {
	userService    service.UserService
	sessionService service.SessionService
	redisService service.RedisService
}

func InitUserHandlerWeb(userService service.UserService, sessionService service.SessionService, redisService service.RedisService) *userHandlerWeb {
	return &userHandlerWeb{userService, sessionService, redisService}
}

func (uhw *userHandlerWeb) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/RegisterPage.html", gin.H{
		"title": "Pendekin - Register",
	})
}

func (uhw *userHandlerWeb) RegisterProcess(c *gin.Context) {
	var register dto.UserRegister
	register.Fullname = c.PostForm("fullname")
	register.Username = c.PostForm("username")
	register.Email = c.PostForm("email")
	register.Password = c.PostForm("password")
	if register.Fullname == "" || register.Username == "" || register.Email == "" || register.Password == "" {
		c.HTML(http.StatusBadRequest, "auth/RegisterPage.html", gin.H{
			"title":   "Pendekin - Register",
			"message": "Please input your valid identity!",
		})
		return
	}

	user, _ := uhw.userService.CheckUserByUsernameEmail(register.Username, register.Email)
	if user.ID != 0 {
		if register.Username == user.Username {
			c.HTML(http.StatusBadRequest, "auth/RegisterPage.html", gin.H{
				"title":   "Pendekin - Register",
				"message": "Username is already taken!",
			})
			return
		}
		if register.Email == user.Email {
			c.HTML(http.StatusBadRequest, "auth/RegisterPage.html", gin.H{
				"title":   "Pendekin - Register",
				"message": "Email is already taken!",
			})
			return
		}
	}

	hashed, err := helper.HashPassword(register.Password)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/register")
		return
	}

	register.Password = hashed

	if err := uhw.userService.StoreUser(&register); err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/register")
		return
	}

	c.HTML(http.StatusOK, "auth/RegisterPage.html", gin.H{
		"title":   "Pendekin - Register",
		"message": "Register success!",
	})
}

func (uhw *userHandlerWeb) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/LoginPage.html", gin.H{
		"title": "Pendekin - Login",
	})
}

func (uhw *userHandlerWeb) LoginProcess(c *gin.Context) {
	var login dto.UserLogin
	identity := c.PostForm("identity")
	if strings.Contains(identity, "@") {
		login.Email = identity
	} else {
		login.Username = identity
	}
	login.Password = c.PostForm("password")
	if identity == "" || login.Password == "" {
		c.HTML(http.StatusBadRequest, "auth/LoginPage.html", gin.H{
			"title":   "Pendekin - Login",
			"message": "Please input your valid identity!",
		})
		return
	}

	user, _ := uhw.userService.CheckUserByUsernameEmailPassword(login.Username, login.Email, login.Password)
	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "auth/LoginPage.html", gin.H{
			"title":   "Pendekin - Login",
			"message": "User not found!",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	jwt, err := helper.GenerateJWT(user.ID, expirationTime)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/login")
		return
	}

	session := dto.Session{
		JWT:            jwt,
		ExpirationTime: expirationTime,
		UserId:         user.ID,
	}

	if err := uhw.sessionService.StoreSession(&session); err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/login")
		return
	}

	c.SetCookie("auth_cookie", session.JWT, 18000, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

func (uha *userHandlerWeb) Logout(c *gin.Context) {
	cookie, _ := c.Request.Cookie("auth_cookie")
	session, err := uha.sessionService.GetSession(cookie.Value)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/logout")
		return
	}

	if err := uha.sessionService.RemoveSession(session.ID); err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/logout")
		return
	}

	if err = uha.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.Redirect(http.StatusInternalServerError, "/user/logout")
		return
	}

	c.SetCookie("auth_cookie", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/user/login")
}
