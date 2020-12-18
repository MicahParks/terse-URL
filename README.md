# shakesearch

## Configuration

Environment variable table:

|Name                |Description                                                                                |Default Value      |Example Value            |
|--------------------|-------------------------------------------------------------------------------------------|-------------------|-------------------------|
|`SHAKESPEARES_WORKS`|The full or relative file path to a text file containing the complete works of Shakespeare.|`completeworks.txt`|`/home/william/works.txt`|

## TODO
- [ ] Dockerfile + docker-compose.
- [ ] Caddyfile + HTTPS + deploy.
- [ ] Highlight results in `app.js`.
- [x] Configure txt file location.
- [x] Handle errors.
- [x] Zap logger.
- [x] Swagger spec.
- [x] Fuzzy matching.
- [x] Line 79 slice out of bounds lol.
- [x] Rate limiter.
