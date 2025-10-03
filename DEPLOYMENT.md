# GeoCatalogo Deployment Guide

**Target**: catalog.dev.geosure.cloud
**Environment**: Development
**Authentication**: Ory Lamarr (auth.geosure.com)
**Infrastructure**: AWS Lambda + API Gateway + OpenSearch

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│  catalog.dev.geosure.cloud                                  │
│  (CloudFront + API Gateway)                                 │
└────────────────┬────────────────────────────────────────────┘
                 │
                 │ HTTPS
                 ↓
┌─────────────────────────────────────────────────────────────┐
│  API Gateway HTTP API                                       │
│  - Custom domain: catalog.dev.geosure.cloud                 │
│  - CORS enabled                                             │
│  - JWT Authorizer (Ory Lamarr)                              │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ├─→ Lambda Authorizer ──→ Ory Lamarr
                 │   (ExploreAuthorizer)    (auth.geosure.com)
                 │
                 ↓
┌─────────────────────────────────────────────────────────────┐
│  GeoCatalogo Lambda (Go ARM64)                              │
│  - Runtime: provided.al2                                    │
│  - Handler: bootstrap                                       │
│  - Memory: 512 MB                                           │
│  - Timeout: 30s                                             │
│  - VPC: Private subnets + NAT Gateway                       │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────────────────────────┐
│  Amazon OpenSearch                                          │
│  - Domain: gro-catalog-dev                                  │
│  - Instance: t3.small.search                                │
│  - Storage: 20 GB EBS                                       │
│  - Index: gro-catalog                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## Prerequisites

### 1. AWS Access

```bash
# Ensure you have AWS SSO configured for geosure-admin
aws sso login --profile geosure-admin

# Verify access
aws sts get-caller-identity --profile geosure-admin
```

### 2. Build Tools

```bash
# Install Go 1.21+
go version  # Should be >= 1.21

# Install AWS Lambda Go adapter
go get github.com/aws/aws-lambda-go/lambda
go get github.com/awslabs/aws-lambda-go-api-proxy
```

### 3. Ory Lamarr Configuration

You need:
- **Ory Issuer URL**: `https://auth.geosure.com` (or your Ory tenant URL)
- **JWT Audience**: Registered audience for the catalog API
- **Client credentials** for testing

---

## Step 1: Build GeoCatalogo Lambda Binary

### 1.1 Create Lambda-Compatible Main File

Create `cmd/lambda/main.go`:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/sirupsen/logrus"

	"github.com/geosure/geocatalogo"
	"github.com/geosure/geocatalogo/api"
)

var log = logrus.New()
var cat *geocatalogo.GeoCatalogue
var adapter *httpadapter.HandlerAdapter

func init() {
	// Load catalog from environment
	var err error
	cat, err = geocatalogo.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize catalog: %v", err)
	}

	// Create HTTP handler
	router := api.NewRouter(cat)
	adapter = httpadapter.New(router)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Log request
	log.WithFields(logrus.Fields{
		"method": req.RequestContext.HTTP.Method,
		"path":   req.RawPath,
	}).Info("Incoming request")

	// Proxy to HTTP handler
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
```

### 1.2 Build for Lambda (ARM64)

```bash
cd /Users/jjohnson/projects/geocatalogo

# Build for Lambda ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
  -tags lambda.norpc \
  -o bootstrap \
  cmd/lambda/main.go

# Verify binary
file bootstrap
# Should output: bootstrap: ELF 64-bit LSB executable, ARM aarch64

# Create deployment package
zip geocatalogo-lambda.zip bootstrap

# Check size
ls -lh geocatalogo-lambda.zip
# Should be around 15-20 MB
```

### 1.3 Upload to S3

```bash
# Get the lambda artifacts bucket from CloudFormation
BUCKET=$(aws cloudformation describe-stacks \
  --stack-name gro-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`LambdaArtifactsBucket`].OutputValue' \
  --output text \
  --profile geosure-admin)

echo "Lambda bucket: $BUCKET"

