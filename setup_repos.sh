#!/bin/bash
set -e

PRESETS=(
    "backstage-cpp"
    "backstage-go"
    "backstage-java"
    "backstage-js"
    "backstage-kitchen-sink"
    "backstage-kotlin"
    "backstage-minimal"
    "backstage-py"
    "backstage-rust"
    "backstage-shell"
)

ORG="BlueCentre"

for preset in "${PRESETS[@]}"; do
    REPO_NAME="$ORG/$preset"
    
    # Convert preset to secret name format (uppercase, hyphens to underscores)
    SECRET_SUFFIX="${preset//-/_}"
    SECRET_NAME="STARTER_DEPLOY_${SECRET_SUFFIX^^}"
    
    echo "----------------------------------------------------------------"
    echo "Processing $preset..."
    echo "Target Repo: $REPO_NAME"
    echo "Secret Name: $SECRET_NAME"

    # 1. Create Repository if it doesn't exist
    if gh repo view "$REPO_NAME" &>/dev/null; then
        echo "✅ Repository $REPO_NAME already exists."
    else
        echo "Creating repository $REPO_NAME..."
        # Create as public, empty repo
        if ! gh repo create "$REPO_NAME" --public --description "Aspect Workflows Template for $preset"; then
            echo "❌ Failed to create repository $REPO_NAME. Check your permissions."
            echo "   You need 'repo' scope and permissions to create repositories in $ORG."
            exit 1
        fi
        echo "✅ Created $REPO_NAME"
    fi

    # 2. Generate SSH Key
    KEY_FILE="/tmp/id_ed25519_$preset"
    rm -f "$KEY_FILE" "$KEY_FILE.pub"
    ssh-keygen -t ed25519 -C "delivery-bot@aspect-workflows" -f "$KEY_FILE" -N "" -q
    echo "🔑 Generated SSH key pair"

    # 3. Add Deploy Key to Target Repo
    echo "Adding deploy key to $REPO_NAME..."
    if ! gh repo deploy-key add "$KEY_FILE.pub" -R "$REPO_NAME" --allow-write --title "Delivery Pipeline Key $(date +%Y-%m-%d)"; then
         echo "❌ Failed to add deploy key to $REPO_NAME."
         exit 1
    fi
    echo "✅ Deploy key added"

    # 4. Add Secret to Current Repo
    # Hardcoded to ensure we target the correct repo for secrets
    CURRENT_REPO="BlueCentre/aspect-workflows-template"
    echo "Setting secret $SECRET_NAME in $CURRENT_REPO..."
    if ! gh secret set "$SECRET_NAME" < "$KEY_FILE" -R "$CURRENT_REPO"; then
        echo "❌ Failed to set secret in $CURRENT_REPO."
        exit 1
    fi
    echo "✅ Secret $SECRET_NAME set"

    # Cleanup
    rm -f "$KEY_FILE" "$KEY_FILE.pub"
done

echo "----------------------------------------------------------------"
echo "🎉 All Backstage presets configured!"
