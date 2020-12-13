# terse URL

[![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/terse-URL)](https://goreportcard.com/report/github.com/MicahParks/terse-URL) [![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/terse-URL.svg)](https://pkg.go.dev/github.com/MicahParks/terse-URL)

Currently under development.

## Configuration

Environment variable table:

|Name                     |Description                                                                                                                                                               |Default Value|Example Value                         |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|--------------------------------------|
|`DEFAULT_TIMEOUT`        |The amount of time to wait before timing out for an incoming (client) or an outgoing (database) request in seconds.                                                       |`60`         |`180`                                 |
|`INVALID_PATHS`          |A comma separated list of paths that cannot be assigned to a shortened URL. Whitespace prefixes and suffixes are trimmed.                                                 |`api,`       |`api, ready ,live`                    |
|`SHORTID_PARANOID`       |Indicate whether randomly generated short URLs should be checked to see if they are already in use. Any value sets the boolean to true. Empty for false.                  |blank        |`true`                                |
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
- [ ] Implement `SHORTID_PARANOID`.
- [ ] Implement `/api/import` endpoints.
- [ ] Reimplement Mongo storage.
- [ ] Redis storage backend.
- [ ] Delete expired Terse.
- [ ] Social media link preview.
- [ ] Move frontend to another repo.
- [ ] A good README.md.
- [ ] Write tests.
- [ ] Implement JWT + JWKs authentication?
- [ ] Change storage configs to files?
- [ ] Change user created warnings to info?
- [x] Take away auth.
- [x] Move things to `/api`.
- [x] Add an `/upsert` endpoint.
- [x] Add an `/randomUltraSafe` endpoint.
