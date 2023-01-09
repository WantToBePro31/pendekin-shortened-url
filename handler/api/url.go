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

type UrlHandlerApi interface {
	GetShortenedUrls(c *gin.Context)
	CreateShortenedUrl(c *gin.Context)
	EditShortenedUrl(c *gin.Context)
	RemoveShortenedUrl(c *gin.Context)
}

type urlHandlerApi struct {
	urlService     service.UrlService
	sessionService service.SessionService
	redisService   service.RedisService
}

func InitUrlHandlerApi(urlService service.UrlService, sessionService service.SessionService, redisService service.RedisService) *urlHandlerApi {
	return &urlHandlerApi{urlService, sessionService, redisService}
}

func (uha *urlHandlerApi) GetShortenedUrls(c *gin.Context) {
	cookie, _ := c.Request.Cookie("auth_cookie")
	urls, err := uha.redisService.Get("urlsfrom" + cookie.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get User's shortened URLs cache!"})
		return
	}

	if len(urls) == 0 {
		session, err := uha.sessionService.GetSession(cookie.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User's session not found!"})
			return
		}

		urls, err = uha.urlService.GetAllShortenedUrls(session.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get User's shortened URLs!"})
			return
		}

		if err := uha.redisService.Store("urlsfrom"+cookie.Value, urls); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new shortened URL cache!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": urls})
}

func (uha *urlHandlerApi) CreateShortenedUrl(c *gin.Context) {
	var pendekin dto.UrlInsert
	if err := c.ShouldBindJSON(&pendekin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your requested data is invalid!"})
		return
	}

	if !strings.HasPrefix(pendekin.RealUrl, "http://") && !strings.HasPrefix(pendekin.RealUrl, "https://") {
		pendekin.RealUrl = "https://" + pendekin.RealUrl
	}

	if pendekin.ShortenedUrl == "" && !pendekin.Randomized {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot automatically create shortened URL without randomized!"})
		return
	}

	if pendekin.Randomized {
		nunique := true
		for nunique {
			new_url := helper.GenerateRandomUrl()
			url, _ := uha.urlService.GetShortenedUrl(new_url)
			if url.ID == 0 {
				pendekin.ShortenedUrl = new_url
				nunique = false
			}
		}
	} else {
		url, _ := uha.urlService.GetShortenedUrl(pendekin.ShortenedUrl)
		if url.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shortened URL is already exist!"})
			return
		}
	}

	cookie, err := c.Request.Cookie("auth_cookie")
	if err == nil {
		session, err := uha.sessionService.GetSession(cookie.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User's session not found!"})
			return
		}
		pendekin.UserId = session.UserId
	}

	if err := uha.urlService.StoreShortenedUrl(&pendekin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new shortened URL!"})
		return
	}

	if cookie != nil {
		if err = uha.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear User's shortened URL cache!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "New shortened URL is successfully stored!"})
}

func (uha *urlHandlerApi) EditShortenedUrls(c *gin.Context) {
	urlId, _ := strconv.ParseUint(c.Param("url_id"), 10, 32)

	var pendekin dto.UrlUpdate
	if err := c.ShouldBindJSON(&pendekin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your requested data is invalid!"})
		return
	}

	if pendekin.RealUrl != "" {
		if !strings.HasPrefix(pendekin.RealUrl, "http://") && !strings.HasPrefix(pendekin.RealUrl, "https://") {
			pendekin.RealUrl = "https://" + pendekin.RealUrl
		}
	}

	if err := uha.urlService.EditShortenedUrl(&pendekin, uint(urlId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit shortened URL!"})
		return
	}

	cookie, _ := c.Request.Cookie("auth_cookie")
	if err := uha.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear User's shortened URL cache!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shortened URL is successfully edited!"})
}

func (uha *urlHandlerApi) RemoveShortenedUrls(c *gin.Context) {
	urlId, _ := strconv.ParseUint(c.Param("url_id"), 10, 32)

	if err := uha.urlService.DeleteShortenedUrl(uint(urlId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shortened URL!"})
		return
	}

	cookie, _ := c.Request.Cookie("auth_cookie")
	if err := uha.redisService.Clear("urlsfrom" + cookie.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear User's shortened URL cache!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shortened URL is successfully deleted!"})
}
