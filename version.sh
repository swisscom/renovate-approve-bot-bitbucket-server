#!/bin/bash
set -e
if [ "$VERSION" != "" ]; then
  echo -n "$VERSION"
  exit 0
fi
git describe --always --tags --dirty