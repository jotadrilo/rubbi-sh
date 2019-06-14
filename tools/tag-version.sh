#!/bin/bash

root=$(git rev-parse --show-toplevel)
version="${root}/VERSION"
git tag "$(cat "$version")"
