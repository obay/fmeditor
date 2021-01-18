#!/bin/bash

# If you want to use Sentry to capture your exceptions, you can define the following
# environment variables
export SENTRY_DSN="use your own Sentry DSN"
export SENTRY_ENVIRONMENT="Dev"
export SENTRY_RELEASE="fmeditor"

go build
./fmeditor --author "Ahmad Obay" --draft true --rootfolder "testfolder/"
