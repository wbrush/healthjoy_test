# go-template-service
This repo is intended to be a template for Go-based microservices. It shows the basic implementation of a REST-based
server with 3 endpoints. These endpoints should be considered "standard" and should be implemented for any future
microservice.
Also this service shows how to implement CRUD operations with a simple _template_ model, how to organize    
the service structure, how to use DAOs with postgresql example realisation. 

# Development Notes

## Unit Testing
Unit tests have been written for the packages in this service. Optii Solutions goal is 80% coverage for unit tests in new development.

## Developer Notes
This should be modified on a per project basis. I would expect this to be general notes from one developer to the next developer who is tasked with looking at this service.

## Required Packages
Any package (or library) that this service uses that requires a "go get" to use MUST be listed here. Any external package that this service uses should be listed here. 

### Swagger
go get -u github.com/go-swagger/go-swagger/cmd/swagger

### Gorilla
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/handlers

### logrus
go get -u github.com/sirupsen/logrus

### go-pg
go get -u github.com/go-pg/pg 
go get -u github.com/go-pg/migrations

# Deployment
Need to document (and/or update) our deployment process here

## Process
 - Pull latest version of the service
 - Verify last deployed version (in the /api/help/ swagger endpoint)
 - Bump the API version for the release
 - Update vendored dependencies
 - Run "go generate" and do a test build locally
 - If successful, check in to github, develop branch
 - Once built and deployed, verify operation (including API version) in dev space
 - In bitbucket, perform PR and merge to master branch
 - Once built and deployed, verify operation (including API version) in QA space
 - In bitbucket, select "Create Release" and enter the release information
 - Once built and deployed, verify operation (including API version) in production

# Services
Every microservice support integrated modules system. Every module located at services directory and must to
implement the *Module* interface.
Every service must to take an arguments for dependency injection 
(like DB connection or config singleton) on New-like constructor func.
When module must run at global service initialisation process, put an instance of your module to 
run list at the params of *RunModules* func call. 

## API module
    the /api directory contains an auto-generated swagger document for the template service.
    there is also code that runs a webdave server on port 8080 with the swagger ui served at
    http://localhost:8080/api/help/index.html
or 
    http://localhost:8080/api/help/

### ToDo
 - Increase unit test coverage in the server package
 
## Daemons module
This is optional module located in template service for example. Daemons module is for periodically 
running background tasks. For example, when you need to expire some records after
some time, it is a good idea to run periodical checking function for records expiration.

Every daemon must to implement *Daemon* interface and must be located inside services/daemons directory.
When daemon must run at global service initialisation process, put an instance of your daemon to 
run list of workers at the *NewDaemons* func. 

# Go-based Development Notes

## Bitbucket Configuration
### "go get" a repo from Bitbucket
Set up bitbucket to use SSH. I think it is required for GO to be able to access the bitbucket private repos.

Configure git to use SSH when accessing bitbucket repos 
 - git config --global user.name "wayne_brush"
 - git config --global user.email "wayne.brush@optiisolutions.com" 
 - git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
 - go get bitbucket.org/optiisolutions/go-template-service

Note: that the ".git" is not included in the command. It will work if the ".git" is included but the directory will not
be named correctly which can affect future import statements.

For new microservices, it is recommended that the repo be created in bitbucket and then the files manually copied from this project to the new one.

## Folder Structure
This repo exemplifies the desired structure for Go-based microservices. Functionality and Features should be broken up
into manageable packages (subfolders) under the root folder. This should allow for easily adding, removing, and editing
the packages in the microservice.

## Repository/Service Names
The repository names (usually also the name the service is deployed under) should follow these conventions: 
 - should be preceded by "ok-" or a similar prefix denoting the project 
 -- "ok" : OptiiKeeper 
 - should be appended by the service type string (table below) 
 -- "mgr" or "manager" : service that performs CRUD operations and persists data 
 -- "svc" or "service" : service that performs some sort of processing, orchestration, or calculations 
 -- "proxy" : service that communicates with another remote service (typically 3rd party) and manages access 
 -- "oe" : service that manages a process or algorithm - an orchestration engine 
 - should be, at most, 2 words that describe the service (i.e. "workflow-generator", "users", etc.)

## Supported Go Compiler Versions
Supported Go compiler versions are: v1.11+ with module support enabled

## Editors
Supported editors:
 - JetBrains GoLand IDE
 - VS Code (golangci-lint linter)

## Libraries/Packages
### 3rd Party Libraries
 - http processing (github.com/gorilla/mux) : standard http handler package allowing us to process query and path parameters easily
 - bitbucket.org/optiisolutions/go-common : library containing multiple packages that can be used in your development
 - go-pg (github.com/go-pg/pg) : Database package that we are using for Postgres database
 - jsonapi (github.com/google/jsonapi) : package allowing for managing JSONAPI marshalling and unmarshalling
 - general mock (github.com/golang/mock/gomock) : package that allows us to mock packages for easier and more complete unit testing
 - .env support (github.com/joho/godotenv) : package that provides general .env (environment variable) support

