#!/bin/sh

# This script will build the docker image for a project based on build arguments
# provided by the project env file and a user-level credential file.
# If user-level credential file does not exist, script will attempt to create it
# and notify user.
# TODO:
#   - Need to tweak docker build to tag using standardized naming convention.
#   - Need to add clean steps to clean up the build image to save space.
#     (may need to set as optional via flag.)
#   - Might also want to have script to push built image to repository if build
#     if successful.

set -e -u

PROJECT_ENV='./files/.env'
BB_CREDS="$HOME/.bb-creds"

FILE_TEXT=$(cat <<"EOF"
# Create an APP PASSWORD in bitbucket and put credentials here.
# Please, do not use your normal account password!
BB_UN=
BB_PW=
EOF
)
#Text Colors
RED='\033[0;31;1m'
YELLOW='\033[0;33m'
RC='\033[0m' # Reset color

# Load project env file so get some build info.
if [ -f $PROJECT_ENV ]; then
    . $PROJECT_ENV
    echo "Using project .env file"
else
    echo "Project .env file not found. Expected location: $PROJECT_ENV"
    echo "We don't have enough info to continue with build, aborting."
    exit
fi

# Get bitbucket credentials so we can pull down private repos.
if [ -f $BB_CREDS ]; then
    . $BB_CREDS
    # We need to make sure the variables in the file are set, in case the
    # script created a file and it has not been populated.
    VAR_UNSET=0
    if [ -z "$BB_UN" ]; then
        echo "${RED}Username not set. Please update in $BB_CREDS${RC}"
        VAR_UNSET=1
    fi
    if [ -z $BB_PW ]; then
        echo "${RED}Password not set.  Please update in $BB_CREDS${RC}"
        VAR_UNSET=1
    fi
    if [ "$VAR_UNSET" != 0 ]; then
        echo "Aborting."
        exit
    fi
else
    # We need to try and create file for bitbucket credentials.
    echo "File $BB_CREDS not found.\nAttempting to create file..."
    echo "$FILE_TEXT" > $BB_CREDS
    chmod 600 $BB_CREDS
    if [ -f $BB_CREDS ]; then
        echo "${YELLOW}File $BB_CREDS was created.\nYou need to add a bitbucket"\
            "app password to this file in order to successfully build a docker" \
            "file.${RC}"
    else
        echo "${RED}File $BB_CREDS could not be automatically created.  Please" \
            "create this file with the following text:\n${RC}"
        echo "$FILE_TEXT"
    fi
    # No matter what, we exit here so user can edit or create file.
    exit
fi

# Set any additional build parameters not provided by files.
APPNAME=`basename "$PWD"`

# Build!
docker build \
  --build-arg BB_UN=$BB_UN \
  --build-arg BB_PW=$BB_PW \
  --build-arg SERVICE_PORT=$PORT \
  --build-arg APPNAME=$APPNAME \
  -t $APPNAME .


# Clean!
# (tbd)
