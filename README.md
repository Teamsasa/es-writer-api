# go-template
This repository is a template for Echo applications following the Clean Architecture principles.

## Usage
1. Create a `.env` file and set the environment variables for the database connection.</br>Example:
    ```env
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=postgres
    DB_PORT=5432
    ```
2. Run `make all` to start the application and initialize the database.