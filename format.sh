#!/bin/sh

git checkout master -- .
git reset .
find src -iname '*.go' -print0 | xargs -0 gofmt -w
