# MockItOut

[![Build 
Status](https://travis-ci.org/madflojo/mockitout.svg)](
https://travis-ci.org/madflojo/mockitout) [![Coverage Status](https://coveralls.io/repos/github/madflojo/mockitout/badge.svg?branch=master)](https://coveralls.io/github/madflojo/mockitout?branch=master)

Test external services faster and better with an HTTP stub server.

MockItOut is a simple to use HTTP stub server. With a small YAML configuration you can quickly create test end-points for HTTP APIs. Unlike other mock servers this project is language agnostic and easy to setup.

## Key Features

* HTTP response stubbing, maching URI with pre-canned body, header and status code replies.
* Logging request data for troubleshooting and diagnostics.
* Runs as a docker container or as a local binary.
* Callable as an external service for unit or functional tests.
* Simple YAML configuration.

## Running MockItOut

To run MockItOut is simple, just run the following commands.

```sh
$ docker run -p 443:8443 madflojo/mockitout:latest
```

This will start the service with our [example mock file](examples/hello_world.yml). To test it you can use `curl`.

```sh
curl -vk https://localhost/hi
```
### Specifying your own mocks file

To add your own mocks file, simply use volume mounts with the `docker run` command.

```sh
$ docker run -p 443:8443 -v stubs/:stubs -e MOCKS_FILE="stubs/mystubs.yml" madflojo/mockitout:latest
```

## Mocks Configuration File

To define end-points create a YAML file with the following format.

```yaml
routes:
  hello:
    path: "/hi"
    response_headers:
      "content-type": "application/json"
      "server": "MockItOut"
    # Multi-line values can be created like this
    body: | 
      {
        "greeting": "Hello",
        "name": "World"
      }
  deny:
    path: "/no"
    response_headers:
      "content-type": "application/json"
      "server": "MockItOut"
    body: |
      {"status": false}
    return_code: 403
  names:
    path: "/names/*"
    response_headers:
      "content-type": "application/json"
      "server": "WalkItOut"
    return_code: 200
    body: |
      {
        "1": {
          "name": "DJ Unk"
        },
        "2": {
          "name": "Andre 3000"
        },
        "3": {
          "name": "Jim Jones"
        }
      }
```

## Configuring with Environment Variables

MockItOut is controlled via environment variables. The below is a list of all environment variables available and what they control.

* `DEBUG` can be `true` or `false`. This will enable or disable debug logs. Default is `false`.
* `DISABLE_LOGGING` can be `true` or `false`. This will disable all logging. Default is `false`.
* `ENABLE_TLS` can be `true` or `false`. This will have the server use HTTPS by default. Default is `true`.
* `LISTEN_ADDR` defines the server listener address and port. Default is `0.0.0.0:8443`
* `CERT_FILE` defines the location of the TLS Certificate file.
* `KEY_FILE` defines the location of the TLS Certificate Key file.
* `GEN_CERTS` can be `true` or `false`. This will enable the server to create temporary testing certs on boot. Default is `true`.
* `MOCKS_FILE` defines the location of the mocks configuration file.


## Contributing
Thank you for your interest in helping develop MockItOut. The time, skills, and perspectives you contribute to this project are valued.

Please reference our [Contributing Guide](CONTRIBUTING.md) for details.

## License
[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)
