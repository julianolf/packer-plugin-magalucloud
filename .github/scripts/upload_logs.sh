#!/bin/bash

set -e

LOG_FILES=(
	"builder/magalucloud/packer_log_magalucloud_builder_test.txt"
	"post-processor/magalucloud/packer_log_magalucloud_importer_test.txt"
)

for FILE in "${LOG_FILES[@]}"; do
	if [ -f "$FILE" ]; then
		echo "Found log file: $FILE. Uploading to Object Storage..."
		FILENAME=$(basename "$FILE")
		aws s3 cp "$FILE" "s3://$BUCKET/logs/${RUN_ID}-$FILENAME" --endpoint-url "$S3_ENDPOINT" --quiet
	fi
done
