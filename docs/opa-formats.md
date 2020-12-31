# OPA input and output needed

This project is using OPA servers for authorizations. This document will show the representation of the input send by OPA center to OPA servers.

## Input

The input data will be a JSON and will have the following structure:

- a `user` key that will contains the user decoded information as a JSON:
  - `preferred_username`: username of connected user
  - `name`: name
  - `given_name`: given name
  - `family_name`: family name
  - `email`: email
  - `email_verified`: email verified boolean flag
- a `tags` key that will contains fixed tags configured (see [OPAServerAuthorizationConfiguration](configuration.md#opaserverauthorizationconfiguration))
- a `data` key that will contains the user action and on which resource
  - `action`: will contains the user action (See [Authorizations](authorizations.md) for more information)
  - `resource`: will contains the resource on which the user is trying the perform the action (See [Authorizations](authorizations.md) for more information)

Here is an example:

```json
{
  "user": {
    "preferred_username": "username",
    "name": "name",
    "given_name": "given name",
    "family_name": "family name",
    "email": "email",
    "email_verified": true
  },
  "tags": {
    "tag1": "value1"
  },
  "data": {
    "action": "action",
    "resource": "resource"
  }
}
```

## Output

The output data will be a JSON and will have the following structure:

- a `result` key that will contains a boolean that indicates if the user is authorized or not

Here is an example:

```json
{
  "result": true
}
```
