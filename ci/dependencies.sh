#!/bin/bash

set -eu

# Add python pip and bash
apt-get update
apt-get install -y python3-pip

# Install docker-compose via pip
pip3 install --no-cache-dir docker-compose
docker-compose -v

# install psql and other tools to confirm db dependency is up
apt-get install -y postgresql-client-common postgresql-client