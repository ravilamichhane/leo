package web

import (
	"bytes"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/ravilmc/leo/types"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func DecodeAndValidate(c echo.Context, v interface{}) error {
	contentType := c.Request().Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		val := reflect.ValueOf(v)
		valType := reflect.TypeOf(v).Elem()
		refVal := val.Elem()

		for i := 0; i < refVal.NumField(); i++ {
			field := refVal.Field(i)
			fieldType := valType.Field(i)

			if (fieldType.Type.AssignableTo(reflect.TypeOf(FileUpload{}))) {

				tag := fieldType.Tag.Get("file")
				if tag == "" {
					continue
				}

				tagslice := strings.Split(tag, ",")

				fieldName := tagslice[0]

				file, err := c.FormFile(fieldName)
				if err != nil {
					continue
				}

				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				buf := new(bytes.Buffer)

				buf.ReadFrom(src)

				var a FileUpload

				a.Data = buf.Bytes()
				a.Ext = filepath.Ext(file.Filename)
				a.Mime = http.DetectContentType(a.Data)

				if a.Ext == "" {
					exts, _ := mime.ExtensionsByType(a.Mime)

					if len(exts) > 0 {
						a.Ext = exts[0]
					} else {
						a.Ext = ".txt"
					}

				}

				if field.CanSet() {
					field.Set(reflect.ValueOf(a))
				}

			}
		}
	}

	if err := c.Bind(v); err != nil {

		return err

	}

	if err := c.Validate(v); err != nil {
		return err
	}

	return nil
}

func GenerateOtp(length int) string {
	available := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	otp := ""
	for i := 0; i < length; i++ {
		otp += string(available[rand.Intn(len(available))])
	}
	return otp
}

func BoolPointer(b bool) *bool {
	return &b
}

func GetActiveUser(c echo.Context) *types.ActiveUser {
	user, ok := c.Get(types.ACTIVE_USER_KEY).(*types.ActiveUser)
	if !ok {
		return nil
	}
	return user
}

func IsAdmin(activeUser *types.ActiveUser) bool {
	if activeUser == nil {
		return false
	}

	return activeUser.Role == types.UserRoleAdmin || activeUser.Role == types.UserRoleSuperAdmin
}

func IsSuperAdmin(activeUser *types.ActiveUser) bool {
	if activeUser == nil {
		return false
	}
	return activeUser.Role == types.UserRoleSuperAdmin
}

func SetActiveUser(c echo.Context, claims *types.CustomClaims) {
	c.Set(types.ACTIVE_USER_KEY, &types.ActiveUser{
		ID:    claims.ID,
		Email: claims.Email,
		Role:  claims.Role,
	})
}

func GetParamUint(c echo.Context, name string) (uint64, error) {
	param := c.Param(name)

	return strconv.ParseUint(param, 10, 64)
}
