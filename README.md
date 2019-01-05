## aws-api-gateway-deploy

This repo contains a Bash deploy script and a CloudFormation template for deploying a serverless API implemented as a
 AWS API Gateway served by a simple Golang lambbda

The template

* Creates the lambda and passes an environment to it
* Creates the API Gateway with two sample endpoints linked to the lambda
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

Examples: 

`./deploy.sh my-api my-domain.com` will deploy an API at `https://my-api-test.my-domain.com`

`./deploy.sh my-api my-domain.com stage` will deploy an API at `https://my-api-stage.my-domain.com`

`./deploy.sh my-api my-domain.com live` will deploy an API at `https://my-api.my-domain.com`

Uncommited code cannot be deployed, and live deploys have these additional checks:

* code must be on the master branch
* code must be sync with the remote origin
* code must be exactly on a tag of the form 

Lastly, deploys to the live platform will present a confirmation prompt.

Note that the first time a stack is created there will be a significant delay before the subdomain is available due to 
propagation but subsequent updates should be quite fast.

### Exporting Swagger JSON and models

The API definition YAML includes a Swagger definition for the API.

This can be exported using the `export.sh` script.

In addition a schema-generator executable can be created from here: https://github.com/merlincox/generate

If this is added to the system path, the `export.sh` will also generate Go structs for the API and optionally replace 
the pkg/models/api.go file if that is out of sync with the API. (Therefore any additionmal models which do not feature 
directly in the API should be placed in the pkg/models/models.go file).

### Endpoints

The `/status` endpoint demonstrates that the environment has been passed to the lambda, and will return the git branch
commit and release tag, the platform and a timestamp for when the lambda was first invoked.

The `/calc` endpoint uses simple maths functions to demonstrate handling of path and query parameters, headers, error-handling and API-level caching.

Usage:

`https://{my-subdomain}[-{platform}].{my-domain}/calc/{op}?val1={val1}&val2={val2}`

where {op} can be one of "add", "substract", "multiply", "divide", "power" or "root" (all of which can be shortened to 
3 letters) and val1 and va12 are numbers.

The Accept-Language request header can optionally be used to format the result.

For example, `/calc/mul?val1=423.456&val2=30.1` with Accept-Language set to "en-GB" will return

```
{
     "locale": "en-GB",
     "op": "multiply",
     "result": "12,746.0256",
     "val1": 423.456,
     "val2": 30.1
}
```
whereas with Accept-Language as "fr-FR" it will produce
```
{
     "locale": "fr-FR",
     "op": "multiply",
     "result": "12 746,0256",
     "val1": 423.456,
     "val2": 30.1
}
```

API-level caching can determined by looking at the x-Timestamp response header. If you repeat a query and the value of 
this header does not change, you are seeing a cached response.


This endpoint also demonstrates error handling.

`/calc/div?val1=423.456&val2=0` will return

```
{
    "message": "Out of limits: 423.456 divide 0",
    "code": 400
}

```

`calc/bad?val1=423.456&val2=123` will return

```
{
    "message": "Unknown calc operation: bad",
    "code": 400
}
```
