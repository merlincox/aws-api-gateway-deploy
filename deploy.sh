#!/usr/bin/env bash

cd "$( dirname "$0" )"

for cmd in "go aws git jq"; do

    if [[ -z "$(which ${cmd})" ]]; then
        echo "${cmd} is required to run this script."  >&2
        exit 1
    fi

done

status=$(git status --porcelain)

if [[ ! -z "$status" ]]; then
    echo "There are uncommitted changes"
    exit 1
fi

git_commit=$(git rev-parse --short=16 HEAD)
timestamp=$(date +"%Y-%m-%dT%H:%M:%S")

if [[ $# -lt 2 ]] ||  [[ $# -gt 3 ]]; then

    echo "USAGE $(basename $0) subdomain_base domain [platform]" >&2
    exit 1

fi

subdomain_base=$1
domain=$2
platform="$3"

if [[ -z "${platform}" ]]; then
    platform=test
fi

set -euo pipefail

template_yaml=api.yaml

if [[ ! -f ${template_yaml} ]]; then
   echo "${template_yaml} not found" >&2
   exit 1
fi

platform_regex="^(test|stage|live)$"
if [[ ! "${platform}" =~ $platform_regex ]]
then
    echo "Platform must live, test or stage"
    exit 1
fi

subdomain_regex="^[a-z0-9-]+$"
if [[ ! "${subdomain_base}" =~ $subdomain_regex ]]
then
    echo "${subdomain_base} does not match pattern ${subdomain_regex}" >&2
    exit 1
fi

if [[ "${platform}" != "live" ]]; then
    subdomain=${subdomain_base}-${platform}
else
    subdomain=${subdomain_base}
fi

domain_zone_id=$(aws route53 list-hosted-zones | jq -r ".HostedZones[] | select(.Name == \"${domain}.\")| .Id" | cut -f3 -d "/")

if [[ -z "${domain_zone_id}" ]]; then
    echo "No hosted-zone record was found for ${domain} domain" >&2
    exit 1
fi

custom_domain="${subdomain}.${domain}"

# NB Cloud Front certs must be in us-east-1 region
certificate_arns=$(aws acm list-certificates --region us-east-1 | jq -r ".CertificateSummaryList[] | select(.DomainName == \"${custom_domain}\") | .CertificateArn")

if [[ -z "${certificate_arns}" ]]; then
     certificate_arns=$(aws acm list-certificates --region us-east-1 | jq -r ".CertificateSummaryList[] | select(.DomainName == \"*.${domain}\") | .CertificateArn")

    if [[ -z "${certificate_arns}" ]]; then
       echo "No SSL certificate was found for ${custom_domain} or *.${domain} patterns in us-east-1" >&2
        exit 1
    fi
fi

for arn in ${certificate_arns}
do

  echo $arn

done
exit 1

git_tag="untagged"

if git describe --tags >/dev/null 2>/dev/null; then
   git_tag=$(git describe --tags)
fi

git_branch=$(git rev-parse --abbrev-ref HEAD)
git_commit=$(git rev-parse --short=16 HEAD)

tag_pattern="^v[0-9]+\.[0-9]+\.[0-9]+$"
if [[ $platform == live ]] ; then

    if [[ "${git_branch}" != "master" ]]; then

       echo "Live deployments must be on the master branch" >&2
       exit 1
    fi

    git fetch
    git_diff=$(git diff origin/master)

    if [[ ! -z "${git_diff}" ]]; then
        echo "Live deployments must be in sync with origin/master" >&2
        exit 1
    fi

    if [[ ! "${git_tag}" =~ ${tag_pattern} ]] ; then

        echo "Live deployments must be tagged with vN.N.N: ${git_tag} does not match" >&2
        exit 1
    fi

    read -p "About to deploy ${git_tag} to live. Confirm? :" -n 1 -r reply
    echo
    case "$reply" in
        y|Y ) echo "Proceeding with deployment";;
        * )   echo "Cancelling deployment" >&2
              exit 1
              ;;
    esac
fi

cf_bucket=cf-api-import-$(date +"%y%m%d%H%M")
cf_stack=api-stack-${platform}

package_yml=$(mktemp /tmp/XXXXXXX.yaml)

go mod download

if   go test ./...
then echo "Tests passed"
else echo "Tests failed"
     exit 1
fi

executable=bin/api
env GOOS=linux go build -o ${executable} api/main.go
chmod +x ${executable}

bucket_created=0

function cleanup {

    if [[ $bucket_created -eq 1 ]] ; then
        aws s3 rm s3://${cf_bucket} --recursive
        aws s3 rb s3://${cf_bucket}
    fi

    if [[ -f ${package_yml} ]]; then
        rm ${package_yml}
    fi
}

trap cleanup EXIT

echo Running at $(date +"%H:%M %d/%m/%y")

aws s3api create-bucket --bucket $cf_bucket --create-bucket-configuration LocationConstraint=$(aws configure get region)

bucket_created=1

aws cloudformation package \
       --template-file ${template_yaml} \
       --s3-bucket $cf_bucket \
       --output-template-file $package_yml

aws cloudformation deploy \
       --template-file $package_yml \
       --stack-name $cf_stack \
       --capabilities CAPABILITY_IAM \
       --parameter-overrides Platform="${platform}" Commit="${git_commit}" \
           CustomDomain="${custom_domain}" HostedZone="${domain_zone_id}" \
           Release="${git_tag}" Branch="${git_branch}" CertificateArn="${certificate_arn}"

