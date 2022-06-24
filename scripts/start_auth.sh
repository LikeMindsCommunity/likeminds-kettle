#!/bin/bash
source const.sh
if [ "$BETA_ENVIRONMENT" = true ]
then
  export $(curl https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/beta-environment)
else
  export $(curl https://prod-likeminds-media.s3.ap-south-1.amazonaws.com/environment/prod-environment)
  echo $BETA_ENVIRONMENT
  echo $GIN_MODE
fi
sudo systemctl restart auth
