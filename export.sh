#!/usr/bin/env bash

set -euo pipefail

if [ $# -ne 1 ] ; then

    echo "USAGE $0 platform"
    exit 1

fi

platform_pattern="^(test|stage|live)+$"

if [[ ! $1 =~ $platform_pattern ]] ; then

    echo "Invalid release-name: $1 does not match $platform_pattern regex"
    exit 1
fi

platform=$1
package=models
package_dir=models
export_dir=export

# Get the API id from the API name using jq command-line json tool. See https://stedolan.github.io/jq

api_id=$( aws apigateway get-rest-apis | jq  -r '.items[] | select(.name == "Sample-API-'${platform}'") | .id' )

# Export a JSON-format Swagger API definition from the AWS Gateway API

aws apigateway get-export --rest-api-id $api_id  --stage-name $platform --export-type swagger $export_dir/$platform.json

# The schema-generator executable can be created from here: https://github.com/merlincox/generate
# It generates a Go source file of struct declarations from the Swagger API definition file

schema_generator=$(which schema-generator)

if [ ! -z "$schema_generator" ]; then

    $schema_generator -p $package -nsk $export_dir/$platform.json > $export_dir/$platform.go

    go fmt $export_dir/$platform.go

fi
