#!/bin/bash

root=$(git rev-parse --show-toplevel)
version="${root}/VERSION"
major=$(cut -d. -f1 "$version")
minor=$(cut -d. -f2 "$version")
patch=$(cut -d. -f3 "$version")

echo "$major.$minor.$((patch+1))" >"$version"
