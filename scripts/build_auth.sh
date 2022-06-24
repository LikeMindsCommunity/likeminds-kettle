#!/bin/bash
source /home/ubuntu/.profile
echo $PATH
cd /home/ubuntu/LikeMinds-Authentication/
BETA_ENVIRONMENT=false
if [ "$BETA_ENVIRONMENT" = true ]
then
  export $(curl https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/beta-environment)
  git pull origin development
else
  export $(curl https://prod-likeminds-media.s3.ap-south-1.amazonaws.com/environment/prod-environment)
  git pull origin master
fi
go build .
