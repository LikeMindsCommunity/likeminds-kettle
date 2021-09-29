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

type VerifyTokenMeta struct {
	AccessUuid         string
	VerifiedMobileNo   string
	CountryCode        string
	AccessTokenExpires int64
	AccessToken        string
}

type LoginTokenMeta struct {
	AccessUuid         string
	AccessToken        string
	AccessTokenExpires int64
	VerifiedMobileNo   string
	CountryCode        string
	UserID             string
}

type RefreshTokenMeta struct {
	RefreshUuid         string
	RefreshToken        string
	RefreshTokenExpires int64
	VerifiedMobileNo    string
	CountryCode         string
	UserID              string
}

func CreateVTM(verifiedMobileNo string, countryCode string) (*VerifyTokenMeta, error) {
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

//CreateLTMAndRTM is used to create login and refresh token meta
func CreateLTMAndRTM(verifiedMobileNo string, countryCode string, userID string) (*LoginTokenMeta, *RefreshTokenMeta, error) {
	ltm := &LoginTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		VerifiedMobileNo:   verifiedMobileNo,
		CountryCode:        countryCode,
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	rtm := &RefreshTokenMeta{
		RefreshUuid:         uuid.NewV4().String(),
		VerifiedMobileNo:    verifiedMobileNo,
		CountryCode:         countryCode,
		RefreshTokenExpires: time.Now().Add(time.Hour * 24 * 31).Unix(),
	}

	var err error
	//Creating login token meta
	os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	ltmClaims := jwt.MapClaims{}
	ltmClaims["access_uuid"] = ltm.AccessUuid
	ltmClaims["verified_mobile_no"] = ltm.VerifiedMobileNo
	ltmClaims["country_code"] = ltm.CountryCode
	ltmClaims["user_id"] = userID
	ltmClaims["exp"] = ltm.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ltmClaims)
	ltm.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	//Creating refresh token meta
	rtmClaims := jwt.MapClaims{}
	rtmClaims["refresh_uuid"] = rtm.RefreshUuid
	rtmClaims["verified_mobile_no"] = rtm.VerifiedMobileNo
	rtmClaims["country_code"] = rtm.CountryCode
	rtmClaims["user_id"] = userID
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
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("user_id is empty")
		}
		return &RefreshTokenMeta{
			RefreshUuid:         refreshUuid,
			UserID:              userID,
			RefreshTokenExpires: rtExpires,
		}, nil
	}
	return nil, err
}
