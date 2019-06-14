#!/bin/bash

root=$(git rev-parse --show-toplevel)
version="${root}/VERSION"
major=$(cut -d. -f1 "$version")
minor=$(cut -d. -f2 "$version")
patch=$(cut -d. -f3 "$version")

new_version_string="$major.$minor.$((patch+1))"
echo "$new_version_string" >"$version"
git commit -am "Bump VERSION to $new_version_string"
