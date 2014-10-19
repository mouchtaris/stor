#!/bin/sh

find src -iname '*.go' -print0 | xargs -0 gofmt -w