# Upload with versioning
VERSION=$(date +%Y%m%d-%H%M%S)
aws s3 cp geocatalogo-lambda.zip \
  s3://${BUCKET}/geocatalogo/${VERSION}/geocatalogo-lambda.zip \
  --profile geosure-admin

echo "Uploaded to: s3://${BUCKET}/geocatalogo/${VERSION}/geocatalogo-lambda.zip"
```

---

## Step 2: Create OpenSearch Domain

### 2.1 Add to app-infra.yml

Add this section to `/Users/jjohnson/projects/geosure/infra/app-infra.yml`:

```yaml
Parameters:
  # ... existing parameters ...

  CatalogLambdaS3Key:
    Type: String
    Description: S3 key for geocatalogo Lambda deployment package
    Default: geocatalogo/latest/geocatalogo-lambda.zip

Resources:
  # ... existing resources ...

  #############################################################################
  # GEOCATALOGO CATALOG SERVICE
  #############################################################################

  CatalogOpenSearchDomain:
    Type: AWS::OpenSearch::Domain
    Properties:
      DomainName: !Sub "${AppName}-catalog-${Environment}"
      EngineVersion: "OpenSearch_2.11"
      ClusterConfig:
        InstanceType: t3.small.search
        InstanceCount: 1
        DedicatedMasterEnabled: false
        ZoneAwarenessEnabled: false
      EBSOptions:
        EBSEnabled: true
        VolumeType: gp3
        VolumeSize: 20
      NodeToNodeEncryptionOptions:
        Enabled: true
      EncryptionAtRestOptions:
        Enabled: true
      DomainEndpointOptions:
        EnforceHTTPS: true
      AccessPolicies:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              AWS: !GetAtt CatalogLambdaRole.Arn
            Action:
              - "es:ESHttp*"
            Resource: !Sub "arn:aws:es:${AWS::Region}:${AWS::AccountId}:domain/${AppName}-catalog-${Environment}/*"
      VPCOptions:
        SubnetIds:
          - !Select [0, !Ref PrivateSubnetIds]
        SecurityGroupIds:
          - !Ref CatalogOpenSearchSecurityGroup

  CatalogOpenSearchSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Security group for Catalog OpenSearch
      VpcId: !Ref VpcId
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 443
          ToPort: 443
          SourceSecurityGroupId: !Ref EC2SecurityGroup
      Tags:
        - Key: Name
          Value: !Sub "${AppName}-catalog-opensearch-${Environment}"

  CatalogLambdaRole:
    Type: AWS::IAM::Role
    Properties:
      RoleName: !Sub "${AppName}-catalog-lambda-${Environment}"
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole
      Policies:
        - PolicyName: OpenSearchAccess
          PolicyDocument:
            Version: "2012-10-17"
            Statement:
              - Effect: Allow
                Action:
                  - "es:ESHttp*"
                Resource: !Sub "arn:aws:es:${AWS::Region}:${AWS::AccountId}:domain/${AppName}-catalog-${Environment}/*"
        - PolicyName: CloudWatchLogs
          PolicyDocument:
            Version: "2012-10-17"
            Statement:
              - Effect: Allow
                Action:
                  - "logs:CreateLogGroup"
                  - "logs:CreateLogStream"
                  - "logs:PutLogEvents"
                Resource: !Sub "arn:aws:logs:${AWS::Region}:${AWS::AccountId}:log-group:/aws/lambda/${AppName}-catalog-${Environment}:*"

  CatalogLambda:
    Type: AWS::Lambda::Function
    Properties:
      FunctionName: !Sub "${AppName}-catalog-${Environment}"
      Handler: bootstrap
      Runtime: provided.al2
      Architectures: [arm64]
      Role: !GetAtt CatalogLambdaRole.Arn
      Timeout: 30
      MemorySize: 512
      Code:
        S3Bucket: !Ref LambdaArtifactsBucket
        S3Key: !Ref CatalogLambdaS3Key
      Environment:
        Variables:
          GEOCATALOGO_REPOSITORY_TYPE: elasticsearch
          GEOCATALOGO_REPOSITORY_URL: !Sub "https://${CatalogOpenSearchDomain.DomainEndpoint}/gro-catalog/records"
          GEOCATALOGO_SERVER_URL: https://catalog.dev.geosure.cloud
          GEOCATALOGO_LOGGING_LEVEL: INFO
          GEOCATALOGO_SERVER_CORS: "true"
          GEOCATALOGO_SERVER_LIMIT: "100"
          GEOCATALOGO_METADATA_IDENTIFICATION_TITLE: "GRO Geospatial Data Catalog"
          GEOCATALOGO_METADATA_IDENTIFICATION_ABSTRACT: "Global Risk Observatory unified geospatial data catalog"
          GEOCATALOGO_METADATA_PROVIDER_NAME: "Geosure"
          GEOCATALOGO_METADATA_PROVIDER_URL: "https://geosure.ai"
      VpcConfig:
        SubnetIds: !Ref PrivateSubnetIds
        SecurityGroupIds:
          - !Ref EC2SecurityGroup

  CatalogApi:
    Type: AWS::ApiGatewayV2::Api
    Properties:
      Name: !Sub "${AppName}-catalog-${Environment}"
      ProtocolType: HTTP
      CorsConfiguration:
        AllowOrigins:
          - https://timeline.dev.geosure.cloud
          - https://gro.dev.geosure.cloud
        AllowMethods:
          - GET
          - POST
          - OPTIONS
        AllowHeaders:
          - Content-Type
          - Authorization
          - Cookie
        MaxAge: 3600

  CatalogApiIntegration:
    Type: AWS::ApiGatewayV2::Integration
    Properties:
      ApiId: !Ref CatalogApi
      IntegrationType: AWS_PROXY
      IntegrationUri: !Sub "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${CatalogLambda.Arn}/invocations"
      PayloadFormatVersion: "2.0"

  CatalogLambdaAuthorizer:
    Type: AWS::ApiGatewayV2::Authorizer
    Properties:
      ApiId: !Ref CatalogApi
      AuthorizerType: REQUEST
      AuthorizerUri: !Sub "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${ExploreAuthorizer.Arn}/invocations"
      AuthorizerPayloadFormatVersion: "2.0"
      EnableSimpleResponses: true
      IdentitySource:
        - "$request.header.Cookie"
      Name: CatalogLambdaAuth

  # Public route for OPTIONS (CORS preflight)
  CatalogOptionsRoute:
    Type: AWS::ApiGatewayV2::Route
    Properties:
      ApiId: !Ref CatalogApi
      RouteKey: "OPTIONS /{proxy+}"
      Target: !Sub integrations/${CatalogApiIntegration}

  # Authenticated routes
  CatalogRootRoute:
    Type: AWS::ApiGatewayV2::Route
    Properties:
      ApiId: !Ref CatalogApi
      RouteKey: "GET /"
      Target: !Sub integrations/${CatalogApiIntegration}
      AuthorizationType: CUSTOM
      AuthorizerId: !Ref CatalogLambdaAuthorizer

  CatalogProxyRoute:
    Type: AWS::ApiGatewayV2::Route
    Properties:
      ApiId: !Ref CatalogApi
      RouteKey: "GET /{proxy+}"
      Target: !Sub integrations/${CatalogApiIntegration}
      AuthorizationType: CUSTOM
      AuthorizerId: !Ref CatalogLambdaAuthorizer

  CatalogStage:
    Type: AWS::ApiGatewayV2::Stage
    Properties:
      ApiId: !Ref CatalogApi
      StageName: "$default"
      AutoDeploy: true
      AccessLogSettings:
        DestinationArn: !GetAtt CatalogApiLogGroup.Arn
        Format: '$context.requestId $context.error.message $context.error.messageString'

  CatalogApiLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: !Sub "/aws/apigateway/${AppName}-catalog-${Environment}"
      RetentionInDays: 7

  CatalogLambdaInvokePerm:
    Type: AWS::Lambda::Permission
    Properties:
      Action: lambda:InvokeFunction
      FunctionName: !Ref CatalogLambda
      Principal: apigateway.amazonaws.com
      SourceArn: !Sub "arn:aws:execute-api:${AWS::Region}:${AWS::AccountId}:${CatalogApi}/*/*"

  CatalogAuthorizerInvokePerm:
    Type: AWS::Lambda::Permission
    Properties:
      Action: lambda:InvokeFunction
      FunctionName: !Ref ExploreAuthorizer
      Principal: apigateway.amazonaws.com
      SourceArn: !Sub "arn:aws:execute-api:${AWS::Region}:${AWS::AccountId}:${CatalogApi}/authorizers/*"

  CatalogDomainMapping:
    Type: AWS::ApiGatewayV2::DomainName
    Properties:
      DomainName: catalog.dev.geosure.cloud
      DomainNameConfigurations:
        - CertificateArn: !Ref ACMCertificateArn
          EndpointType: REGIONAL

  CatalogApiMapping:
    Type: AWS::ApiGatewayV2::ApiMapping
    Properties:
      ApiId: !Ref CatalogApi
      DomainName: !Ref CatalogDomainMapping
      Stage: !Ref CatalogStage

