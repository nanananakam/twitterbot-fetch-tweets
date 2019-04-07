#!/usr/bin/env sh
printenv
aws s3 cp s3://${AWS_S3_BUCKET}/tweets.tar.xz .
tar Jxvf tweets.tar.xz
rm tweets.tar.xz
/main
tar Jcvf tweets.tar.xz tweets.db
aws s3 cp tweets.tar.xz s3://${AWS_S3_BUCKET}/tweets.tar.xz