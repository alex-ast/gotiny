Implementaion of TinyURL-like REST API service. 
Uses Redis for caching, MongoDB as database, Prometheus for matrics collection and Docker for containerization and for the build.

Supports Create, Delete, Get operations.

# Build: 

Prerequisites: [make](https://www.gnu.org/software/make), [go](https://go.dev) and [Docker](https://www.docker.com)

List of available targets:


- `make local`    # Default target. Builds on the local machine, built binaries are in app/
- `make start`    # starts Redis, MongoDB and API and http server locally
...
- `make stop`     # stops Redis/Mongo.

Other targets:
- `make docker`   # builds docker container for API server
- `make gen`      # generates [API markdown](api/gotiny_api_spec.md) and API structures out of [API spec](api/gotiny-api.yaml)

<hr />

### Build and run

`make local && make start`

### Create

`curl -X POST http://localhost:81/api/url -d "{\"longUrl\":\"http://alex.astapchuk.com\"}"`

Ouput:

    {
        "status": {
            "success": true
        },
        "urlInfo": {
            "longUrl": "http://alex.astapchuk.com",
            "shortId": "hQ9baMJ"
        }
    }


### Retrieve

`curl http://localhost:81/api/url/hQ9baMJ`

    {
        "source": "db",
        "status": {
            "success": true
        },
        "urlInfo": {
            "longUrl": "http://alex.astapchuk.com",
            "shortId": "hQ9baMJ"
        }
    }

Redirect:

`curl -i --max-redirs 0 http://localhost/hQ9baMJ`

    HTTP/1.1 303 See Other
    Access-Control-Allow-Origin: *
    Content-Type: text/html; charset=utf-8
    Location: http://alex.astapchuk.com
    Content-Length: 52

    <a href="http://alex.astapchuk.com">See Other</a>.

## Retrieve bypassing cache (debug mode only)

`curl http://localhost:81/api/url/hQ9baMJ?cache=0`

    {
        "source": "db",
        "status": {
            "success": true
        },
        "urlInfo": {
            "longUrl": "http://alex.astapchuk.com",
            "shortId": "hQ9baMJ"
        }
    }

### Delete

`curl -X DELETE http://localhost:81/api/url/hQ9baMJ`

    {
        "shortId": "hQ9baMJ",
        "status": {
            "success": true
        }
    }

Check it's really deleted:

`curl http://localhost:81/api/url/hQ9baMJ`

    {
        "message": "URL not found"
    }