Outputs:
  # ... existing outputs ...

  CatalogApiEndpoint:
    Description: Catalog API endpoint
    Value: !GetAtt CatalogApi.ApiEndpoint
    Export:
      Name: !Sub "${AppName}-${Environment}-catalog-api-endpoint"

  CatalogDomainName:
    Description: Catalog custom domain name
    Value: !Ref CatalogDomainMapping
    Export:
      Name: !Sub "${AppName}-${Environment}-catalog-domain"

  CatalogOpenSearchEndpoint:
    Description: Catalog OpenSearch domain endpoint
    Value: !GetAtt CatalogOpenSearchDomain.DomainEndpoint
    Export:
      Name: !Sub "${AppName}-${Environment}-catalog-opensearch-endpoint"
```

---

## Step 3: Deploy CloudFormation Stack

### 3.1 Validate Template

```bash
cd /Users/jjohnson/projects/geosure/infra

aws cloudformation validate-template \
  --template-body file://app-infra.yml \
  --profile geosure-admin
```

### 3.2 Deploy Stack Update

```bash
# Set parameters
STACK_NAME="gro-dev"
S3_KEY="geocatalogo/20251003-132000/geocatalogo-lambda.zip"  # Use your uploaded version

# Update stack
aws cloudformation update-stack \
  --stack-name ${STACK_NAME} \
  --template-body file://app-infra.yml \
  --parameters \
    ParameterKey=CatalogLambdaS3Key,ParameterValue=${S3_KEY} \
    ParameterKey=OryIssuer,UsePreviousValue=true \
    ParameterKey=JwtAudiences,UsePreviousValue=true \
  --capabilities CAPABILITY_NAMED_IAM \
  --profile geosure-admin

