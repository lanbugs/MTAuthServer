#!/usr/bin/env bash

########################################################################################################################
APP_NAME="mtauthsrv"
VERSION_NAME="v1.0.3"

RELEASES_DIR="release"

TARGETS=("linux/amd64" "windows/amd64" "darwin/amd64" "darwin/arm64")
########################################################################################################################

$HOME/go/bin/swag init


if [ ! -d "$RELEASES_DIR" ]; then
  mkdir "$RELEASES_DIR"
fi

for target in "${TARGETS[@]}"; do
  IFS='/' read -ra parts <<< "$target"
  os="${parts[0]}"
  arch="${parts[1]}"

  target_dir="$RELEASES_DIR/$VERSION_NAME/$os/$arch"
  mkdir -p "$target_dir"

  env GOOS="$os" GOARCH="$arch" go build -o "$target_dir/$APP_NAME"

  if [ $? -eq 0 ]; then
    echo "Successful for $os/$arch compiled"
  else
    echo "Error on $os/$arch compilation"
  fi
done

