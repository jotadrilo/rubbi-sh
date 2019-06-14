#!/bin/bash

echo STABLE_VERSION "$(cat VERSION)"
echo STABLE_COMMIT "$(git rev-parse HEAD)"
echo STABLE_DATE "$(date +%Y-%m-%d)"

