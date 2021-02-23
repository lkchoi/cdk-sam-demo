import * as cdk from '@aws-cdk/core';
import * as apigw from '@aws-cdk/aws-apigateway'
import * as lambda from '@aws-cdk/aws-lambda'
import * as dynamodb from '@aws-cdk/aws-dynamodb'
import { resolve } from 'path'

const PROJECT_ROOT = resolve(__dirname, '../../')

export class GeminiCdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Lambda Function
    const handler = new lambda.Function(this, 'GeminiFunction', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'main',
      code: lambda.Code.fromAsset(resolve(PROJECT_ROOT, 'dist')),
      timeout: cdk.Duration.seconds(5),
    });

    // DynamoDB Table
    const table = new dynamodb.Table(this, 'GeminiTable', {
      tableName: 'Gemini',
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

    // Permissions
    table.grantReadWriteData(handler)

    // API Gateway REST API
    new apigw.LambdaRestApi(this, 'GeminiApi', { handler })
  }
}
