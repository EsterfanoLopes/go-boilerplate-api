set -x
awslocal s3 mb s3://advertiser-config
echo '{}' > dev-configs.json
awslocal s3 cp dev-configs.json s3://advertiser-config/
awslocal s3 mb s3://datalake-bucket
awslocal s3 mb s3://email-attachments-bucket
set +x