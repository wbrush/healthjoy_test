# healthjoy_test
This service was written as a take home coding test for a position at HealthJoy. 

## Test Status
Unfortunately, I was not able to complete it during the weekend that I alloted for the task. The status is given below. I don't know if they want me 
to complete it or will accept it as is. With my workload at my current position, I don't believe that I will be able to spend any time on it until 
next weekend. Most of my time working on this challenge was spent standing up a development environment on my personal machine so that I could start 
on the test.

## Test Task
GOAL
Create a web service which will fork it's own Github repo to a user's account.
REQUIREMENTS
Feature:
1. A user should be able to trigger the "fork" from a browser (e.g. clicking on a link)
2. Once complete, the forked Github repo should be accessible in the user's Github account
3. Ensure your code is deploy ready
Technical:
1. Python 3.7 or the language of your choice utilizing a framework of your choice
2. Include instructions on how to run the solution locally.
Please add any clarifications needed for design decisions made.

## Design tradeoffs

- In reviewing this issue I realized that since it is a single endpoint, it might be better to implement this using a serverless implementation. However, I have not had much experience with that technology and did not feel that I had the time to research it, implement it, and test it for this test. Therefore, I decided to use a microservice template that I have used in the past to create a microservice that would handle this issue. It has what I consider to be the basic needs of a microservice and allows me to get a base microservice up and running in about an hour.

- In an effort to track versions, this service allows the compiler to insert the git repo commit version and the build date/time in the executable. I usually build in unix so the command is:
go build -ldflags "-X main.commit=`git rev-parse --short HEAD` -X main.builtAt=`date +%FT%T%z`"

This will insert the current values in to the executable and the commit can be used as a version tracker. Unfortunately, this command works only in unix (or unix shell like Ming64) and I haven't found the correct syntax for windows yet.

# Development Notes

## Unit Testing
Some unit tests have been written for the packages in this service. 

## Developer Notes
These are notes that I recommend providing to the next dev who might be working on this.

### Status
The service is mostly functional and runs locally. The following endpoints have been verified as working:
 - /api/info : returns system information that is useful in debugging some issues
 - /api/ping : returns nothing but is useful in determining connectivity
 - /api/help : returns the swagger page which documents the endpoints functionality
 - /api/v1/copy_repo : would do the actual task that the test was to implement of forking a git repo to my repos. This is not fully implemented but is started. The steps are given in the handler function as comments on what needs to be done.

## Required Packages
This service uses go-modules so the modules are documented in the go.mod file in the mainline. However, key modules are listed here for completeness.

### Swagger
Used for compiling the swagger comments into a swagger.json file.

go get -u github.com/go-swagger/go-swagger/cmd/swagger

### Gorilla
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/handlers

### logrus
go get -u github.com/sirupsen/logrus

# Deployment
Need to document (and/or update) our deployment process here

## Process
This would spell out the deployment steps once developed and tested.

## Running locally
To run this locally, you would need to have the go compiler installed. I can add the windows executable to the repo if needed. To run locally:
1. run the command "go get github.com/wbrush/healthjoy_test" from the command line
2a. if on unix, run " go build -ldflags "-X main.commit=`git rev-parse --short HEAD` -X main.builtAt=`date +%FT%T%z`" "
2b. if on windows, run "go build"
3a. if on unix, run "./healthjoy_test" 
3b. if on windows, run "healthjoy_test.exe"
4. you should see startup messages similar to:

time="2021-07-11 19:32:30.811276" level=debug msg="loading environmental variables"
time="2021-07-11 19:32:30.812310" level=warning msg="HOST env variable not found. Using Default!"
time="2021-07-11 19:32:30.812310" level=info msg=------------------------------
time="2021-07-11 19:32:30.812310" level=info msg="Starting healthjoy_test"
time="2021-07-11 19:32:30.812310" level=info msg="Version:; Build Date:"
time="2021-07-11 19:32:30.812310" level=info msg=------------------------------
time="2021-07-11 19:32:30.815068" level=info msg="Starting REST Server on port 8000..."

5. The API docemented above should be available at "localhost:8000"
