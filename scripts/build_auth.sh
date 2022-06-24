#!/bin/bash
source const.sh
source /home/ubuntu/.profile
cd /home/ubuntu/LikeMinds-Authentication/
if [ "$BETA_ENVIRONMENT" = true ]
then
  git pull origin development
else
  git pull origin master
fi
go build .
