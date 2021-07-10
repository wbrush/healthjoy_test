GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get -u
LD_FLAGS="-X main.commit=`git rev-parse --short HEAD` -X main.builtAt=`date +%FT%T%z`"
COVERPROFILE=cover.out
BINARY_NAME="$(notdir $(CURDIR))"

all: unit run
build:
	$(GOBUILD) -ldflags $(LD_FLAGS)
unit:
	$(GOTEST) -tags=unit -v ./...
integration:
	$(GOTEST) -tags=integration -v ./...
cover:
	$(GOTEST) -tags=unit -coverprofile $(COVERPROFILE) ./...
	rm -f $(COVERPROFILE)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
deps:
	$(GOGET)

mockgen:
	#TODO add all packages and interfaces that need to be mocked
	mockgen -destination ./dao/mock_dao/main.go bitbucket.org/optiisolutions/go-template-service/dao DataAccessObject