#!/bin/bash
source /home/ubuntu/.profile
echo $PATH
cd /home/ubuntu/LikeMinds-Authentication/
if [ "$BETA_ENVIRONMENT" = true]
  curl $(export https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/beta-environment)
  git pull origin development
else
  curl $(export https://prod-likeminds-media.s3.ap-south-1.amazonaws.com/environment/prod-environment)
  git pull origin master
go build .
