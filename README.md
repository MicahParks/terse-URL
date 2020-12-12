# terse URL

Currently under development.

## Configuration

Environment variable table:

|Name                     |Description                                                                                                                                                               |Default Value|Example Value                         |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|--------------------------------------|
|`DEFAULT_TIMEOUT`        |The amount of time to wait before timing out for an incoming (client) or an outgoing (database) request in seconds.                                                       |`60`         |`180`                                 |
|`INVALID_PATHS`          |A comma separated list of paths that cannot be assigned to a shortened URL. Whitespace prefixes and suffixes are trimmed.                                                 |`api,`       |`api, ready ,live`                    |
|`KEYCLOAK_BASE_URL`      |The base URL for the Keycloak server. It has the HTTP prefix, hostname, and port.                                                                                         |blank        |`http://keycloak:8080`                |
|`KEYCLOAK_ID`            |The ID of the Keycloak client that is a service account.                                                                                                                  |blank        |`terseBackend`                        |
|`KEYCLOAK_REALM`         |The realm Keycloak uses for this service.                                                                                                                                 |blank        |`terseURL`                            |
|`KEYCLOAK_SECRET`        |The secret for the service account's Keycloak client.                                                                                                                     |blank        |`123e4567-e89b-12d3-a456-426614174000`|
|`SHORTID_SEED`           |The seed to give the random shortened URL generator. Unsigned 64 bit integer.                                                                                             |System clock |`2301015`                             |
|`TERSE_MONGO_COLLECTION` |The MongoDB collection to store Terse pairs in.                                                                                                                           |`terseStore` |`terseStore`                          |
|`TERSE_MONGO_DATABASE`   |The MongoDB database used to store the Terse collection in. Default `terseURL`.                                                                                           |`terseURL`   |`terseURL`                            |
|`TERSE_MONGO_URI`        |The MongoDB URI for the MongoDB server used for the database with the Terse collection in it.                                                                             |blank        |`mongodb://mongodb0.example.com:27017`|
|`TERSE_STORE_TYPE`       |The type of storage backend for Terse pairs.                                                                                                                              |`memory`     |`memory` or `mongo`                   |
|`VISITS_MONGO_COLLECTION`|The MongoDB collection to store Visits in.                                                                                                                                |`visitsStore`|`visitsStore`                         |
|`VISITS_MONGO_DATABASE`  |The MongoDB database to store the Visits collection in.                                                                                                                   |`terseURL`   |`terseURL`                            |
|`VISITS_MONGO_URI`       |The MongoDB URI for the MongoDB server used for the database with the Visits collection in it.                                                                            |blank        |`mongodb://mongodb0.example.com:27017`|
|`VISITS_STORE_TYPE`      |The type of storage backend for Visits. Leave blank to not track visits.                                                                                                  |blank        |blank or `memory` or `mongo`          |
|`WORKER_COUNT`           |The quantity of workers to use for incoming request async tasks like performing async database calls. Not the number of incoming requests that can be handled at one time.|`4`          |`10`                                  |

## TODO

- [ ] Address TODOs.
- [ ] Move things to `/api`?
- [ ] Add an `/randomUltraSafe` endpoint?
- [ ] Add an `/upsert` endpoint?
- [ ] Social media link preview.
- [ ] Move frontend to another repo.
- [ ] Change storage configs to files?
- [ ] Take away auth?
- [ ] Change user created warnings to info?
- [ ] A good README.md.
- [ ] Write tests.
