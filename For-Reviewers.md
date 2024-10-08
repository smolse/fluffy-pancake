# For Reviewers

This project was developed with the help of the [Gin](https://github.com/gin-gonic/gin) web framework for Go and follows
the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) directory structure. The Gin
framework was chosen because of its simplicity, performance and popularity within the Golang community.

The entrypoint for the web service is in the [cmd/risk-api/main.go](./cmd/risk-api/main.go) file, while the rest of the
application code is in the [internal](./internal) directory.

Below I explain some assumptions and decisions made during the development of the initial version of this project.
Due to the given time constraints, several trade-offs were made and the current implementation is not completely
production-ready. The following list describes the most important points behind the implementation and provides
additional ideas for future improvements:

- Application architecture

    - Due to the time constraints and simplicity of the required business logic, for the purpose of not over-engineering
    the solution, I decided not to strictly follow any specific architectural pattern like Clean architecture or
    Hexagonal architecture and instead implemented a simple layered architecture. If the service is expected to
    significantly grow in complexity in the future, it could be beneficial to refactor the codebase to utilize a more
    advanced architectural approach.

- OpenAPI specification

    - I didn't create a Swagger/OpenAPI specification for the API endpoints, but it would be beneficial to have one. It
    would help document the API in a standardized way and provide a clear contract for clients to interact with the
    service. Furthermore, it could be used to generate client SDKs and server stubs, so that there would be less
    boilerplate code to write manually for routing, models and request/response validation. We could use tools like
    [openapi-generator](https://github.com/OpenAPITools/openapi-generator) to generate the service skeleton. However,
    since the goal of this assignment was to perform some hands-on Go coding, I decided to write things from scratch
    without using any code generation tools. In a real-world scenario, the OpenAPI specification would be really
    valuable and could be also used for deploying the API with Infrastructure as Code tooling. Such specification would
    live in the `/api` directory in the project root.

- Exponential backoff with jitter

    - Currently, the service layer ([service/service.go](service/service.go)) does not implement any retry mechanism in
    case of failures when interacting with the data store (not really needed for the local in-memory implementation,
    but would be essential for a real decoupled database). Implementing an exponential backoff with jitter retry
    mechanism would make the service more resilient to transient failures and improve the overall reliability of the
    service.

- Pagination for returning risks

    - Currently, the `GET /v1/risks` endpoint returns all risks at once. For the production-ready service, it would be
    crucial to implement pagination to limit the number of risks returned in a single response.

- Cache for specific risks

    - The current implementation does not include any caching mechanism. Adding a cache layer for specific risks
    returned by the `GET /v1/risks/{id}` endpoint could improve the performance of the service by reducing the number
    of database queries. This could be implemented using a key-value store like Redis or Memcached.

- Health check endpoint

    - While the task requirements didn't mention implemenetation of a health check endpoint (e.g., `/v1/health`), it's
    important to implement one for proper operationalization of the service. The health check endpoint can be used to
    monitor the service's availability, and it would allow for self-healing and auto-scaling mechanisms to be
    implemented, which are critical for running the service in a production environment.

- Graceful shutdown

    - I added a signal handler to catch `SIGINT` (sent by Ctrl+C) and `SIGTERM` (often sent by container orchestrators
    to stop or replace a container) signals to gracefully shut down the server. This allows the service to drain
    existing connections and finish ongoing requests instead of abruptly terminating them.

- Custom error types

    - The current implementation uses the standard Go `error` type for error handling. It would be beneficial to define
    custom error types to provide more context and information about the errors that occur in the service. This would
    make it cleaner to propagate errors throughout the service and provide better error messages to the clients.

- Models Improvement

    - The current implementation uses a simple `Risk` struct to represent a risk, which is used by both the service and
    storage layers. Depending on the complexity of the risks and the requirements of the service, it might be beneficial
    to separate the internal representation of a risk from the external representation;

    - Also, I didn't create structs for the API response models and instead used `gin.H` maps. It would be better to
    actually define formal response models, but again, this probably would be best achieved if an OpenAPI specification
    was used.

- Logs, metrics and tracing

    - The current implementation does not include any metrics or tracing collection capabilities. Adding metrics and
    tracing would allow for better observability of the service, which is crucial for monitoring and debugging purposes.

    - Logging was implemented using the simplistic standard Go `log` package. For a production-ready service, it would
    be better to make use of a more advanced logging library like `zap`, `logrus` or `zerolog` to be able to capture
    structured logs, have better control over log levels and other features.

- Testing

    - I added some unit tests here and there (e.g., for [service](./internal/service/), [handler](./internal/handler/)
    and [models](./internal/models/)), but the test coverage is far from ideal. Would be good to add more tests when
    the time permits.

    - Ideally, we can also implement end-to-end tests for this service. I would suggest using
    [godog](https://github.com/cucumber/godog), which is an official Cucumber BDD framework implementation for Golang.
    This code would live in the `/test` directory in the project root.

- CI/CD

    - I didn't create any GitHub Actions workflows in the initial implementation, but ideally this project needs
    at least 2 workflows:

        - One to execute when a pull request is opened or updated, which would run the tests and linter for the code
        changes;

        - Another to execute when a new tag is pushed, which would build the application binaries and Docker image,
        and then make them available in the release assets and GitHub Container Registry, respectively. Currently, the
        binaries are simply stored in the [bin](./bin/) directory in the repository.

- Linter

    - The current implementation does not use any linting tools during testing. Adding a linter like
    [golangci-lint](https://github.com/golangci/golangci-lint) would be nice to ensure code consistency and quality.
    It would ideally run prior to executing unit tests.
