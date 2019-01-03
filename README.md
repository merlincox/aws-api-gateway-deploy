## aws-api-gateway-deploy

This repo contains a Bash deploy script and a CloudFormation template for deploying a serverless API implemented as a
 AWS API Gateway served by a simple Golang lambbda

The template

* Creates the lambda and passes an environment to it
* Creates the API Gateway with an endpoint linked to the lambda
* Creates a subdomain mapped to the API

Altogether, the deployment script enables you to create to `https://{your-sub-domain}.{your-domain}/status` endpoint,
 an example environment for the API (in the supplied example git tag, branch, and commit information and a 
'platform' variable). Redirection of HTTP to HTTPS and CloudWatch logging are automatically supplied by the gateway.

### Prerequsisites

* A domain with a hosted-zone record in AWS Route 53
* A SSL certificate for that domain in the AWS Certificate Manager. Note that this certificate has to be in the 
`us-east-1` AWS region because it is deployed to CloudFront.
* The AWS command line interface `aws` installed and suitably set up with credentials for your AWS account
* `go` installed
* `glide` (a go dependency manager) installed
* `git` installed
* `jq` installed (`jq` is a very useful command-line tool for manipulating JSON. See https://stedolan.github.io/jq.)

### Deployment

The deployment script usage is:
 
 `./deploy.sh subdomain_base domain [platform]`

`platform` is intended to be 'test', 'stage' or 'test' ('test' is the default)

'-test' and '-stage' are appended to the `subdomain_base`

Thus: 

`./deploy.sh my-api my-domain.com` will deploy an API at `https://my-api-test.my-domain.com`

`./deploy.sh my-api my-domain.com stage` will deploy an API at `https://my-api-stage.my-domain.com`

`./deploy.sh my-api my-domain.com live` will deploy an API at `https://my-api.my-domain.com`

Uncommited code cannot be deployed, and live deploys have these additional checks:

* code must be on the master branch
* code must be sync with the remote origin
* code must be exactly on a tag of the form 

Lastly, live deploys have a confirmation prompt.

### Exporting Swagger JSON and models

The API definition YAML includes a Swagger definition for the API.

This can be exported using the `export.sh` script.

In addition a schema-generator executable can be created from here: https://github.com/merlincox/generate

If this is added to the system path, the `export.sh` will also generate Go structs for the API and optionally replace 
the pkg/models/api.go file if that is out of sync with the API. (Therefore any additionmal models which do not feature 
directly in the API should be placed in the pkg/models/models.go file).

