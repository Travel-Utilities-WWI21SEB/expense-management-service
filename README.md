# expense-management-service

[![Continious Integration](https://github.com/Travel-Utilities-WWI21SEB/expense-management-service/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Travel-Utilities-WWI21SEB/expense-management-service/actions/workflows/ci.yml)
[![Continuous Delivery](https://github.com/Travel-Utilities-WWI21SEB/expense-management-service/actions/workflows/cd.yml/badge.svg?branch=main)](https://github.com/Travel-Utilities-WWI21SEB/expense-management-service/actions/workflows/cd.yml)
[![CodeScene Code Health](https://codescene.io/projects/39275/status-badges/code-health?component-name=Expense-management-service)](https://codescene.io/projects/39275/architecture/biomarkers?component=Expense-management-service)

Backend Service for the Expense Management Tool
 
## Developing

If you cloned the repository locally ensure that you have `go 1.20` installed. To install the dependencies run the following command:

```bash
go mod download
```

After installing the dependencies you need to populate a local `.env` file in the root directory with the following variables:

```bash
# DEV:    For local development
# DOCKER: For builds with docker-compose
# PROD:   For prod builds or previews
ENVIRONMENT=DEV or DOCKER or PROD

# Secret key for generating JWT tokens
JWT_SECRET="some super long and secure secret"

# For dev database connection
DEV_DB_HOST=localhost
DEV_DB_PORT=5432
DEV_DB_USER=admin
DEV_DB_PASSWORD=password
DEV_DB_NAME=travel_expenses-db

# For docker network database connection
DOCKER_DB_HOST=expense-db-container
DOCKER_DB_PORT=5432
DOCKER_DB_USER=admin
DOCKER_DB_PASSWORD=password
DOCKER_DB_NAME=travel_expenses-db

# For prod database connection
PROD_DB_HOST=ask luca
PROD_DB_PORT=ask luca
PROD_DB_USER=ask luca
PROD_DB_PASSWORD=ask luca
PROD_DB_NAME=ask luca

# Mail Server
MAILGUN_API_KEY=ask kevin
MAILGUN_DOMAIN=mail.costventures.works
```

For the dev-enviroment you will also need a local postgres cluster running on port 5432 with the credentials as specified in the example above. 
The database itself should be setup with the expense_db.sql file in the root directory of this repository. More information on how to do this
can be found [here](https://github.com/Travel-Utilities-WWI21SEB/expense-management-docs#datenbank).

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

To contribute to this template, clone this repository locally and commit your code on a separate branch. 
Try using conventional commits. When you're ready to submit your PR, check if your application builds by running the following commands:

```bash
go build ./...
```

Finally check if you can build and run the container:

```bash
bash ./scripts/start
```

If the checks pass open a pull request on the `main` branch. Once the pull request is approved and merged, your changes will be automatically deployed to production.
