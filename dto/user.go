package dto

type UserRegister struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}
