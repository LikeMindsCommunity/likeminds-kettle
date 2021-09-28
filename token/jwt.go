package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/myesui/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type VerifyTokenMeta struct {
	AccessUuid       string
	VerifiedMobileNo string
	CountryCode      string
	ATExpires        int64
	AToken           string
}

type LoginTokenMeta struct {
	AccessUuid       string
	VerifiedMobileNo string
	CountryCode      string
	UserID           string
	ATExpires        int64
	AToken           string
}

func CreateVerifyOTPToken(verifiedMobileNo string, countryCode string) (*VerifyTokenMeta, error) {
	vtm := &VerifyTokenMeta{
		AccessUuid:       uuid.NewV4().String(),
		VerifiedMobileNo: verifiedMobileNo,
		CountryCode:      countryCode,
		ATExpires:        time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = vtm.AccessUuid
	atClaims["verified_mobile_no"] = verifiedMobileNo
	atClaims["country_code"] = countryCode
	atClaims["exp"] = vtm.ATExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	vtm.AToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return vtm, nil
}

func CreateLoginToken(meta *VerifyTokenMeta, userID string) (*LoginTokenMeta, error) {
	td := &LoginTokenMeta{
		AccessUuid:       uuid.NewV4().String(),
		VerifiedMobileNo: meta.VerifiedMobileNo,
		CountryCode:      meta.CountryCode,
		ATExpires:        time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["verified_mobile_no"] = meta.VerifiedMobileNo
	atClaims["country_code"] = meta.CountryCode
	atClaims["user_id"] = userID
	atClaims["exp"] = td.ATExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
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
		return []byte("JWT_SECRET"), nil
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
			return nil, errors.New("access_uuid is empty")
		}
		atExpires := int64(claims["exp"].(float64))
		if atExpires == 0 {
			return nil, errors.New("exp is empty")
		}
		verifiedMobileNo, ok := claims["verified_mobile_no"].(string)
		if !ok {
			return nil, errors.New("verified_mobile_no is empty")
		}
		countryCode, ok := claims["country_code"].(string)
		if !ok {
			return nil, errors.New("country_code is empty")
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("user_id is empty")
		}
		return &LoginTokenMeta{
			AccessUuid:       accessUuid,
			VerifiedMobileNo: verifiedMobileNo,
			CountryCode:      countryCode,
			UserID:           userID,
			ATExpires:        atExpires,
		}, nil
	}
	return nil, err
}