### Library Development Best Practices
#### Recommended Library Code Format
This format allows for IoC and easier unit testing by allowing the library to be mocked. However, there may be instances where implementing this format would be difficult.
    
    type {model} struct {
    }
    
    func {package}Factory() (rv {model}) {
        ...
        return
    }
    
    func (this {model}FunctionOne() {
        ...
    }
    
    func (this *{model}FunctionTwo() {
        ...
    }
    
#### Library Development Process
It is best to start out by developing packages in the project that requires the functionality. If it is decided that the package would be helpful in current or future projects, then it can be moved to the go-common repo. To request this, please contact the project lead, architecture lead, or Wayne Brush.

## REST Endpoint Best Practices
### Handling HTTP Request "Content-Type" & Response "Accept"
{GWB} need to determine if we are going to support JSONAPI

Keep in mind that "application/json" and "application/vnd.api+json" are different data formats and need to be accounted for separately. The UI team has requested that we support the JSON API standard but I find it is easier to debug using plain JSON structures. Since it is easy to marshal and unmarshal based on the "Accept" or "Content-Type" strings, it is possible to support both. However, "application/vnd.api+json" support is required while supporting any other format is optional.  

Since we haven't defined a JSON API structure for errors, returned error codes and messages should either be in "application/text" or "application/json" format.

### Returning HTTP Status Values
HTTP status codes should be implemented as follows:
 - 200 : status ok : request was processed successfully (includes case where no data is returned)
 - 201 : status created : request was processed successfully and new resources were added
 - 204 : no data : request was processed successfully and no data was returned
 - 400 : bad request : request can not be processed due to a client error (generally, syntax error) 
 - 404 : not found : resource or service was not found
 - 415 : bad media type : invalid or unsupported "Accept" or "Content-Type" value
 - 500 : internal server error : server error that prevented request from being completed
 - 501 : not implemented : endpoint functionality has not been implemented

## Unit Tests
Unit tests should be written for each package. The current process recommends coverage of 80% or higher. Note: Test Driven Development (TDD) stipulates that the tests are written first and then the code is developed to make the test pass. To reduce the possibility of recurring errors, it is recommended that TDD be used (if possible) when debugging issues/bugs. Ideally, unit tests will cover edge cases as well as the "happy" path.

# Supported Features
## Auto-versioning Support
This service supports auto-versioning. The build command in Dockerfile will load the short git commit number into the "commit" variable in the main() function. It will also load the current time and date into the "builtAt" variable in the main() function. These variables are loaded into the configuration so they can be accessed by project packages. These values are printed out when the service boots up and are returned as part of the "info" endpoint data.

## Swagger
This service shows the swagger auto-generation capability. By adding swagger comments in the code, it is possible to have "go generate" create a new and current swagger.json file in the "/api" directory. This file is served using the swagger-ui package at the "/api/help/" endpoint. The swagger.json file is updated as part of the build process so the deployed swagger should always be current (although not necessarily correct).

## .env Support
This service supports a ".env" file which can be used to load environment variables when the service starts. It is expected that the .env file is in {root directory}/files. If the environment variable is already defined, the package will not overwrite it's value. However, if the variable is not defined, the package will set it's value to the value in the specified file.

# Recommended Development Approach
It is recommended that new service development starts by following these steps:
 - create a repo in bitbucket following the naming convention
 - load latest go-template-service code
 - copy the files from this project into the new project (develop branch)
 - edit the files for the correct import path to the new project
 - edit the .env file to specify required environment variables
 - verify changes by:
 -- "go generate"
 -- "go build"
 -- "go test ./... -covermode=count"
 -- run locally and test "/api/info", "/api/ping", and "/api/help/" endpoints
 - check in basic functionality to develop branch and merge to master
 - continue service development

#### Makefile
```bash
make (all)
```
> Runs tests, builds binary, runs binary

```bash
make build
```
> Builds binary for PC platform

```bash
make test
```
> Runs tests

```bash
make cover
```
> Generates test coverage report to a _cover.out_ file

```bash
make clean
```
> Cleans the project, removes the binary

```bash
make run
```
> Builds binary, runs binary

```bash
make deps
```
> Installs all the project deps by _go get -u_

```bash
make mockgen
```
> Generate mocks for provided packages and interfaces.   
> Can be changed directly in _Makefile_.


### Integration Testing
Integration tests have been added to this service. The tests are located in an "integrations" subdirectory and the existing
code has been refactored to allow for running those tests. Things to notice are:
* each test has it's own shard (i.e. update role has shard 5)
* the mainline initialization code has been refactored to allow for "blocking" and "non-blocking" operation
* nil is not always equal to nil in golang. found a case where we had to type the nil value to get it to work
* I've inlcuded a script "integration-tests.sh" in the root directory which will run the integration tests from the command line
* All test files must be tagged with either "unit" or "integration" directive lines
* the tests can be debugged by accessing the postgresql DBs using pgAdmin when the containers are running (up)

