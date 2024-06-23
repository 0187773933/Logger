#!/bin/bash

sudo chown -R morphs:morphs /home/morphs/SAVE_FILES

GIT_URL="https://github.com/0187773933/Logger.git"
HASH_FILE="/home/morphs/git.hash"
REMOTE_HASH=$(git ls-remote "$GIT_URL" HEAD | awk '{print $1}')

if [ -f "$HASH_FILE" ]; then
	STORED_HASH=$(sudo cat "$HASH_FILE")
else
	STORED_HASH=""
fi

if [ "$REMOTE_HASH" == "$STORED_HASH" ]; then
	echo "No New Updates Available"
	cd /home/morphs/DockerBuild
	LOG_LEVEL=debug exec /home/morphs/DockerBuild/server "$@"
else
	echo "New updates available. Updating and Rebuilding Go Module"
	echo "$REMOTE_HASH" | sudo tee "$HASH_FILE"
	cd /home/morphs
	sudo rm -rf /home/morphs/DockerBuild
	git clone "$GIT_URL"
	sudo chown -R morphs:morphs /home/morphs/DockerBuild
	cd /home/morphs/DockerBuild
	cp -r /home/morphs/SAVE_FILES/ /home/morphs/DockerBuild/SAVE_FILES/
	sudo chown -R morphs:morphs /home/morphs/DockerBuild/SAVE_FILES/
	/usr/local/go/bin/go mod tidy
	GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o /home/morphs/DockerBuild/server
	LOG_LEVEL=debug exec /home/morphs/DockerBuild/server "$@"
fi