# Monitor deployment
aws cloudformation wait stack-update-complete \
  --stack-name ${STACK_NAME} \
  --profile geosure-admin

echo "Stack updated successfully!"
```

### 3.3 Get Outputs

```bash
# Get API endpoint
aws cloudformation describe-stacks \
  --stack-name ${STACK_NAME} \
  --query 'Stacks[0].Outputs[?OutputKey==`CatalogApiEndpoint`].OutputValue' \
  --output text \
  --profile geosure-admin

# Get OpenSearch endpoint
aws cloudformation describe-stacks \
  --stack-name ${STACK_NAME} \
  --query 'Stacks[0].Outputs[?OutputKey==`CatalogOpenSearchEndpoint`].OutputValue' \
  --output text \
  --profile geosure-admin
```

---

## Step 4: Configure DNS (Route 53)

### 4.1 Get API Gateway Domain Name

```bash
DOMAIN_TARGET=$(aws apigatewayv2 get-domain-name \
  --domain-name catalog.dev.geosure.cloud \
  --query 'DomainNameConfigurations[0].ApiGatewayDomainName' \
  --output text \
  --profile geosure-admin)

echo "Target: $DOMAIN_TARGET"
```

### 4.2 Create DNS Record

```bash
# Get hosted zone ID for dev.geosure.cloud
HOSTED_ZONE_ID=$(aws route53 list-hosted-zones \
  --query 'HostedZones[?Name==`dev.geosure.cloud.`].Id' \
  --output text \
  --profile geosure-admin | cut -d/ -f3)

