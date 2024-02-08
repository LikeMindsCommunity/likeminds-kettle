#!/bin/bash
branch_name="$(git rev-parse --abbrev-ref HEAD)"
curl -o ./.env https://likeminds-configs-prod.s3.ap-south-1.amazonaws.com/application-dot-envs-prod/kettle-prod/kettle-prod-dot-env-public
git pull origin master
#if [ "$branch_name" = "master" ]; then
#  curl -o ./.env https://likeminds-configs-prod.s3.ap-south-1.amazonaws.com/application-dot-envs-prod/kettle-prod/kettle-prod-dot-env-public
#  git pull origin master
#else
#  curl -o ./.env https://likeminds-configs-beta.s3.ap-south-1.amazonaws.com/application-dot-envs-beta/kettle-beta/kettle-beta-dot-env-public
#  git pull origin "$branch_name"
#fi
set -a
source ./.env
docker compose -f ./monitoring/prometheus/docker-compose-prometheus.yml up -d --build
