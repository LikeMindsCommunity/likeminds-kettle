#!/bin/bash
source /home/ubuntu/.profile
cd /home/ubuntu/LikeMinds-Authentication/
branch_name="$(git rev-parse --abbrev-ref HEAD)"
if [ "$branch_name" = "master" ]
then
  curl -o ./.env https://prod-likeminds-media.s3.ap-south-1.amazonaws.com/environment/kettle-prod-dot-env-public
  git pull origin master
else
  curl -o ./.env https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/Kettle-Beta-Dot-Env/.env
  git pull origin "$branch_name"
fi
set -a
source ./.env
docker compose -f ./prometheus/docker-compose-prometheus.yml up -d --build