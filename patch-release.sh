#!/bin/bash

echo "Build check..."
go build || exit 1

LATEST_TAG=$(git tag -l --sort=-version:refname | head -n1)

if [ -z "$LATEST_TAG" ]; then
    # If no previous tag found
    NEW_VERSION="v0.0.1"
else
    VERSION_NUM=$(echo $LATEST_TAG | sed 's/v//')
    IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUM"

    # increase only patch version
    NEW_PATCH=$((PATCH + 1))
    NEW_VERSION="v${MAJOR}.${MINOR}.${NEW_PATCH}"
fi

echo "Current version: ${LATEST_TAG:-none}"
echo "New     version: $NEW_VERSION"

# Tag create and push
git tag $NEW_VERSION
git push origin $NEW_VERSION

echo "Done: $NEW_VERSION"
