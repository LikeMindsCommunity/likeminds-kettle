package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/myesui/uuid"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func CreateOTM(apiKey string) (*constants.OnboardingTokenMeta, error) {
	otm := &constants.OnboardingTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * 15).Unix(),
	}

	var err error

	otmClaims := jwt.MapClaims{}
	otmClaims[TokemAccessUUID] = otm.AccessUuid
	otmClaims[TokenExp] = otm.AccessTokenExpires
	otmClaims[TokenAPIKey] = apiKey
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, otmClaims)
	otm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable(utils.EnvAccessSecret)))
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
	vtmClaims[TokemAccessUUID] = vtm.AccessUuid
	vtmClaims[TokenExp] = vtm.AccessTokenExpires
	vtmClaims[TokenAPIKey] = apiKey

	if emailId != "" {
		vtmClaims["email_id"] = emailId

	} else if mobileNo != "" && countryCode != "" {
		vtmClaims["mobile_no"] = mobileNo
		vtmClaims["country_code"] = countryCode
	} else {
		return nil, errors.New("Email ID or Mobile no should be present!")
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, vtmClaims)
	vtm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable(utils.EnvAccessSecret)))
	if err != nil {
		return nil, err
	}
	return vtm, nil
}

// CreateLTMAndRTM is used to create login and refresh token meta
func CreateLTMAndRTM(userUniqueID string, apiKey string, tokenExpiryBeta int64, rtmExpiryBeta int64, isGuestUser bool, deviceID string,
) (*constants.LoginTokenMeta, *constants.RefreshTokenMeta, error) {

	isBeta := environment.GoDotEnvVariable("BETA_ENVIRONMENT")

	// LTM & RTM token expiry
	LTMTokenExpiryTime := time.Duration(PROD_AUTH_TOKEN_EXPIRY)
	RTMTokenExpiryTime := time.Duration(time.Hour * REFRESH_TOKEN_EXPIRY)

	if isBeta == "true" {

		// Setting default LTM token expiry to 60 minutes for Beta
		if tokenExpiryBeta <= 0 {
			tokenExpiryBeta = BETA_AUTH_TOKEN_EXPIRY
		}

		LTMTokenExpiryTime = time.Duration(tokenExpiryBeta)

		if rtmExpiryBeta > 0 {
			RTMTokenExpiryTime = time.Duration(time.Minute * time.Duration(rtmExpiryBeta))
		}
	}

	ltm := &constants.LoginTokenMeta{
		AccessUuid:         uuid.NewV4().String(),
		AccessTokenExpires: time.Now().Add(time.Minute * LTMTokenExpiryTime).Unix(),
	}

	rtm := &constants.RefreshTokenMeta{
		RefreshUuid:         uuid.NewV4().String(),
		RefreshTokenExpires: time.Now().Add(RTMTokenExpiryTime).Unix(),
	}

	var err error
	//Creating login token meta
	ltmClaims := jwt.MapClaims{}
	ltmClaims[TokemAccessUUID] = ltm.AccessUuid
	ltmClaims[TokenUserUniqueId] = userUniqueID
	ltmClaims[TokenAPIKey] = apiKey
	ltmClaims[TokenIsGuest] = isGuestUser
	ltmClaims[TokenDeviceID] = deviceID

	bytesData, err := json.Marshal(ltmClaims)

	if err == nil {
		encryptedData := utils.Encrypt(bytesData)

		if encryptedData != "" {
			ltmClaims = jwt.MapClaims{"data": encryptedData}
		}
	}

	ltmClaims[TokenExp] = ltm.AccessTokenExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ltmClaims)
	ltm.AccessToken, err = at.SignedString([]byte(environment.GoDotEnvVariable(utils.EnvAccessSecret)))
	if err != nil {
		return nil, nil, err
	}

	//Creating refresh token meta
	rtmClaims := jwt.MapClaims{}
	rtmClaims["refresh_uuid"] = rtm.RefreshUuid
	rtmClaims[TokenUserUniqueId] = userUniqueID
	rtmClaims[TokenAPIKey] = apiKey
	rtmClaims[TokenIsGuest] = isGuestUser
	rtmClaims[TokenDeviceID] = deviceID

	bytesData, err = json.Marshal(rtmClaims)

	if err == nil {
		encryptedData := utils.Encrypt(bytesData)

		if encryptedData != "" {
			rtmClaims = jwt.MapClaims{"data": encryptedData}
		}
	}

	rtmClaims[TokenExp] = rtm.RefreshTokenExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtmClaims)
	rtm.RefreshToken, err = rt.SignedString([]byte(environment.GoDotEnvVariable(utils.EnvAccessSecret)))
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
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(environment.GoDotEnvVariable(utils.EnvAccessSecret)), nil
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
		accessUuid, ok := claims[TokemAccessUUID].(string)

		if !ok {
			return nil, err
		}

		apiKey, _ := claims[TokenAPIKey].(string)

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
		accessUuid, ok := claims[TokemAccessUUID].(string)

		if !ok {
			return nil, err
		}

		apiKey, _ := claims[TokenAPIKey].(string)
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

		atExpires, accessUUID, userUniqueId, apiKey, isGuest, err, deviceID := extractParamsFromLTM(jwt_claims, claims)
		if err != nil {
			return nil, err
		}

		return &constants.LoginTokenMeta{
			AccessUuid:         accessUUID,
			UserUniqueID:       userUniqueId,
			AccessTokenExpires: atExpires,
			ApiKey:             apiKey,
			IsGuest:            isGuest,
			DeviceID:           deviceID,
		}, nil
	}
	return nil, err
}

