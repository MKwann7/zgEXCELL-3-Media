#! /bin/bash

clear
reset

BUILD=local
APP_NAME=zgexcell_media

project_path="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "project_path = ${project_path}"

docker-compose -f docker/docker-compose.local.yml --project-name ${APP_NAME} up --build