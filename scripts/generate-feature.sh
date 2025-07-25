#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATES_DIR="$SCRIPT_DIR/templates"

# Load .env
ENV_FILE=".env"
if [ -f "$ENV_FILE" ]; then
  export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

# Validasi GO_MODULE
if [ -z "$GO_MODULE" ]; then
  echo "❌ GO_MODULE is not defined in .env"
  exit 1
fi

FEATURE=$1
FLAG=$2

FEATURE=$(echo "$FEATURE" | xargs)

if [ -z "$FEATURE" ]; then
  echo "Usage: ./generate-feature.sh <feature-name> [--resource]"
  echo "Example: ./generate-feature.sh Account --resource"
  exit 1
fi

IS_RESOURCE=false
if [ "$FLAG" == "--resource" ]; then
  IS_RESOURCE=true
  echo "🔧 Using resource-based templates..."
fi

FEATURE_SANITIZED=$(echo "$FEATURE" | tr -cs '[:alnum:]' '_' | tr '[:upper:]' '[:lower:]')
FEATURE_SANITIZED=$(echo "$FEATURE_SANITIZED" | sed 's/_$//')
FEATURE_CAPITALIZED=$(echo "${FEATURE_SANITIZED^}")

BASE_DIR="internal/domains/$FEATURE_SANITIZED"
mkdir -p "$BASE_DIR/controller" "$BASE_DIR/entity" "$BASE_DIR/params" "$BASE_DIR/repository" "$BASE_DIR/service"

# Template selector based on resource flag
select_template() {
  local default="$1"
  if [ "$IS_RESOURCE" = true ]; then
    echo "$TEMPLATES_DIR/resources/resource_$default"
  else
    echo "$TEMPLATES_DIR/$default"
  fi
}

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

      sed \
        -e "s|{{Feature}}|$FEATURE_CAPITALIZED|g" \
        -e "s|{{feature}}|$FEATURE_SANITIZED|g" \
        -e "s|{{Module}}|$GO_MODULE|g" \
        "$template"
    } > "$filepath"
  fi
}

# Generate from templates
create_file_from_template "$BASE_DIR/entity/${FEATURE_SANITIZED}_entity.go" "entity" "$TEMPLATES_DIR/entity.tpl"
create_file_from_template "$BASE_DIR/params/${FEATURE_SANITIZED}_request.go" "params" "$TEMPLATES_DIR/params_request.tpl"
create_file_from_template "$BASE_DIR/params/${FEATURE_SANITIZED}_response.go" "params" "$TEMPLATES_DIR/params_response.tpl"

create_file_from_template "$BASE_DIR/controller/${FEATURE_SANITIZED}_controller.go" "controller" "$(select_template controller.tpl)"
create_file_from_template "$BASE_DIR/controller/${FEATURE_SANITIZED}_controller_impl.go" "controller" "$(select_template controller_impl.tpl)"
create_file_from_template "$BASE_DIR/repository/${FEATURE_SANITIZED}_repository.go" "repository" "$(select_template repository.tpl)"
create_file_from_template "$BASE_DIR/repository/${FEATURE_SANITIZED}_repository_impl.go" "repository" "$(select_template repository_impl.tpl)"
create_file_from_template "$BASE_DIR/service/${FEATURE_SANITIZED}_service.go" "service" "$(select_template service.tpl)"
create_file_from_template "$BASE_DIR/service/${FEATURE_SANITIZED}_service_impl.go" "service" "$(select_template service_impl.tpl)"

# ROUTE
sed \
  -e "s|{{Feature}}|$FEATURE_CAPITALIZED|g" \
  -e "s|{{feature}}|$FEATURE_SANITIZED|g" \
  -e "s|{{Module}}|$GO_MODULE|g" \
  "$TEMPLATES_DIR/route.tpl" > "$BASE_DIR/route.go"

# MODULE
sed \
  -e "s|{{Feature}}|$FEATURE_CAPITALIZED|g" \
  -e "s|{{feature}}|$FEATURE_SANITIZED|g" \
  -e "s|{{Module}}|$GO_MODULE|g" \
  "$TEMPLATES_DIR/module.tpl" > "$BASE_DIR/module.go"


# Success message
GREEN='\033[0;32m'
NC='\033[0m'
echo -e "${GREEN}✅ Generated feature '$FEATURE' in $BASE_DIR with templates (resource=$IS_RESOURCE)${NC}"
