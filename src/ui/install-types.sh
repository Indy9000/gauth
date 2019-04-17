#!/bin/bash
set -e #exit if error

mkdir built
#knockout types
npm install --save @types/knockout

#google sign in api types
npm install --save @types/gapi.auth2
npm install --save @types/gapi

mkdir externals
cd externals
wget https://cdnjs.cloudflare.com/ajax/libs/knockout/3.5.0/knockout-min.js
wget https://requirejs.org/docs/release/2.3.6/minified/require.js
