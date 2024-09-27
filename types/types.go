package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type ActiveUser struct {
	ID    uint     `json:"id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	ID    uint     `json:"id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	ID uint `json:"id"`
}

type (
	VerificationType string
	UserRole         string
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	VerificationTypeEmail         VerificationType = "email"
	VerificationTypePasswordReset VerificationType = "password_reset"
	UserRoleAdmin                 UserRole         = "admin"
	UserRoleUser                  UserRole         = "user"
	UserRoleSuperAdmin            UserRole         = "super_admin"
)

const ACTIVE_USER_KEY = "activeuser"

type Paginator interface {
	GetLimit() int
	GetPage() int
	GetSort() string
	SetTotal(total int)
	SetTotalPages(totalPages int)
}

type PageReturn struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
}

func (p *PageReturn) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PageReturn) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *PageReturn) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at"
	}
	return p.Sort
}

func (p *PageReturn) SetTotal(total int) {
	p.Total = total
}

func (p *PageReturn) SetTotalPages(totalPages int) {
	p.TotalPages = totalPages
}

func A(paginator Paginator) {

}

type UserResponse struct {
	PageReturn
	Data string
}


func (u *UserResponse) Sort () {
	u.PageReturn.Sort = "sdsd "
}

func S(){
	var userResponse UserResponse

	A(&userResponse)
}