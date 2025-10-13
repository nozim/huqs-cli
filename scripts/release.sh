#!/bin/bash

# Exit on error
set -e

GITHUB_TOKEN=${GITHUB_TOKEN}     # Replace with your GitHub token or set as an environment variable

GIT_SHORT_REV=${GIT_REV}

# Variables (Update as needed)
REPO_OWNER="nozim"           # Replace with your GitHub username or organization
REPO_NAME="huqs-cli"          # Replace with your repository name
TAG=$GIT_SHORT_REV                      # Replace with your release tag (e.g., v1.0.0)
BUILD_DIR="build"                    # Directory to store builds


# Ensure dependencies
command -v jq >/dev/null 2>&1 || { echo "jq is required. Install it and try again."; exit 1; }
command -v zip >/dev/null 2>&1 || { echo "zip is required. Install it and try again."; exit 1; }

for FILE in $BUILD_DIR/*; do
  BASENAME=$(basename $FILE)

  if [[ $FILE == *.exe ]]; then
    zip -j $BUILD_DIR/${BASENAME}.zip $FILE
  else
    tar -czvf $BUILD_DIR/${BASENAME}.tar.gz -C $BUILD_DIR $BASENAME
  fi
done


RELEASE="huqs-$TAG"

RELEASE_RESPONSE=$(curl -s \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/tags/$RELEASE")

RELEASE_ID=$(echo "$RELEASE_RESPONSE" | jq -r .id)

if [ "$RELEASE_ID" == "null" ]; then
  echo "Release with tag $TAG not found. Creating a new release..."
  RELEASE="huqs-$TAG"

  RELEASE_RESPONSE=$(curl -s \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github+json" \
    https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases \
    -d "{\"tag_name\": \"$RELEASE\", \"name\": \"$RELEASE\", \"body\": \"Release created via script.\"}")

  echo $RELEASE_RESPONSE

  RELEASE_ID=$(echo "$RELEASE_RESPONSE" | jq -r .id)
fi

if [ "$RELEASE_ID" == "null" ]; then
  echo "Failed to create or retrieve release. Exiting."
  echo "response $RELEASE_RESPONSE"
  exit 1
fi

echo "Release ID: $RELEASE_ID"
#
# Step 4: Upload assets
echo "Uploading assets to GitHub release"

UPLOAD_URL="https://uploads.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/$RELEASE_ID/assets"
echo "Uploading assets to GitHub release..."

for FILE in $BUILD_DIR/*.{tar.gz,zip}; do
  BASENAME=$(basename $FILE)
  echo "Uploading $BASENAME to $UPLOAD_URL..."
  if [ ! -e "$FILE" ]; then
      # Your processing logic here
      echo "Processing: $FILE"
      continue
  fi
  curl -s \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Content-Type: $(file --mime-type -b $FILE)" \
    --data-binary @$FILE \
    "$UPLOAD_URL?name=$BASENAME"

  echo "$BASENAME uploaded."
done
