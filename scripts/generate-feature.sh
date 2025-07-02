#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATES_DIR="$SCRIPT_DIR/templates"

FEATURE=$1
FEATURE=$(echo "$FEATURE" | xargs)


if [ -z "$FEATURE" ]; then
  echo "Usage: ./generate-feature.sh <feature-name>"
  echo "Example: ./generate-feature.sh Account"
  exit 1
fi

FEATURE_SANITIZED=$(echo "$FEATURE" | tr -cs '[:alnum:]' '_' | tr '[:upper:]' '[:lower:]')
FEATURE_SANITIZED=$(echo "$FEATURE_SANITIZED" | sed 's/_$//')

FEATURE_CAPITALIZED=$(echo "${FEATURE_SANITIZED^}")
BASE_DIR="internal/domains/$FEATURE_SANITIZED"

mkdir -p "$BASE_DIR/controller" "$BASE_DIR/entity" "$BASE_DIR/params" "$BASE_DIR/repository" "$BASE_DIR/service"

create_file_from_template() {
  local filepath=$1
  local package=$2
  local template=$3

  if [ -e "$filepath" ]; then
    echo "⚠️  File $filepath already exists, skipping..."
  else
    {
      echo "package $package"
      echo
      sed "s/{{Feature}}/$FEATURE_CAPITALIZED/g" "$template"
    } > "$filepath"
  fi
}

create_file_from_template "$BASE_DIR/entity/${FEATURE_SANITIZED}_entity.go" "entity" "$TEMPLATES_DIR/entity.tpl"
create_file_from_template "$BASE_DIR/params/${FEATURE_SANITIZED}_request.go" "params" "$TEMPLATES_DIR/params_request.tpl"
create_file_from_template "$BASE_DIR/params/${FEATURE_SANITIZED}_response.go" "params" "$TEMPLATES_DIR/params_response.tpl"
create_file_from_template "$BASE_DIR/controller/${FEATURE_SANITIZED}_controller.go" "controller" "$TEMPLATES_DIR/controller.tpl"
create_file_from_template "$BASE_DIR/controller/${FEATURE_SANITIZED}_controller_impl.go" "controller" "$TEMPLATES_DIR/controller_impl.tpl"
create_file_from_template "$BASE_DIR/repository/${FEATURE_SANITIZED}_repository.go" "repository" "$TEMPLATES_DIR/repository.tpl"
create_file_from_template "$BASE_DIR/repository/${FEATURE_SANITIZED}_repository_impl.go" "repository" "$TEMPLATES_DIR/repository_impl.tpl"
create_file_from_template "$BASE_DIR/service/${FEATURE_SANITIZED}_service.go" "service" "$TEMPLATES_DIR/service.tpl"
create_file_from_template "$BASE_DIR/service/${FEATURE_SANITIZED}_service_impl.go" "service" "$TEMPLATES_DIR/service_impl.tpl"

# ROUTE
{
  echo "package $FEATURE_SANITIZED"
  echo
  sed "s/{{Feature}}/$FEATURE_CAPITALIZED/g" "$TEMPLATES_DIR/route.tpl"
} > "$BASE_DIR/route.go"

GREEN='\033[0;32m'
NC='\033[0m'
echo -e "${GREEN}✅ Generated feature '$FEATURE' in $BASE_DIR with templates${NC}"
