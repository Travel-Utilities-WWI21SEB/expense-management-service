# expense-management-service

Backend Service for the Expense Management Tool

## Developing

If you cloned the repository locally ensure that you have `go 1.20` installed. To install the dependencies run the following command:

```bash
go mod download
```

After installing the dependencies you need to populate a local `.env` file in the root directory with the following variables:

````bash
# The port that the application will be exposed on
ENVIROMENT=DEV or PROD

# For DEV-enviroment only
DEV_DB_HOST=localhost
DEV_DB_PORT=5432
DEV_DB_USER=admin
DEV_DB_PASSWORD=password
DEV_DB_NAME=travel_expenses-db

# For PROD-enviroment only
PROD_DB_HOST=ask luca
PROD_DB_PORT=ask luca
PROD_DB_USER=ask luca
PROD_DB_PASSWORD=ask luca
PROD_DB_NAME=ask luca

For the dev-enviroment you will also need a local postgres cluster running on port 5432 with the credentials as specified in the example above. The database itself should be setup with the expense_db.sql file in the root directory of this repository.

It is planned to provide a docker-compose-file in the future to simplify the setup of the local development enviroment.

## Building

To run the application in a local container run the following command:

```bash
bash ./scripts/start
# exposed at localhost:8081
````

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
