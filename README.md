# Uploadly

Easily store and manage your files in the cloud

**No tests yet 😢**

[SwaggerHub](https://app.swaggerhub.com/apis/uploadly/uploadly)

## Set up
### Requirements
- Go 1.20 or higher
- SQLite3

### Configuration
Copy an instance of the configuration file and save it as a file `config.yaml`
```shell
cp config.dist.yaml config.yaml
```
... and configure as you need. All configuration instructions are given in the configuration example

### Launch

Install the dependencies in the vendor directory and run the application
```shell
go mod vendor
go build -v ./...
go run -v ./...
```
Enjoy!
