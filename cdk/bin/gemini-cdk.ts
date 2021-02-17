#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { GeminiCdkStack } from '../lib/gemini-cdk-stack';

const app = new cdk.App();
new GeminiCdkStack(app, 'GeminiCdkStack');
