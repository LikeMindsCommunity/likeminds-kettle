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
	AccessUuid         string
	VerifiedMobileNo   string
	CountryCode        string
	AccessTokenExpires int64
	AccessToken        string
}

type LoginTokenMeta struct {
	AccessUuid          string
	RefreshUuid         string
	VerifiedMobileNo    string
	CountryCode         string
	UserID              string
	AccessTokenExpires  int64
	AccessToken         string
	RefreshTokenExpires int64
	RefreshToken        string
}

func CreateVerifyOTPToken(verifiedMobileNo string, countryCode string) (*VerifyTokenMeta, error) {
	vtm := &VerifyTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		VerifiedMobileNo:   verifiedMobileNo,
		CountryCode:        countryCode,
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = vtm.AccessUuid
	atClaims["verified_mobile_no"] = verifiedMobileNo
	atClaims["country_code"] = countryCode
	atClaims["exp"] = vtm.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	vtm.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return vtm, nil
}

func CreateLTM(verifiedMobileNo string, countryCode string, userID string) (*LoginTokenMeta, error) {
	ltm := &LoginTokenMeta{
		AccessUuid:          uuid.NewV4().String(),
		RefreshUuid:         uuid.NewV4().String(),
		VerifiedMobileNo:    verifiedMobileNo,
		CountryCode:         countryCode,
		AccessTokenExpires:  time.Now().Add(time.Minute * 15).Unix(),
		RefreshTokenExpires: time.Now().Add(time.Hour * 24 * 31).Unix(),
	}

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = ltm.AccessUuid
	atClaims["verified_mobile_no"] = ltm.VerifiedMobileNo
	atClaims["country_code"] = ltm.CountryCode
	atClaims["user_id"] = userID
	atClaims["exp"] = ltm.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	ltm.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = ltm.RefreshUuid
	rtClaims["verified_mobile_no"] = ltm.VerifiedMobileNo
	rtClaims["country_code"] = ltm.CountryCode
	rtClaims["user_id"] = userID
	rtClaims["exp"] = ltm.RefreshTokenExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	ltm.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return ltm, nil
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

func ExtractVTM(r *http.Request) (*VerifyTokenMeta, error) {
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

func ExtractAccessLTM(r *http.Request) (*LoginTokenMeta, error) {
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
			AccessUuid:         accessUuid,
			VerifiedMobileNo:   verifiedMobileNo,
			CountryCode:        countryCode,
			UserID:             userID,
			AccessTokenExpires: atExpires,
		}, nil
	}
	return nil, err
}

func ExtractRefreshLTM(r *http.Request) (*LoginTokenMeta, error) {
	token, err := VerifyToken(r)
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
			RefreshUuid:         refreshUuid,
			VerifiedMobileNo:    verifiedMobileNo,
			CountryCode:         countryCode,
			UserID:              userID,
			RefreshTokenExpires: rtExpires,
		}, nil
	}
	return nil, err
}
