package types

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var ErrTokenExpired = jwt.ErrTokenExpired
var ErrRecordNotFound = gorm.ErrRecordNotFound
