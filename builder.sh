#!/bin/bash
set -e #exit on error
clear 
echo "building backend ..."

go build -o ./bin/gauth ./src

#!/bin/sh
set -e #exit if error
cd ./src/ui/
echo "building frontend ..."
tsc --strict
cd - > /dev/null
echo "copying files ..."
mkdir -p ./bin/ui
cp -R src/ui/built ./bin/ui
cp -R src/ui/externals ./bin/ui
# cp -R src/ui/favicon ./bin/ui
cp src/ui/*.html ./bin/ui
# cp src/ui/*.css ./bin/ui
# cp src/ui/*.js ./bin/ui
echo "SUCCESS"