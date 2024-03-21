## Overview
Amazon Bedrock Metadata is an open-source project that processes Amazon Bedrock model invocation logs to generate insights focusing on usage and cost. This tool consumes model invocation logs stored in an Amazon S3 bucket, redacts input prompts and completion responses, generates metadata related to usage, cost and carbon footprint, and stores the enhanced metadata in another S3 bucket. Users can leverage Amazon Athena and Quicksight for visualizing and analyzing their Amazon Bedrock model usage, facilitating better decision-making regarding cost allocation, optimizing usage and estimating environmental impact.

## Features
- **Log Consumption**: Automatically fetches Amazon Bedrock model invocation logs in near real-time from a designated S3 bucket.
- **Redact Prompt and Response**: Strips away user prompts and completion responses to minimize the exposure risk of personally identifiable information (PII).
- **Metadata Generation**: Enriches logs with metadata including per model invocation cost calculation, carbon footprint estimation, and usage tags aggregation.
- **Data Storage**: Stores the processed metadata in a specified S3 bucket, ready for analysis using Amazon Athena and Quicksight.
- **Integration**: Seamlessly integrates with [Greenscale AI](https://www.greenscale.ai) platform for richer and customizable insights and recommendations on usage and cost optimization.

## Prerequisites
Amazon S3 buckets with [Amazon Bedrock model invocation logs](https://docs.aws.amazon.com/bedrock/latest/userguide/model-invocation-logging.html#setup-s3-destination) needs to be configured.

## Quick Deploy
[![Launch Stack](https://cdn.rawgit.com/buildkite/cloudformation-launch-stack-button-svg/master/launch-stack.svg)](https://us-east-1.console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=AmazonBedrockMetadata&templateURL=https://greenscale-ai-public.s3.amazonaws.com/amazon-bedrock-metadata/template.json)

## Usage
Once the setup is complete, the Amazon Bedrock Metadata will automatically process new invocation logs as they are generated. You can then use Amazon Athena to query the metadata and Amazon Quicksight for in-depth analysis and visualization.

## License
This project is open-source and available under the MIT License.

## Support and Contact
For support, feature requests, or bug reporting, please open an issue in the GitHub repository. For direct inquiries, contact the project maintainers via email (hello@greenscale.ai).
