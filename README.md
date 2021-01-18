# terseurl

[![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/terseurl)](https://goreportcard.com/report/github.com/MicahParks/terseurl) [![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/terseurl.svg)](https://pkg.go.dev/github.com/MicahParks/terseurl)

Currently under development.

## Configuration

Environment variable table:

|Name               |Description                                                                                                                                                                    |Default Value                  |Example Value                                |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------|---------------------------------------------|
|`DEFAULT_TIMEOUT`  |The amount of time to wait before timing out for an incoming (client) or an outgoing (database) request in seconds.                                                            |`60`                           |`180`                                        |
|`INVALID_PATHS`    |A comma separated list of paths that cannot be assigned to a shortened URL. Whitespace prefixes and suffixes are trimmed. All swagger endpoints like `api` are invalid.        |swagger endpoints and frontend |`ready ,live, v2`                            |
|`SHORTID_PARANOID` |Indicate whether randomly generated short URLs should be checked to see if they are already in use. Any value sets the boolean to true. Empty for false.                       |blank                          |`true`                                       |
|`SHORTID_SEED`     |The seed to give the random shortened URL generator. Unsigned 64 bit integer. It is recommend to set this in a production setting.                                             |System clock                   |`2301015`                                    |
|`TEMPLATE_PATH`    |The full or relative path to the HTML template to use when a shortened URL is requested and JavaScript fingerprinting or social media link previews are on.                    |`socialMediaLinkPreview.gohtml`|`customTemplate.gohtml`                      |
|`TERSE_STORE_JSON` |The JSON formatted storage configuration for the TerseStore. If empty, it will try to read the file at `terseStore.json`. If not found it will use an in memory implementation.|blank                          |`{"type":"bbolt","bboltPath":"terse.bbolt"}` |
|`VISITS_STORE_JSON`|The JSON formatted storage configuration for the VisitsStore. If empty, it will try to read the file at `visitsStore.json`. If not found, visits will not be tracked.          |blank                          |`{"type":"bbolt","bboltPath":"visits.bbolt"}`|
|`WORKER_COUNT`     |The quantity of workers to use for incoming request async tasks like performing async database calls. Not the number of incoming requests that can be handled at one time.     |`4`                            |`10`                                         |

### JSON formatted storage configuration

TODO

## Deployment
```bash
touch terse.bbolt visits.bbolt
docker-compose up
```

## TODO

- [ ] Address TODOs.
- [ ] Social media link preview `inherit` mode that gets the Original URL and uses the meta tags for that.
- [ ] Write a utility that will export `.bbolt` to JSON.
- [ ] Implement `SHORTID_PARANOID`.
- [ ] Allow for shortened URLs of the form `{owner}/{shortened}` in `/api/write/{operation}` endpoint.
  - [ ] Only allow for random shortened URLs in top level.
- [ ] Implement fingerprinting with fingerprintjs, but remove HTML canvas extraction. Embed minified in single HTML
  template.
- [ ] Find all potential nil pointer dereferences due to data being stored in backend storage that does not conform to
  swagger spec.
- [ ] Implement Redis storage backend?
- [ ] Implement pluggable store interface implementations.
- [ ] Visit counts in TerseStore.
- [ ] Reimplement Mongo storage.
- [ ] Write tests.
- [ ] Write a good README.md.
- [ ] Move frontend to another repo.
- [ ] Implement JWT + JWKS authentication?
- [x] Implement social media link previews.
- [x] Implement `/api/import` endpoints.
- [x] Flag strategy.
- [x] Implement bbolt storage backend.
- [x] Configure storage backends via config files?
- [x] Change user created warnings to info.
- [x] Take away auth.
- [x] Move things to `/api`.
- [x] Add an `/upsert` endpoint.
- [x] Add an `/randomUltraSafe` endpoint.
