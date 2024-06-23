#!/bin/bash
APP_NAME=$(<.Docker-name)
id=$(sudo docker restart $APP_NAME)
sudo docker logs -f $id