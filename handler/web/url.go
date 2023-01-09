package handler

import (
	"net/http"
	"pendekin/dto"
	"pendekin/helper"
	"pendekin/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UrlHandlerWeb interface {
	Index(c *gin.Context)
	Home(c *gin.Context)
	Real(c *gin.Context)
	Dashboard(c *gin.Context)
	CreateShortenedUrl(c *gin.Context)
	Result(c *gin.Context)
	Update(c *gin.Context)
	EditShortenedUrls(c *gin.Context)
	RemoveShortenedUrls(c *gin.Context)
}

type urlHandlerWeb struct {
	urlService     service.UrlService
	sessionService service.SessionService
	redisService   service.RedisService
}

func InitUrlHandlerWeb(urlService service.UrlService, sessionService service.SessionService, redisService service.RedisService) *urlHandlerWeb {
	return &urlHandlerWeb{urlService, sessionService, redisService}
}

func (uhw *urlHandlerWeb) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/url")
}

func (uhw *urlHandlerWeb) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "general/HomePage.html", gin.H{
		"title": "Pendekin - Home",
	})
}

func (uhw *urlHandlerWeb) Real(c *gin.Context) {
	shortened_url := c.Param("shortened")
	url, err := uhw.urlService.GetShortenedUrl(shortened_url)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/ErrorPage.html", gin.H{
			"title": "Pendekin - Error",
		})
		return
	}

	c.Redirect(http.StatusFound, url.RealUrl)
}

func (uhw *urlHandlerWeb) Dashboard(c *gin.Context) {
	cookie, _ := c.Request.Cookie("auth_cookie")
	urls, err := uhw.redisService.Get("urlsfrom" + cookie.Value)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/dashboard")
		return
	}

	if len(urls) == 0 {
		session, err := uhw.sessionService.GetSession(cookie.Value)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/dashboard")
			return
		}

		urls, err = uhw.urlService.GetAllShortenedUrls(session.UserId)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/dashboard")
			return
		}

		if err := uhw.redisService.Store("urlsfrom"+cookie.Value, urls); err != nil {
			c.Redirect(http.StatusInternalServerError, "/dashboard")
			return
		}
	}

	c.HTML(http.StatusOK, "main/DashboardPage.html", gin.H{
		"title": "Pendekin - Dashboard",
		"urls":  urls,
	})
}

func (uhw *urlHandlerWeb) CreateShortenedUrl(c *gin.Context) {
	var pendekin dto.UrlInsert
	pendekin.RealUrl = c.PostForm("real_url")
	pendekin.ShortenedUrl = c.PostForm("shortened_url")
	pendekin.Randomized, _ = strconv.ParseBool(c.PostForm("randomized"))
	if pendekin.RealUrl == "" {
		c.HTML(http.StatusBadRequest, "general/HomePage.html", gin.H{
			"title":   "Pendekin - Home",
			"message": "Please input your valid URL!",
		})
		return
	}

	if !strings.HasPrefix(pendekin.RealUrl, "http://") && !strings.HasPrefix(pendekin.RealUrl, "https://") {
		pendekin.RealUrl = "https://" + pendekin.RealUrl
	}

	if pendekin.ShortenedUrl == "" && !pendekin.Randomized {
		c.HTML(http.StatusBadRequest, "general/HomePage.html", gin.H{
			"title":   "Pendekin - Home",
			"message": "Please input your shortened URL!",
		})
		return
	}

	if pendekin.Randomized {
		nunique := true
		for nunique {
			new_url := helper.GenerateRandomUrl()
			url, _ := uhw.urlService.GetShortenedUrl(new_url)
			if url.ID == 0 {
				pendekin.ShortenedUrl = new_url
				nunique = false
			}
		}
	} else {
		url, _ := uhw.urlService.GetShortenedUrl(pendekin.ShortenedUrl)
		if url.ID != 0 {
			c.HTML(http.StatusBadRequest, "general/HomePage.html", gin.H{
				"title":   "Pendekin - Home",
				"message": "Shortened URL is already exist!",
			})
			return
		}
	}

	cookie, err := c.Request.Cookie("auth_cookie")
	if err == nil {
		session, err := uhw.sessionService.GetSession(cookie.Value)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/url/create")
			return
		}
		pendekin.UserId = session.UserId
	}

	if err := uhw.urlService.StoreShortenedUrl(&pendekin); err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/create")
		return
	}

	if cookie != nil {
		if err = uhw.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
			c.Redirect(http.StatusInternalServerError, "/url/create")
			return
		}
	}

	c.Redirect(http.StatusFound, "/url/"+pendekin.ShortenedUrl)
}

func (uhw *urlHandlerWeb) Result(c *gin.Context) {
	shortened_url := c.Param("shortened")
	url, err := uhw.urlService.GetShortenedUrl(shortened_url)
	if err != nil {
		if err.Error() == "record not found" {
			c.HTML(http.StatusNotFound, "error/ErrorPage.html", gin.H{
				"title": "Pendekin - Error",
			})
			return
		}
		c.Redirect(http.StatusInternalServerError, "/url/create")
		return
	}

	c.HTML(http.StatusOK, "general/ResultPage.html", gin.H{
		"title":        "Pendekin - Result",
		"realUrl":      url.RealUrl,
		"shortenedUrl": url.ShortenedUrl,
		"message":      "success create new shortened url",
	})
}

func (uhw *urlHandlerWeb) Update(c *gin.Context) {
	urlId := c.Param("url_id")
	iurlId, _ := strconv.ParseUint(urlId, 10, 32)

	url, err := uhw.urlService.GetShortenedUrlById(uint(iurlId))
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/update/"+urlId)
		return
	}

	c.HTML(http.StatusOK, "main/UpdatePage.html", gin.H{
		"title": "Pendekin - Update",
		"id": urlId,
		"url": url,
	})
}

func (uhw *urlHandlerWeb) EditShortenedUrls(c *gin.Context) {
	urlId := c.Param("url_id")
	iurlId, _ := strconv.ParseUint(urlId, 10, 32)

	var pendekin dto.UrlUpdate
	pendekin.RealUrl = c.PostForm("real_url")
	pendekin.ShortenedUrl = c.PostForm("shortened_url")

	if err := uhw.urlService.EditShortenedUrl(&pendekin, uint(iurlId)); err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/update/"+urlId)
		return
	}

	cookie, _ := c.Request.Cookie("auth_cookie")
	if err := uhw.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/update/"+urlId)
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (uhw *urlHandlerWeb) RemoveShortenedUrls(c *gin.Context) {
	urlId := c.Param("url_id")
	iurlId, _ := strconv.ParseUint(urlId, 10, 32)

	if err := uhw.urlService.DeleteShortenedUrl(uint(iurlId)); err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/delete/"+urlId)
		return
	}

	cookie, _ := c.Request.Cookie("auth_cookie")
	if err := uhw.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.Redirect(http.StatusInternalServerError, "/url/delete/"+urlId)
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}
