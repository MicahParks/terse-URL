# terse URL

[![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/terse-URL)](https://goreportcard.com/report/github.com/MicahParks/terse-URL) [![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/terse-URL.svg)](https://pkg.go.dev/github.com/MicahParks/terse-URL)

Currently under development.

## Configuration

Environment variable table:

|Name                     |Description                                                                                                                                                               |Default Value    |Example Value                         |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------|--------------------------------------|
|`DEFAULT_TIMEOUT`        |The amount of time to wait before timing out for an incoming (client) or an outgoing (database) request in seconds.                                                       |`60`             |`180`                                 |
|`INVALID_PATHS`          |A comma separated list of paths that cannot be assigned to a shortened URL. Whitespace prefixes and suffixes are trimmed. All swagger endpoints like `api` are invalid.   |swagger endpoints|`ready ,live, v2`                     |
|`SHORTID_PARANOID`       |Indicate whether randomly generated short URLs should be checked to see if they are already in use. Any value sets the boolean to true. Empty for false.                  |blank            |`true`                                |
|`SHORTID_SEED`           |The seed to give the random shortened URL generator. Unsigned 64 bit integer. It is recommend to set this in a production setting.                                        |System clock     |`2301015`                             |
|`TERSE_STORE_TYPE`       |The type of storage backend for Terse pairs.                                                                                                                              |`memory`         |`memory`                              |
|`VISITS_STORE_TYPE`      |The type of storage backend for Visits. Leave blank to not track visits.                                                                                                  |blank            |blank or `memory`                     |
|`WORKER_COUNT`           |The quantity of workers to use for incoming request async tasks like performing async database calls. Not the number of incoming requests that can be handled at one time.|`4`              |`10`                                  |

## TODO

- [ ] Address TODOs.
- [ ] Delete expired Terse.
- [ ] Implement `SHORTID_PARANOID`.
- [ ] Implement `/api/import` endpoints.
- [ ] Implement social media link previews.
- [ ] Implement Redis storage backend.
- [ ] Reimplement Mongo storage.
- [ ] Write tests.
- [ ] Write a good README.md.
- [ ] Move frontend to another repo.
- [ ] Implement JWT + JWKs authentication?
- [ ] Configure storage backends via config files?
- [x] Change user created warnings to info.
- [x] Take away auth.
- [x] Move things to `/api`.
- [x] Add an `/upsert` endpoint.
- [x] Add an `/randomUltraSafe` endpoint.