func extractParamsFromLTM(jwtClaims jwt.MapClaims, claims jwt.MapClaims,
) (int64, string, string, string, bool, error, string) {

	atExpires, accessUUID, userUniqueId, apiKey, isGuest, deviceID := int64(0), "", "", "", false, ""

	atExpiresFloat, ok := jwtClaims[TokenExp].(float64)
	if !ok {
		atExpiresFloat, _ = claims[TokenExp].(float64)
	}

	if atExpiresFloat != 0 {
		atExpires = int64(atExpiresFloat)
	}

	if atExpires == 0 {
		return atExpires, accessUUID, userUniqueId, apiKey, isGuest, errors.New("exp is empty"), deviceID
	} else if atExpires < time.Now().Unix() {
		return atExpires, accessUUID, userUniqueId, apiKey, isGuest, errors.New("LTM expired!"), deviceID
	}

	accessUUID, ok = claims[TokemAccessUUID].(string)
	if !ok {
		return atExpires, accessUUID, userUniqueId, apiKey, isGuest, errors.New("access_uuid is empty"), deviceID
	}

	userUniqueId, ok = claims[TokenUserUniqueId].(string)
	if !ok {
		return atExpires, accessUUID, userUniqueId, apiKey, isGuest, errors.New("user_unique_id is empty"), deviceID
	}
	isGuest, ok = claims[TokenIsGuest].(bool)
	if !ok {
		return atExpires, accessUUID, userUniqueId, apiKey, isGuest, errors.New("is_guest is empty"), deviceID
	}
	deviceID, ok = claims[TokenDeviceID].(string)

	apiKey, _ = claims[TokenAPIKey].(string)

	return atExpires, accessUUID, userUniqueId, apiKey, isGuest, nil, deviceID
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

		rtExpires, refreshUuid, userUniqueID, apiKey, isGuest, err, deviceID := extractParamsFromRTM(jwt_claims, claims)
		if err != nil {
			return nil, err
		}

		return &constants.RefreshTokenMeta{
			RefreshUuid:         refreshUuid,
			UserUniqueID:        userUniqueID,
			RefreshTokenExpires: rtExpires,
			ApiKey:              apiKey,
			IsGuest:             isGuest,
			RefreshToken:        ExtractToken(bearerToken),
			DeviceID:            deviceID,
		}, nil
	}
	return nil, err
}

func extractParamsFromRTM(jwtClaims jwt.MapClaims, claims jwt.MapClaims,
) (int64, string, string, string, bool, error, string) {

	rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, deviceID := int64(0), "", "", "", false, ""

	rtExpiresFloat, ok := jwtClaims[TokenExp].(float64)
	if !ok {
		rtExpiresFloat, _ = claims[TokenExp].(float64)
	}
	if rtExpiresFloat != 0 {
		rtExpires = int64(rtExpiresFloat)
	}
	if rtExpires == 0 {
		return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, errors.New("exp is empty"), deviceID
	} else if rtExpires < time.Now().Unix() {
		return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, errors.New("RTM expired!"), deviceID
	}

	refreshUUID, ok = claims["refresh_uuid"].(string)
	if !ok {
		return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, errors.New("access_uuid is empty"), deviceID
	}

	userUniqueId, ok = claims[TokenUserUniqueId].(string)
	if !ok {
		return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, errors.New("user_unique_id is empty"), deviceID
	}

	isGuest, ok = claims[TokenIsGuest].(bool)
	if !ok {
		return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, errors.New("is_guest is empty"), deviceID
	}

	deviceID, _ = claims[TokenDeviceID].(string)

	apiKey, _ = claims[TokenAPIKey].(string)

	return rtExpires, refreshUUID, userUniqueId, apiKey, isGuest, nil, deviceID
}
