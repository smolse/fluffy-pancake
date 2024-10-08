# fluffy-pancake

This project implements a simple RESTful web service for risk management. The service currently exposes the following
API endpoints, as specified in the assignment:
- `GET /v1/risks`: Returns a list of all risks.
- `GET /v1/risks/{id}`: Returns a specific risk by its ID.
- `POST /v1/risks`: Creates a new risk.

See the [For-Reviewers.md](./For-Reviewers.md) document for more information on the assumptions and decisions made
during the development of this project.

## Requirements

This application has been primarily intended to run in a Linux environment, though it should work in other major
operating systems (such as Windows and MacOS, though it wasn't thoroughly tested in them) as well without any issues.
If you want to run the compiled application natively on your machine, no additional dependencies should be required. If
you want to run it in a container, you will need to have the Docker Engine installed.

GNU Make v4 is used to orchestrate the testing and building process, which is expected to be executed in a Unix-like
environment.

## Installation

There are multiple ways to obtain executables for this application. Compiled binaries for AMD64 Linux, MacOS, and
Windows platforms are available in the [bin](./bin/) directory. You may need to do `chmod +x` to make them executable
before running them.

If you want to build the application from source for your specific OS and CPU architecture, you can execute the
following command to build the application on your local machine (need to have the Go toolchain installed):

```shell
$ make GOOS=linux GOARCH=amd64
```

You can also build binaries for your host machine within a Docker container, which is especially useful for ensuring reproducibility across different environments or if Go is not installed locally but you have access to Docker:

```shell
$ make build-in-docker GOOS=linux GOARCH=amd64
```

Compiled binaries will be available in the [bin](./bin/) directory.

If you actually want to run the application in a container, the Docker image can be built by running the following
command:

```shell
$ make build-docker
```

This will create a Docker image named `risk-api` with the latest tag.

All of these targets will run unit tests prior to building the application, though it's also possible to run them
separately by executing the `make test` target.  The Makefile also includes a `clean-docker` target to remove the
generated Docker images from your local Docker registry.

## Usage

Once you have the compiled binary or Docker image in place, you can start the web service in either of the following
ways:

```shell
$ ./bin/risk-api-linux-amd64
```

or

```shell
$ docker run -p 8080:8080 risk-api
```

### Configuration

The application uses environment variables for configuration. The following environment variables are currently
available:

Name                               | Default Value | Description
-----------------------------------|---------------|-----------------
`SERVER_PORT`                      | `8080`        | The port on which the server listens
`SERVER_GRACEFUL_SHUTDOWN_TIMEOUT` | `5s`          | How long the server waits for active connections to finish before shutting down
`DATASTORE_TYPE`                   | `syncmap`     | The type of the data store to use (`syncmap` is the only supported implementation for now)
`GIN_MODE`                         | `release`     | The mode in which the Gin framework runs (`release` or `debug`)

### Running the application

Assuming the application is running on the default port (8080), you can interact with the API using `curl` (optionally
with `jq` for pretty-printing JSON responses):

```
➜  fluffy-pancake git:(main) ✗ ./bin/risk-api-linux-amd64
2024/10/08 13:15:06 Starting server on :8080
```

In another terminal, first list all risks (should be empty initially):

```
➜  fluffy-pancake git:(main) ✗ curl -s localhost:8080/v1/risks | jq
{
  "risks": []
}
```

Create a couple of risks:

```
➜  fluffy-pancake git:(main) ✗ curl -s -XPOST localhost:8080/v1/risks -d '{"state": "open", "title": "foo", "description": "bar"}' | jq
{
  "id": "6729f259-e7bd-4d09-ba2a-761e30871e1f"
}
```

```
➜  fluffy-pancake git:(main) ✗ curl -s -XPOST localhost:8080/v1/risks -d '{"state": "open", "title": "foo", "description": "bar"}' |jq
{
  "id": "e8452a33-038a-439a-94e6-8f23d8213da4"
}
```

List all risks again:

```
➜  fluffy-pancake git:(main) ✗ curl -s localhost:8080/v1/risks | jq                                                               
{
  "risks": [
    {
      "id": "6729f259-e7bd-4d09-ba2a-761e30871e1f",
      "state": "open",
      "title": "foo",
      "description": "bar"
    },
    {
      "id": "e8452a33-038a-439a-94e6-8f23d8213da4",
      "state": "open",
      "title": "foo",
      "description": "bar"
    }
  ]
}
```

Retrieve a specific risk:

```
➜  fluffy-pancake git:(main) ✗ curl -s localhost:8080/v1/risks/6729f259-e7bd-4d09-ba2a-761e30871e1f | jq
{
  "risk": {
    "id": "6729f259-e7bd-4d09-ba2a-761e30871e1f",
    "state": "open",
    "title": "foo",
    "description": "bar"
  }
}
```

Try to retrieve a non-existent risk:

```
➜  fluffy-pancake git:(main) ✗ curl -s localhost:8080/v1/risks/00000000-0000-0000-0000-000000000000 | jq
{
  "error": "risk 00000000-0000-0000-0000-000000000000 was not found"
}
```

Try to create a risk with an invalid state:

```
fluffy-pancake git:(main) ✗ curl -s -XPOST localhost:8080/v1/risks -d '{"state": "test"}' | jq
{
  "error": "invalid risk state: test"
}
```
