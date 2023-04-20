package utility

import (
	"github.com/nateshr/likeminds-authentication/utils"
)

func AuthenticateAPIKeyInternally(headers map[string]interface{}, api_key string) (map[string]string, error) {

	var response map[string]string

	if api_key == "" {
		return response, nil
	}

	// Internally call /api/sdk/authenticate
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, SDKAuthenticateEndPoint, utils.GETRequest, headers, nil, nil)

	if err != nil {
		return nil, err
	}

	dataResponse := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)

	// Parse response
	if dataResponse != nil {
		var communityId string

		if dataResponse[ParamCommunityID] != nil {
			communityId = utils.ParseInterfaceToString(dataResponse[ParamCommunityID])
		}

		response = map[string]string{
			ParamCommunityID: communityId,
		}

	}

	return response, nil
}
