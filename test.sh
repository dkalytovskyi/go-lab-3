#!/usr/bin/env sh

cd ./integration || return
CGO_ENABLED=0 bood
cat ./out/bin/test.txt
