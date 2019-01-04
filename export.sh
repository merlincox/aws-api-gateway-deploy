#!/usr/bin/env bash

schema_generator=$(which schema-generator)

set -euo pipefail

cd "$(dirname "$0")"

if [[ $# -ne 1 ]] ; then
    echo "USAGE: $(basename $0) platform" >&2
    exit 1
fi

platform_pattern="^(test|stage|live)$"

if [[ ! $1 =~ $platform_pattern ]] ; then
    echo "Invalid release-name: $1 does not match ${platform_pattern} regex" n>&2
    exit 1
fi

platform=$1

package=models
package_dir=pkg/models
export_dir=export

export_json=${export_dir}/swagger_${platform}.json
exported_go=${export_dir}/auto_${platform}.go
current_go=${package_dir}/api.go

if [[ ! -d ${export_dir} ]]; then
    mkdir ${export_dir}
fi

# Get the API id from the API name using jq command-line json tool. See https://stedolan.github.io/jq

api_id=$( aws apigateway get-rest-apis | jq  -r '.items[] | select(.name == "Sample-API-'${platform}'") | .id' )

# Export a JSON-format Swagger API definition from the AWS Gateway API

aws apigateway get-export --rest-api-id $api_id  --stage-name ${platform} --export-type swagger ${export_json}

# The schema-generator executable can be created from here: https://github.com/merlincox/generate
# It generates a Go source file of struct declarations from the Swagger API definition file

if [[ ! -z "${schema_generator}" ]]; then

    ${schema_generator} -p ${package} -nsk ${export_json} > ${exported_go}

    go fmt ${exported_go}

    if diff ${exported_go} ${current_go}; then
        echo "Exported API models for ${platform} match current API models at ${current_go}"
    else
        echo
        read -p "Update current API models? :" -n 1 -r reply
        echo
        case "$reply" in
            y|Y ) cp -f ${current_go} $export_dir/api.go_old
                  cp -f ${exported_go} ${current_go}
                  ;;
            * )   echo Not overwritten
                  ;;
        esac
    fi
    mv ${exported_go} ${exported_go}_new
else
    echo "No schema generator found"
fi
