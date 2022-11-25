# Go Coverage

This tool is used to export summary of XML coverage report into a SQL database.

## Usage

```bash
export ANALYTICS_DATABASE_HOST="localhost"
export ANALYTICS_DATABASE_PORT="5433"
export ANALYTICS_DATABASE_USERNAME="postgres"
export ANALYTICS_DATABASE_PASSWORD=""
export ANALYTICS_DATABASE_NAME="analytics"
export DRONE_REPO="openware/go-coverage"
export ANALYTICS_COMPONENT="go-coverage"        # use it to differenciate different applications in a mono-repo.
export DRONE_TAG="1.0.0"                        # trigger this script on drone tag to track only stable versions

go run ./ coverage.xml
```
