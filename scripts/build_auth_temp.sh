#!/bin/bash
source /home/ubuntu/.profile
cd /home/ubuntu/LikeMinds-Authentication/
branch_name="$(git rev-parse --abbrev-ref HEAD)"
if [ "$branch_name" = "development" ]
then
  curl -o .env https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/beta-environment
  git pull origin development
else
  curl -o .env https://prod-likeminds-media.s3.ap-south-1.amazonaws.com/environment/prod-environment
  git pull origin master
fi
go build .