echo "Hosted Zone: $HOSTED_ZONE_ID"

# Create change batch
cat > /tmp/catalog-dns-change.json <<EOF
{
  "Changes": [{
    "Action": "UPSERT",
    "ResourceRecordSet": {
      "Name": "catalog.dev.geosure.cloud",
      "Type": "CNAME",
      "TTL": 300,
      "ResourceRecords": [{
        "Value": "${DOMAIN_TARGET}"
      }]
    }
  }]
}
EOF

# Apply change
aws route53 change-resource-record-sets \
  --hosted-zone-id ${HOSTED_ZONE_ID} \
  --change-batch file:///tmp/catalog-dns-change.json \
  --profile geosure-admin
```

---

## Step 5: Initialize OpenSearch Index

### 5.1 Create Index with Mapping

```bash
# Get OpenSearch endpoint
OS_ENDPOINT=$(aws cloudformation describe-stacks \
  --stack-name gro-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`CatalogOpenSearchEndpoint`].OutputValue' \
  --output text \
  --profile geosure-admin)

echo "OpenSearch: https://${OS_ENDPOINT}"

# Create index with geospatial mapping
curl -X PUT "https://${OS_ENDPOINT}/gro-catalog" \
  -H 'Content-Type: application/json' \
  -d '{
    "mappings": {
      "properties": {
        "geometry": {
          "type": "geo_shape"
        },
        "properties": {
          "properties": {
            "title": {
              "type": "text",
              "fields": {
                "keyword": { "type": "keyword" }
              }
            },
            "abstract": {
              "type": "text"
            },
            "collection": {
              "type": "keyword"
            },
            "datetime": {
              "type": "date"
            }
          }
        }
      }
    }
  }'
```

### 5.2 Bulk Load Catalog Records

```bash
# Convert catalog to bulk format
cd /Users/jjohnson/projects/geosure/catalog

python3 <<'PYEOF'
import json

with open('data/geocatalogo_records.json') as f:
    records = json.load(f)

with open('/tmp/bulk-catalog.ndjson', 'w') as out:
    for record in records:
        # Index action
        action = {"index": {"_index": "gro-catalog", "_id": record["id"]}}
        out.write(json.dumps(action) + '\n')
        # Document
        out.write(json.dumps(record) + '\n')

print(f"Created bulk file with {len(records)} records")
PYEOF

# Upload to OpenSearch
curl -X POST "https://${OS_ENDPOINT}/_bulk" \
  -H 'Content-Type: application/x-ndjson' \
  --data-binary @/tmp/bulk-catalog.ndjson

# Verify count
curl "https://${OS_ENDPOINT}/gro-catalog/_count"
```

---

## Step 6: Test Deployment

### 6.1 Test Without Authentication (Should Fail)

```bash
curl https://catalog.dev.geosure.cloud/?q=wildfire
# Should return 401 Unauthorized
```

### 6.2 Authenticate with Ory

```bash
# Get access token from Ory
# (Replace with your Ory client credentials)

ORY_CLIENT_ID="your-client-id"
ORY_CLIENT_SECRET="your-client-secret"
ORY_ISSUER="https://auth.geosure.com"

ACCESS_TOKEN=$(curl -X POST "${ORY_ISSUER}/oauth2/token" \
  -d "grant_type=client_credentials" \
  -d "client_id=${ORY_CLIENT_ID}" \
  -d "client_secret=${ORY_CLIENT_SECRET}" \
  -d "scope=catalog:read" \
  | jq -r '.access_token')

