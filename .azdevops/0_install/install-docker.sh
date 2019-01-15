#!/bin/bash

set -o nounset

echo; echo "Removing any previously installed Docker versions..."
sudo apt-get remove --purge docker*
sudo apt-get autoremove --purge
sudo apt-get autoclean

echo; echo "Accessing Docker repository..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository \
	"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
	$(lsb_release -cs) \
	stable"

echo; echo "Installing Docker..."
sudo apt-get update
sudo apt-get install docker-ce=$DOCKER_VERSION

sudo service docker restart
