#!/bin/bash

FEATURE=$1

if [ -z "$FEATURE" ]; then
  echo "Usage: ./generate-feature.sh <feature-name>"
  exit 1
fi

# LOWERCASE FEATURE
FEATURE_LOWER=$(echo "$FEATURE" | tr '[:upper:]' '[:lower:]')
BASE_DIR="internal/domains/$FEATURE_LOWER"

mkdir -p "$BASE_DIR/controller"
mkdir -p "$BASE_DIR/entity"
mkdir -p "$BASE_DIR/params"
mkdir -p "$BASE_DIR/repository"
mkdir -p "$BASE_DIR/service"

# Fungsi untuk membuat file dan menuliskan package
create_file_with_package() {
  local filepath=$1
  local package=$2
  echo "package $package" > "$filepath"
}

# ENTITY
create_file_with_package "$BASE_DIR/entity/${FEATURE_LOWER}_entity.go" "entity"

# PARAMS
create_file_with_package "$BASE_DIR/params/${FEATURE_LOWER}_request.go" "params"
create_file_with_package "$BASE_DIR/params/${FEATURE_LOWER}_response.go" "params"

# CONTROLLER
create_file_with_package "$BASE_DIR/controller/${FEATURE_LOWER}_controller.go" "controller"
create_file_with_package "$BASE_DIR/controller/${FEATURE_LOWER}_controller_impl.go" "controller"

# REPOSITORY
create_file_with_package "$BASE_DIR/repository/${FEATURE_LOWER}_repository.go" "repository"
create_file_with_package "$BASE_DIR/repository/${FEATURE_LOWER}_repository_impl.go" "repository"

# SERVICE
create_file_with_package "$BASE_DIR/service/${FEATURE_LOWER}_service.go" "service"
create_file_with_package "$BASE_DIR/service/${FEATURE_LOWER}_service_impl.go" "service"

# ROUTE
echo "package $FEATURE_LOWER" > "$BASE_DIR/route.go"

echo "✅ Generated feature '$FEATURE' in $BASE_DIR with packages"
