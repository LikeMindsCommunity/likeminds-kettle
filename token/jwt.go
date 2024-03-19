package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/myesui/uuid"
	"github.com/nateshr/likeminds-authentication/constants"
	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/utils"
)

func CreateOTM(api_key string) (*constants.OnboardingTokenMeta, error) {
	otm := &constants.OnboardingTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error

	otmClaims := jwt.MapClaims{}
	otmClaims["access_uuid"] = otm.AccessUuid
	otmClaims["exp"] = otm.AccessTokenExpires
	otmClaims["api_key"] = api_key
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, otmClaims)
	otm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return otm, nil
}

func CreateVTM(apiKey string, emailId string, mobileNo string, countryCode string) (*constants.VerifyTokenMeta, error) {
	vtm := &constants.VerifyTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	var err error

	vtmClaims := jwt.MapClaims{}
	vtmClaims["access_uuid"] = vtm.AccessUuid
	vtmClaims["exp"] = vtm.AccessTokenExpires
	vtmClaims["api_key"] = apiKey

	if emailId != "" {
		vtmClaims["email_id"] = emailId

	} else if mobileNo != "" && countryCode != "" {
		vtmClaims["mobile_no"] = mobileNo
		vtmClaims["country_code"] = countryCode
	} else {
		return nil, errors.New("Email ID or Mobile no should be present!")
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, vtmClaims)
	vtm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return vtm, nil
}

// CreateLTMAndRTM is used to create login and refresh token meta
func CreateLTMAndRTM(userUniqueID string, api_key string, token_expiry_beta int64, isGuestUser bool) (*constants.LoginTokenMeta, *constants.RefreshTokenMeta, error) {

	isBeta := environment.GoDotEnvVariable("BETA_ENVIRONMENT")

	// Setting LTM token expiry to 15 minutes for Prod
	LTMTokenExpiryTime := time.Duration(PROD_AUTH_TOKEN_EXPIRY)

	if isBeta == "true" {

		// Setting default LTM token expiry to 60 minutes for Beta
		if token_expiry_beta <= 0 {
			token_expiry_beta = BETA_AUTH_TOKEN_EXPIRY
		}

		LTMTokenExpiryTime = time.Duration(token_expiry_beta)
	}

	ltm := &constants.LoginTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * LTMTokenExpiryTime).Unix(),
	}

	rtm := &constants.RefreshTokenMeta{
		RefreshUuid:         uuid.NewV4().String(),
		RefreshTokenExpires: time.Now().Add(time.Hour * 24 * 31).Unix(),
	}

	var err error
	//Creating login token meta
	//os.Setenv("ACCESS_SECRET", "JWT_SECRET") //this should be in an env file
	ltmClaims := jwt.MapClaims{}
	ltmClaims["access_uuid"] = ltm.AccessUuid
	ltmClaims["user_unique_id"] = userUniqueID
	ltmClaims["exp"] = ltm.AccessTokenExpires
	ltmClaims["api_key"] = api_key
	ltmClaims["is_guest"] = isGuestUser

	bytesData, err := json.Marshal(ltmClaims)

	if err == nil {
		encryptedData := utils.Encrypt(bytesData)

		if encryptedData != "" {
			ltmClaims = jwt.MapClaims{"data": encryptedData}
		}
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ltmClaims)
	ltm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable("ACCESS_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	//Creating refresh token meta
	rtmClaims := jwt.MapClaims{}
	rtmClaims["refresh_uuid"] = rtm.RefreshUuid
	rtmClaims["user_unique_id"] = userUniqueID
	rtmClaims["exp"] = rtm.RefreshTokenExpires
	rtmClaims["api_key"] = api_key
	rtmClaims["is_guest"] = isGuestUser

	bytesData, err = json.Marshal(rtmClaims)

	if err == nil {
		encryptedData := utils.Encrypt(bytesData)

		if encryptedData != "" {
			rtmClaims = jwt.MapClaims{"data": encryptedData}
		}
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtmClaims)
	rtm.RefreshToken, err = rt.SignedString([]byte(environment.GoDotEnvVariable("ACCESS_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	return ltm, rtm, nil
}

// ExtractToken to extract token from headers or return string value
func ExtractToken(bearerToken string) string {
	//Normally Authorization the_token_xxx
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	//Return simple token string
	return bearerToken
}

// VerifyToken to verify token
func VerifyToken(bearerToken string) (*jwt.Token, error) {
	tokenString := ExtractToken(bearerToken)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//os.Setenv("ACCESS_SECRET", "JWT_SECRET")
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(environment.GoDotEnvVariable("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractOTM is used to return OTM and check if bearer token is valid or not
func ExtractOTM(bearerToken string) (*constants.OnboardingTokenMeta, error) {
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

		apiKey, _ := claims["api_key"].(string)

		return &constants.OnboardingTokenMeta{
			AccessUuid: accessUuid,
			ApiKey:     apiKey,
		}, nil

	}

	return nil, err
}

// ExtractVTM is used to return VTM and check if bearer token is valid or not
func ExtractVTM(bearerToken string) (*constants.VerifyTokenMeta, error) {
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

		apiKey, _ := claims["api_key"].(string)
		emailId, _ := claims["email_id"].(string)
		mobileNo, _ := claims["mobile_no"].(string)
		countryCode, _ := claims["country_code"].(string)

		return &constants.VerifyTokenMeta{
			AccessUuid:  accessUuid,
			ApiKey:      apiKey,
			EmailID:     emailId,
			MobileNo:    mobileNo,
			CountryCode: countryCode,
		}, nil

	}

	return nil, err
}

// ExtractLTM is used to return LTM and check if bearer token is valid or not
func ExtractLTM(bearerToken string) (*constants.LoginTokenMeta, error) {
	var claims jwt.MapClaims

	token, err := VerifyToken(bearerToken)
	if err != nil {
		return nil, err
	}
	jwt_claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		encryptedData, ok := jwt_claims["data"].(string)
		if ok {
			claimsBytes := utils.Decrypt(encryptedData)
			json.Unmarshal(claimsBytes, &claims)
		} else {
			claims = jwt_claims
		}

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
		isGuest, ok := claims["is_guest"].(bool)
		if !ok {
			return nil, errors.New("is_guest is empty")
		}
		apiKey, _ := claims["api_key"].(string)
		return &constants.LoginTokenMeta{
			AccessUuid:         accessUuid,
			UserUniqueID:       userUniqueID,
			AccessTokenExpires: atExpires,
			ApiKey:             apiKey,
			IsGuest:            isGuest,
		}, nil
	}
	return nil, err
}

// ExtractRTM is used to return RTM and check if bearer token is valid or not
func ExtractRTM(bearerToken string) (*constants.RefreshTokenMeta, error) {
	var claims jwt.MapClaims

	token, err := VerifyToken(bearerToken)
	if err != nil {
		return nil, err
	}
	jwt_claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		encryptedData, ok := jwt_claims["data"].(string)
		if ok {
			claimsBytes := utils.Decrypt(encryptedData)
			json.Unmarshal(claimsBytes, &claims)
		} else {
			claims = jwt_claims
		}

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
		isGuest, ok := claims["is_guest"].(bool)
		if !ok {
			return nil, errors.New("is_guest is empty")
		}
		apiKey, _ := claims["api_key"].(string)
		return &constants.RefreshTokenMeta{
			RefreshUuid:         refreshUuid,
			UserUniqueID:        userUniqueID,
			RefreshTokenExpires: rtExpires,
			ApiKey:              apiKey,
			IsGuest:             isGuest,
			RefreshToken:        ExtractToken(bearerToken),
		}, nil
	}
	return nil, err
}