echo "Access token: ${ACCESS_TOKEN:0:20}..."
```

### 6.3 Test Authenticated Request

```bash
# Test with bearer token
curl -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  https://catalog.dev.geosure.cloud/?q=wildfire

# Test with cookie (if using session)
curl --cookie "id_token=${ACCESS_TOKEN}" \
  https://catalog.dev.geosure.cloud/?q=wildfire
```

### 6.4 Test Various Queries

```bash
# Search by keyword
curl --cookie "id_token=${ACCESS_TOKEN}" \
  'https://catalog.dev.geosure.cloud/?q=california'

# Get by ID
curl --cookie "id_token=${ACCESS_TOKEN}" \
  'https://catalog.dev.geosure.cloud/?recordids=db_h3_l8_union_cities_urban'

# Pagination
curl --cookie "id_token=${ACCESS_TOKEN}" \
  'https://catalog.dev.geosure.cloud/?q=database&from=0&size=5'
```

---

## Step 7: CI/CD Setup (GitHub Actions)

Create `.github/workflows/deploy-catalog.yml` in the geocatalogo repository:

```yaml
name: Deploy GeoCatalogo to AWS

on:
  push:
    branches: [gro]
  workflow_dispatch:

env:
  AWS_REGION: us-east-1
  GO_VERSION: "1.21"

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::${{ secrets.AWS_ACCOUNT_ID }}:role/GitHubActionsRole
          aws-region: ${{ env.AWS_REGION }}

      - name: Build Lambda binary
        run: |
          GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
            -tags lambda.norpc \
            -o bootstrap \
            cmd/lambda/main.go

          zip geocatalogo-lambda.zip bootstrap

          ls -lh geocatalogo-lambda.zip

      - name: Upload to S3
        run: |
          VERSION=$(date +%Y%m%d-%H%M%S)-${GITHUB_SHA:0:7}

          aws s3 cp geocatalogo-lambda.zip \
            s3://${{ secrets.LAMBDA_BUCKET }}/geocatalogo/${VERSION}/geocatalogo-lambda.zip

          # Update "latest" pointer
          aws s3 cp geocatalogo-lambda.zip \
            s3://${{ secrets.LAMBDA_BUCKET }}/geocatalogo/latest/geocatalogo-lambda.zip

          echo "Deployed version: ${VERSION}"

      - name: Update Lambda function
        run: |
          aws lambda update-function-code \
            --function-name gro-catalog-dev \
            --s3-bucket ${{ secrets.LAMBDA_BUCKET }} \
            --s3-key geocatalogo/latest/geocatalogo-lambda.zip

          # Wait for update to complete
          aws lambda wait function-updated \
            --function-name gro-catalog-dev

          echo "Lambda function updated successfully"

      - name: Test deployment
        run: |
          # Get function info
          aws lambda get-function \
            --function-name gro-catalog-dev \
            --query 'Configuration.[FunctionName,LastModified,State]' \
            --output table
```

Add GitHub secrets:
- `AWS_ACCOUNT_ID`
- `LAMBDA_BUCKET`

---

## Monitoring and Operations

### View Lambda Logs

```bash
# Stream logs
aws logs tail /aws/lambda/gro-catalog-dev \
  --follow \
  --profile geosure-admin

# View recent errors
aws logs filter-log-events \
  --log-group-name /aws/lambda/gro-catalog-dev \
  --filter-pattern "ERROR" \
  --start-time $(date -u -d '1 hour ago' +%s)000 \
  --profile geosure-admin
```

### View API Gateway Logs

```bash
aws logs tail /aws/apigateway/gro-catalog-dev \
  --follow \
  --profile geosure-admin
```

### Monitor OpenSearch

```bash
# Cluster health
curl "https://${OS_ENDPOINT}/_cluster/health?pretty"

# Index stats
curl "https://${OS_ENDPOINT}/gro-catalog/_stats?pretty"

# Search performance
curl "https://${OS_ENDPOINT}/gro-catalog/_search" \
  -H 'Content-Type: application/json' \
  -d '{
    "profile": true,
    "query": {
      "match": {
        "properties.title": "wildfire"
      }
    }
  }'
