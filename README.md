# Keycloak REST Identity store for mxisd

This is a proof of concept of the [REST Identity store](https://github.com/kamax-matrix/mxisd/blob/master/docs/stores/rest.md) for Keycloak

### Endpoints
| Name           | Implemented |
|----------------|-------------|
| Authentication | Yes         |
| Directory      | Yes (Not fully tested yet) |
| Identity       | Yes         |
| Profile        | Yes (Not fully tested yet) |

# Keycloak Setup
> Using the `master` realm is not recomended
## Client Setup
You will need to create a new client in order to use this with

| Name                | Value        |
|---------------------|--------------|
| Access Type         | confidential |
| Valid Redirect URIs | https://matrix.example.com/* |

The `Client ID` and `Credentials`>`Secret` values are needed in the [Run](#running) step.

## User Setup
For this backend to work, it needs to be able to see and query all users available in the realm. In order to do so, you will need to setup a single user with the `view-users` role enabled in `Role Mappings`>`Client Roles`>`realm-management`>`Available Roles`


# Building
Install the dependencies:
```
go get golang.org/x/oauth2
```
Then compile the backend:
```
go build -o backend .
```

# Running
In order to run this backend you will need to setup the following environmental variables and save it to `.env`

```
CLIENT_ID=clientid
CLIENT_SECRET=clientsecret
KEYCLOAK_URL=http://local-ip-of-keycloak:port
KEYCLOAK_REALM=realm
KEYCLOAK_USER=user
KEYCLOAK_PASSWORD=password
KEYCLOAK_ATTRIBUTES_VALID_3PIDS="email,msisdn"
```
> `KEYCLOAK_URL` - The address should be an IP/host name that provides a direct connection to Keycloak.

Run with
```
source .env && ./backend
```

# Integration
To integrate this backend with [mxisd](https://github.com/kamax-matrix/mxisd), you will need to enable this REST backend in `mxisd.yaml` 

```
rest:
  enabled: true
  host: 'http://ip-of-this-backend:8091'
```
> `host` - The host should be an IP/host name that provides a direct connection to this backend. This MUST NOT be a public address, and SHOULD NOT go through a reverse proxy.


In order to be able to log in to your `matrix` homeserver using credentials stored in `Keycloak`, set up the Basic/Advanced [Authentication](https://github.com/kamax-matrix/mxisd/blob/master/docs/features/authentication.md) feature in `mxisd`

# Known Issues
- For now `Keycloak` usernames should comply with matrix username rules
- Users can have up to two `email` addresses and one `msisdn` (Phone number) as 3PIDs
