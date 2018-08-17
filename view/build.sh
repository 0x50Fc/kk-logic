#!/bin/sh

if [ ! -d bin ]; then
    mkdir bin
fi

uglifyjs kk.js kk.extend.js kk.app.js \
    -o ./bin/kk.min.js
