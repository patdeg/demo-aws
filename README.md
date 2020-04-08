# AWS Demo with S3 and Athena

Showcasing automated processing when loading data to S3 into a parquet file to be analyzed in Athena

# Prerequisites

Install AWS CLI following https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html, and Elastic Beanstalk CLI with:

```
pip install awsebcli --upgrade --user
```

# Configure AWS & Elastic Beanstalk 

Set up AWS CLI with:

```
aws configure
```

Set up Elastic Beanstalk in the cloned folder:

```
cd demo-aws
eb init
eb create
```

# Deploy the application

You deploy the application to Elastic Beanstalk with

```
eb deploy
```

# Set up Notification

... TO BE FURTHER REFINE ...

In AWS console for Simple Notification System, create a new topic, and subscribe your Elastic Beanstalk application with the url, for example http://gotest-env.eba-12345.us-west-1.elasticbeanstalk.com/event. You can subscribe your email as well to debug notification.

Now, copy the ARN for the Topic, and enter it in the S3 Events settings under the Properties menu of your bucket. 

# Test the processing

You can now process data. 

Check the content of your folders:

```
aws s3 ls s3://deglon/data/
aws s3 ls s3://deglon/processed/
```

Copy the file ```test.json``` to s3 with

```
aws s3 cp test.json s3://deglon/data/
```

Within a second, you should have the parquet file. Check it with:
```
aws s3 ls s3://deglon/processed/
```

# Analyze the data in Athena

Create a database in AWS Glue.

In Athena, select the new database and create a table from the s3 /processed/ folder.

*Et voila!*


