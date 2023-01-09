package dto

type UrlInsert struct {
	RealUrl      string `json:"real_url" binding:"required"`
	ShortenedUrl string `json:"shortened_url"`
	Randomized   bool   `json:"randomized" binding:"required"`
	UserId       uint   `json:"user_id"`
}

type UrlUpdate struct {
	RealUrl      string `json:"real_url"`
	ShortenedUrl string `json:"shortened_url" binding:"required"`
}
