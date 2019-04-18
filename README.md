# gauth
minimal google sign-in prototype with typescript knockout and golang

Authenticates with Google Sign In and uses a session token for subsequent api call authorization. Readily useable in other projects

Implements the process described here: https://developers.google.com/identity/sign-in/web/backend-auth

In order to use this, create an Oauth2 client credentials in https://console.developers.google.com/apis/credentials
and complete the settings. You must use this client ID in the index.html and main.go instead of the provided one.
