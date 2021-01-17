# Changelog

## v0.1.0

Initial Release

Working dynamodb backend with assumptions:

 - Must have both a hash key and range key on your primary partition
 - Assumed PK and SK for the key names but can be overridden with environment variables
 - Tests exist to prove functionality