```

### Update Catalog Data

```bash
# Re-run converter in catalog repo
cd /Users/jjohnson/projects/geosure/catalog
python3 scripts/convert_to_geocatalogo.py

# Upload to OpenSearch (see Step 5.2)
```

---

## Troubleshooting

### Issue: 401 Unauthorized

**Cause**: JWT token invalid or expired

**Solution**:
```bash
# Check Ory issuer configuration
aws lambda get-function-configuration \
  --function-name gro-catalog-dev-authorizer \
  --query 'Environment.Variables.ORY_ISSUER' \
  --profile geosure-admin

# Verify it matches your Ory tenant URL
```

### Issue: 500 Internal Server Error

**Cause**: Lambda error or OpenSearch connection issue

**Solution**:
```bash
# Check Lambda logs
aws logs tail /aws/lambda/gro-catalog-dev --follow --profile geosure-admin

# Check OpenSearch connectivity from Lambda VPC
aws lambda invoke \
  --function-name gro-catalog-dev \
  --payload '{"rawPath": "/", "requestContext": {"http": {"method": "GET"}}}' \
  /tmp/response.json \
  --profile geosure-admin

cat /tmp/response.json | jq
```

### Issue: No search results

**Cause**: Index not populated or mapping incorrect

**Solution**:
```bash
# Check index exists
curl "https://${OS_ENDPOINT}/_cat/indices?v"

# Check document count
curl "https://${OS_ENDPOINT}/gro-catalog/_count"

# Re-index if needed (see Step 5.2)
```

---

## Security Considerations

1. **Ory Lamarr Integration**:
   - All routes (except OPTIONS) require valid JWT
   - Token verified by Lambda Authorizer (ExploreAuthorizer)
   - Supports both cookie-based and bearer token auth

2. **VPC Configuration**:
   - Lambda runs in private subnets
   - OpenSearch in VPC (not internet-accessible)
   - Security groups restrict access

3. **Encryption**:
   - HTTPS enforced (API Gateway + CloudFront)
   - OpenSearch encryption at rest enabled
   - OpenSearch node-to-node encryption enabled

4. **IAM Roles**:
   - Lambda has minimal permissions (OpenSearch + CloudWatch only)
   - OpenSearch access policy restricts to Lambda role only

---

## Cost Estimate

**Development Environment (dev)**:

| Service | Configuration | Monthly Cost |
|---------|--------------|--------------|
| Lambda | 512 MB, ~1M requests | $1.50 |
| API Gateway | HTTP API, 1M requests | $1.00 |
| OpenSearch | t3.small.search, 20 GB | $45.00 |
| Data Transfer | ~10 GB outbound | $0.90 |
| CloudWatch Logs | 7-day retention, 1 GB | $0.50 |
| **Total** | | **~$49/month** |

**Production Environment (prod)**:
- OpenSearch: 3× t3.medium.search = $195/month
- Total: ~$250-300/month (depending on traffic)

---

## Maintenance

### Monthly Tasks

1. Review CloudWatch metrics and logs
2. Check OpenSearch cluster health
3. Update catalog data if schema changes
4. Review and rotate API credentials
5. Update dependencies (go mod tidy)

### Quarterly Tasks

1. Review and update Ory Lamarr configuration
2. Performance testing and optimization
3. Disaster recovery testing
4. Cost optimization review

---

## References

- **GeoCatalogo Docs**: https://github.com/geosure/geocatalogo
- **GRO Integration**: `/Users/jjohnson/projects/geocatalogo/GRO_INTEGRATION.md`
- **Catalog Repository**: `/Users/jjohnson/projects/geosure/catalog/`
- **Ory Documentation**: https://www.ory.sh/docs
- **AWS Lambda Go**: https://github.com/aws/aws-lambda-go
- **OpenSearch**: https://opensearch.org/docs/latest/

---

**Last Updated**: October 3, 2025
**Maintainer**: @data-platform
**Status**: Ready for deployment
