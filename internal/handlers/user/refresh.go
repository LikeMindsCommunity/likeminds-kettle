package user

import (
	"encoding/json"
	"net/http"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/token"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	TokenExpiryBeta int `json:"token_expiry_beta,omitempty"`
}

// Refresh to generate new LTM and RTM tokens
func Refresh(c *gin.Context) {

	// Parse Request
	request := RefreshRequest{}
	jsonData, _ := c.GetRawData()
	json.Unmarshal(jsonData, &request)

	//Check if request has RTM token or not
	currentRTM, ok := c.MustGet(constants.ParamRTM).(*constants.RefreshTokenMeta)
	if !ok {
		utils.GeneralAPIError(c, utils.ErrorInvalidRTM)
		return
	}

	ltmExpiry := token.BETA_AUTH_TOKEN_EXPIRY
	if request.TokenExpiryBeta > 0 {
		ltmExpiry = request.TokenExpiryBeta
	}

	//Create login and refresh token meta from ltm
	ltm, _, err := token.CreateLTMAndRTM(currentRTM.UserUniqueID, currentRTM.ApiKey, int64(ltmExpiry), token.DEFAULT_TOKEN_EXPIRY, currentRTM.IsGuest, currentRTM.DeviceID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	//Generate Token Object
	token := map[string]interface{}{
		constants.ParamAccessToken:  ltm.AccessToken,
		constants.ParamRefreshToken: currentRTM.RefreshToken,
	}

	//Generate Response
	utils.GenerateResponse(c, token, false)
}
