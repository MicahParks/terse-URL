# terseurl

[![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/terseurl)](https://goreportcard.com/report/github.com/MicahParks/terseurl) [![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/terseurl.svg)](https://pkg.go.dev/github.com/MicahParks/terseurl)

Currently under development.

## Terms

* *client*: Something using editor access to an instance of terseurl through the REST API.
* *user*: Something visiting an instance of terseurl with a web browser without any special access.
* *Visits data*: Data about *users* that have visited a particular shortened URL.
* *Terse data*: Data that describes the process of how to get from a shortened URL to an original URL. This data
  includes the shortened URL and original URL.

## Features

### Redirection with shortened URLs

Create web browser redirections to original URLs through shortened URLs. Shortened URLs are unique URL safe strings. A
*client* may provide one, or the server can generate one.

*Example:* A *client* created *Terse data* with the shortened URL, `myblog`. The *Terse data* has the original URL of
`http://example.com/blogs/my/1`. The link `https://terseurl.com/myblog` is shared with other *users*. When a *user*
visits `https://terseurl.com/myblog`, their web browser will redirect them to `http://example.com/blogs/my/1`.

### Multiple redirection types

Currently, the project supports the following redirection types:

* HTTP 301
* HTTP 302
* HTML `<meta>`
* JavaScript

If there are more redirection types (that are widely accepted by web browsers) suggest them to the developers.

### Social media link previews

If *Terse data* is configured to perform a redirect via HTML `<meta>` tags or JavaScript, there is the option to add
social media link previews. This is done by adding HTML `<meta>` tags for [Open Graph](https://ogp.me) and
[Twitter](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/markup).

This can be added manually to *Terse data*. It can also be inherited from the original URL by using an API endpoint, or
a button on the frontend.

### *Visits data*

By default, the project will not keep track of *Visits data*. If the project is configured to, it can track visits to
shortened URLs. All gathered *Visits data* is placed in backend storage and accessible via the web frontend or API.

The types of *Visits data* collected can vary. It can include IP address, HTTP headers, and information gathered form
JavaScript.

### Control *Terse data* and *Visits data*

*Terse data* and *Visits data* is accessible through the web interface and API. Data can easily be imported and exported
in JSON format via the frontend and API endpoints. Data can also be interacted with directly via the frontend or other
API endpoints.

### Customizable storage options

Currently, the project natively supports these storage backends:

* memory
* bbolt (file on disk)

However, the project can support any storage backend that implements its respective storage interface. TODO

## Configuration

Environment variable table:

|Name                 |Description                                                                                                                                                                                              |Default Value                  |Example Value                                                                    |
|---------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------|---------------------------------------------------------------------------------|
|`DEFAULT_TIMEOUT`    |The amount of time to wait before timing out for an incoming (client) or an outgoing (database) request in seconds.                                                                                      |`60`                           |`180`                                                                            |
|`FRONTEND_STATIC_DIR`|The path to the directory that contains the static frontend assets to be served out of `/frontend/*`. If empty, the embedded assets will be used.                                                        |blank                          |`./frontend2`                                                                    |
|`HTTP_PREFIX`        |The HTTP prefix all shortened URLs will have. This is used by the frontend.                                                                                                                              |`https://terseurl.com/`        |`https://example.com/`                                                           |
|`INVALID_PATHS`      |A comma separated list of paths that cannot be assigned to a shortened URL. Whitespace prefixes and suffixes are trimmed. All swagger endpoints like `api` are invalid.                                  |swagger endpoints and frontend |`ready ,live, v2`                                                                |
|`JWKS_URL`           |The full URL to the Java Web Key Store where trusted JWTs are signed from. Only functional if `AUTH` is `true`                                                                                           |blank                          |`http://keycloak.terseurl.com/auth/realms/terseurl/protocol/openid-connect/certs`|
|`SHORTID_PARANOID`   |Indicate whether randomly generated short URLs should be checked to see if they are already in use. Any value except for `true` sets the boolean to false.                                               |blank                          |`true`                                                                           |
|`SHORTID_SEED`       |The seed to give the random shortened URL generator. Unsigned 64 bit integer. It is recommend to set this in a production setting.                                                                       |System clock                   |`2301015`                                                                        |
|`TEMPLATE_PATH`      |The full or relative path to the HTML template to use when a shortened URL is requested and JavaScript fingerprinting or social media link previews are on. If empty, the embedded template will be used.|`redirect.gohtml`              |`customTemplate.gohtml`                                                          |
|`USE_AUTH`           |Turn authentication and authorization on or off. Any value except for `true` sets the boolean to false.                                                                                                  |blank                          |`true`                                                                           |
|`SUMMARY_STORE_JSON` |The JSON formatted storage configuration for the SummaryStore. If empty, it will try to read the file at `summaryStore.json`. If not found it will use an in memory implementation.                      |blank                          |`{"type":"memory"}`                                                              |
|`TERSE_STORE_JSON`   |The JSON formatted storage configuration for the TerseStore. If empty, it will try to read the file at `terseStore.json`. If not found it will use an in memory implementation.                          |blank                          |`{"type":"bbolt","bboltPath":"terse.bbolt"}`                                     |
|`VISITS_STORE_JSON`  |The JSON formatted storage configuration for the VisitsStore. If empty, it will try to read the file at `visitsStore.json`. If not found, visits will not be tracked.                                    |blank                          |`{"type":"bbolt","bboltPath":"visits.bbolt"}`                                    |
|`WORKER_COUNT`       |The quantity of workers to use for incoming request async tasks like performing async database calls. Not the number of incoming requests that can be handled at one time.                               |`4`                            |`10`                                                                             |

### JSON formatted storage configuration

TODO

## Deployment

To deploy terseurl for local development, follow the below instructions. These can be adapted for production. It is
recommended to use embedded assets in production.

### `docker-compose` development

The following command will create a fresh instance of terseurl hosted via HTTPS exposing ports `80` and `443`.

If using bbolt, the files must be created before a `docker-compose` starts the service. If the files are not present at
this time, `docker-compose` will create the files as directories.

```bash
rm -rf terse.bbolt visits.bbolt
touch terse.bbolt visits.bbolt
docker-compose up
```

### Local development

```bash
HOST=0.0.0.0 PORT=30000 FRONTEND_STATIC_DIR=frontend TEMPLATE_PATH=redirect.gohtml USE_AUTH=false go run -race cmd/terseurl-server/main.go
```

## TODO

- [ ] Address source code TODOs.
- [ ] Inherit HTML title.
- [ ] Make write operations atomic?
- [ ] Implement an Authorization store.
- [ ] Update deployment instructions and `docker-compose.yml` for auth.
- [ ] Completely remove ErrShortenedNotFound? Use zero values to communicate that?
- [ ] Change bbolt data structure for Visits to something more scalable.
- [ ] Truncate frontend table data so it doesn't run off the screen.
- [ ] Add more logic to rate limiter for frontend use case.
- [ ] Copy to clipboard button for shortened URL.
- [ ] Show full shortened URL in table data.
- [ ] Hyperlinks for shortened URL and original URL.
- [ ] Add referer URL to query parameters?
- [ ] Implement `SHORTID_PARANOID` environment variable.
- [ ] Implement JavaScript tracking.
- [ ] Implement JavaScript fingerprinting.
  - [ ] Remove things that break on Firefox by default, like canvas extraction?
  - [ ] Embed into redirect template through variable.
- [ ] Implement frontend code for `USE_AUTH`. May require endpoint.
- [ ] Expand Visits model to include untyped data for custom VisitsStore implementation or interceptor.
- [ ] Create Visits data visualizer.
- [ ] Create profile viewer.
- [ ] Move that GitHub button to the collapsable side bar.
- [ ] Visit data interceptor for data purging or whatever before it goes to backend storage.
- [ ] Implement `SHORTID_PARANOID`.
- [ ] Allow for shortened URLs of the form `{shortened}/{extended}` in `/api/write/{operation}` endpoint.
  - [ ] Only allow random shortened URLs in top level?
- [ ] Implement Redis storage backend?
- [ ] Implement custom data stores, interceptor, and HTTP middlewares modules the same
  way [Caddy does](https://caddyserver.com/docs/extending-caddy).
- [ ] Reimplement Mongo storage.
- [ ] Write tests.
- [ ] Write a good README.md.
