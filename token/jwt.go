package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/myesui/uuid"
	"os"
	"strings"
	"time"
)

const HeaderAuthorization = "Authorization"
const ParamAccessToken = "access_token"
const ParamRefreshToken = "refresh_token"
const ParamVTM = "vtm"
const ParamLTM = "ltm"
const ParamRTM = "rtm"
const ErrorInvalidVTM = "Invalid VTM!"
const ErrorInvalidLTM = "Invalid LTM!"
const ErrorInvalidRTM = "Invalid RTM!"

type VerifyTokenMeta struct {
	AccessUuid         string
	AccessTokenExpires int64
	AccessToken        string
}

type LoginTokenMeta struct {
	AccessUuid         string
	AccessToken        string
	AccessTokenExpires int64
	UserUniqueID       string
}

type RefreshTokenMeta struct {
	RefreshUuid         string
	RefreshToken        string
	RefreshTokenExpires int64
	UserUniqueID        string
}

func CreateVTM() (*VerifyTokenMeta, error) {
	vtm := &VerifyTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	vtmClaims := jwt.MapClaims{}
	vtmClaims["access_uuid"] = vtm.AccessUuid
	vtmClaims["exp"] = vtm.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, vtmClaims)
	vtm.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return vtm, nil
}

//CreateLTMAndRTM is used to create login and refresh token meta
func CreateLTMAndRTM(userUniqueID string) (*LoginTokenMeta, *RefreshTokenMeta, error) {
	ltm := &LoginTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	rtm := &RefreshTokenMeta{
		RefreshUuid:         uuid.NewV4().String(),
		RefreshTokenExpires: time.Now().Add(time.Hour * 24 * 31).Unix(),
	}

	var err error
	//Creating login token meta
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	ltmClaims := jwt.MapClaims{}
	ltmClaims["access_uuid"] = ltm.AccessUuid
	ltmClaims["user_unique_id"] = userUniqueID
	ltmClaims["exp"] = ltm.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ltmClaims)
	ltm.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	//Creating refresh token meta
	rtmClaims := jwt.MapClaims{}
	rtmClaims["refresh_uuid"] = rtm.RefreshUuid
	rtmClaims["user_unique_id"] = userUniqueID
	rtmClaims["exp"] = rtm.RefreshTokenExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtmClaims)
	rtm.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	return ltm, rtm, nil
}

//ExtractToken to extract token from headers or return string value
func ExtractToken(bearerToken string) string {
	//Normally Authorization the_token_xxx
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	//Return simple token string
	return bearerToken
}

//VerifyToken to verify token
func VerifyToken(bearerToken string) (*jwt.Token, error) {
	tokenString := ExtractToken(bearerToken)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		os.Setenv("ACCESS_SECRET", "JWT_SECRET")
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//ExtractVTM is used to return VTM and check if bearer token is valid or not
func ExtractVTM(bearerToken string) (*VerifyTokenMeta, error) {
	token, err := VerifyToken(bearerToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		return &VerifyTokenMeta{
			AccessUuid: accessUuid,
		}, nil
	}
	return nil, err
}

//ExtractLTM is used to return LTM and check if bearer token is valid or not
func ExtractLTM(bearerToken string) (*LoginTokenMeta, error) {
	token, err := VerifyToken(bearerToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("access_uuid is empty")
		}
		atExpires := int64(claims["exp"].(float64))
		if atExpires == 0 {
			return nil, errors.New("exp is empty")
		}
		userUniqueID, ok := claims["user_unique_id"].(string)
		if !ok {
			return nil, errors.New("user_unique_id is empty")
		}
		return &LoginTokenMeta{
			AccessUuid:         accessUuid,
			UserUniqueID:       userUniqueID,
			AccessTokenExpires: atExpires,
		}, nil
	}
	return nil, err
}

//ExtractRTM is used to return RTM and check if bearer token is valid or not
func ExtractRTM(bearerToken string) (*RefreshTokenMeta, error) {
	token, err := VerifyToken(bearerToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, errors.New("access_uuid is empty")
		}
		rtExpires := int64(claims["exp"].(float64))
		if rtExpires == 0 {
			return nil, errors.New("exp is empty")
		}
		userUniqueID, ok := claims["user_unique_id"].(string)
		if !ok {
			return nil, errors.New("user_unique_id is empty")
		}
		return &RefreshTokenMeta{
			RefreshUuid:         refreshUuid,
			UserUniqueID:        userUniqueID,
			RefreshTokenExpires: rtExpires,
		}, nil
	}
	return nil, err
}
