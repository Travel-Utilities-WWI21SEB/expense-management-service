# expense-management-service

Backend Service for the Expense Management Tool

## Developing

If you cloned the repository locally ensure that you have `go 1.20` installed. To install the dependencies run the following command:

```bash
go mod download
```

## Building

To run the application in a local container run the following command:

```bash
bash ./scripts/start
# exposed at localhost:8081
```

To stop the container run the following command:

```bash
bash ./scripts/stop
```

To run the application locally without a container run the following command:

```bash
go run ./...
# exposed at localhost:8080
```

## Contributing

To contribute to this template, clone this repository locally and commit your code on a separate branch. Try using conventional commits. When you're ready to submit your PR, check if your application builds by running the following commands:

```bash
go build ./...
```

Finally check if you can build and run the container:

```bash
bash ./scripts/start
```

If the checks pass open a pull request on the `main` branch. Once the pull request is approved and merged, your changes will be automatically deployed to production.
