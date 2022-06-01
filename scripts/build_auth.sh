#!/bin/bash
source /home/ubuntu/.profile
echo $PATH
cd /home/ubuntu/LikeMinds-Authentication/
git pull origin development
go build .
