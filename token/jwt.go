package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/myesui/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type Details struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type VerifyTokenMeta struct {
	AccessUuid       string
	VerifiedMobileNo string
	CountryCode      string
}

type LoginTokenMeta struct {
	AccessUuid       string
	AtExpires        int64
	VerifiedMobileNo string
	CountryCode      string
	UserID           string
}

func CreateVerifyOTPToken(verifiedMobileNo string, countryCode string) (*Details, error) {
	td := &Details{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["verified_mobile_no"] = verifiedMobileNo
	atClaims["country_code"] = countryCode
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateLoginToken(meta *VerifyTokenMeta, userID string) (*Details, error) {
	td := &Details{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["verified_mobile_no"] = meta.VerifiedMobileNo
	atClaims["country_code"] = meta.CountryCode
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

func ExtractVerifyTokenMeta(r *http.Request) (*VerifyTokenMeta, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		verifiedMobileNo, ok := claims["verified_mobile_no"].(string)
		if !ok {
			return nil, err
		}
		countryCode, ok := claims["country_code"].(string)
		if !ok {
			return nil, err
		}
		return &VerifyTokenMeta{
			AccessUuid:       accessUuid,
			VerifiedMobileNo: verifiedMobileNo,
			CountryCode:      countryCode,
		}, nil
	}
	return nil, err
}

func ExtractLoginTokenMeta(r *http.Request) (*LoginTokenMeta, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		atExpires, ok := claims["exp"].(int64)
		if !ok {
			return nil, err
		}
		verifiedMobileNo, ok := claims["verified_mobile_no"].(string)
		if !ok {
			return nil, err
		}
		countryCode, ok := claims["country_code"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		return &LoginTokenMeta{
			AccessUuid:       accessUuid,
			AtExpires:        atExpires,
			VerifiedMobileNo: verifiedMobileNo,
			CountryCode:      countryCode,
			UserID:           userID,
		}, nil
	}
	return nil, err
}
