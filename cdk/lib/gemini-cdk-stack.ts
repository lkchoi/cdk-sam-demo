import * as cdk from '@aws-cdk/core';
import * as apigw from '@aws-cdk/aws-apigateway'
import * as lambda from '@aws-cdk/aws-lambda'
import * as dynamodb from '@aws-cdk/aws-dynamodb'
import { resolve } from 'path'

const PROJECT_ROOT = resolve(__dirname, '../../')

export class GeminiCdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const handler = new lambda.Function(this, 'GeminiFunction', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'main',
      code: lambda.Code.fromAsset(resolve(PROJECT_ROOT, 'dist')),
      timeout: cdk.Duration.seconds(30),
    });

    const table = new dynamodb.Table(this, 'GeminiTable', {
      partitionKey: {
        name: 'PK',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'SK',
        type: dynamodb.AttributeType.STRING,
      },
      // encryption: dynamodb.TableEncryption.AWS_MANAGED
    });
    table.grantReadWriteData(handler)
    new apigw.LambdaRestApi(this, 'GeminiApi', { handler })
  }
}
