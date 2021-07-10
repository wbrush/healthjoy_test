#/bin/sh

# this script will run the integration tests. Note that we are using the same methodology that the .NET guys are using. We
# are pulling a docker compose file and it contains the DB and pubsub functionality so that we can actually hit the DB and
# pubsub.
docker-compose -f ./docker-compose.integration-test.yml up -d --remove-orphans

# perform integration test suite
go test -tags=integration -v ./...

# clean up
docker-compose -f docker-compose.integration-test.yml down
