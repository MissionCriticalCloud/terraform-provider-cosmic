#!/bin/bash

set -e

usage()
{
  echo "usage: release.sh [[-v version ] | [-h]]"
}


version=
while getopts v:h: option
do
  case "${option}"
  in
  v) version=${OPTARG};;
  h) usage;;
  esac
done

if [ "$version" == "" ]; then
  echo "Please specify version/tag to release in format X.X.X"
  exit 1
else
  # Add tag
  echo "git tag -a v${version} -m \"Release v$version\""
  git tag -a v${version} -m "Release v$version"
  # Push
  echo "git push origin v${version}"
  git push origin v${version}
  # Create release
  echo "goreleaser --rm-dist"
  goreleaser --rm-dist
fi