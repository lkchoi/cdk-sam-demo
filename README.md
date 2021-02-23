# CDK SAM Demo

This project uses CDK and SAM to

1. Install and run dynamodb locally
https://hub.docker.com/r/amazon/dynamodb-local/

2. install nvm (node version manager)
https://github.com/nvm-sh/nvm

```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.37.2/install.sh | bash
nvm install --lts
```

3. install and run dynamodb admin (GUI)
https://github.com/aaronshaf/dynamodb-admin

4. clone (TODO push gemini)

5. CDK won't deploy to DynamoDB local, must use dynamodb-admin GUI or "create-table" script in Makefile

6. Set DynamoDB endpoint URL based on env
local: "http://docker.for.mac.localhost:8000"
prod: https://docs.aws.amazon.com/general/latest/gr/rande.html#ddb_regionhttps://docs.aws.amazon.com/general/latest/gr/rande.html#ddb_region
