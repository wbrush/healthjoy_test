# Args defined here will apply to all build stages *IF* you also reference
# the name inside that step.
ARG APPNAME="appname"

# BUILD CONTAINER
FROM golang:1.16.2 AS builder
# NOTE:  Exposing credentials here is bad, even if arguments are only used in
# build stage, as they are recorded into layer and exposed via Docker history
# command.  Still researching best way to manage secrets used via Dockerfile.
ARG BB_UN
ARG BB_PW
ARG APPNAME
# Need to supply creds to pull our private repos
# (This can be combined into RUN step later *or* put into a pre-build script.)
RUN git config --global url."https://${BB_UN}:${BB_PW}@bitbucket.org/".insteadOf "https://bitbucket.org/"
# Copy the local package files to the container's workspace.
WORKDIR /go/src/${APPNAME}
COPY . .
# Need to add swagger.exe to the go/bin directory
#RUN go get github.com/go-swagger/go-swagger/cmd/swagger

# Set build env vars
ENV GO111MODULE=on
# One problem of using make file is that it will *always* download dependencies again
# If we re-use build container, it may be possible to keep dependencies to save
# some size, bandwidth, and compile time. Food for thought.
# Another thing to think about is building a container this way requires static
# linking so the build command is different than what we might normally do in a
# local dev environment.
# I set LD_FLAGS as env var here first just to reduce complexity of nested quotes
# and provide a little more legibility in the mess of chained commands.
#RUN make
RUN LD_FLAGS="-X main.commit=`git rev-parse --short HEAD` -X main.builtAt=`date +%FT%T%z`" \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "${LD_FLAGS}" -o ${APPNAME} .

# prepare designated set of swagger-ui static files in the build container to be copied over to runtime later
# only really need tatic files from it, so need to keep go modules versioning happy
ENV REGEX_VAR="^/go/pkg/mod/${SWAGGER_MODULE_PATH}\(/.*\)?${SWAGGER_MODULE_NAME}\(.*\)?swagger-ui-static"
RUN mkdir /swagger-pkg && find /go/pkg/mod -mindepth 1 -type d -regex ${REGEX_VAR} -exec cp -r --parents {} /swagger-pkg \;

# APP container
FROM alpine:latest AS app-container
ARG SERVICE_PORT=8080
ARG APPNAME
ENV APPBINARY="$APPNAME"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy all essential files from build container.
# Right now we need multiple copy steps because app looks for files at different
# locations.
COPY --from=builder /go/src/${APPNAME}/${APPNAME} .
COPY --from=builder /go/src/${APPNAME}/files/.env ./files/.env
COPY --from=builder /go/src/${APPNAME}/docs/swagger.json ./docs/swagger.json
COPY --from=builder /go/src/${APPNAME}/dao/postgres/base_migrations/ ./dao/postgres/base_migrations/
COPY --from=builder /go/src/${APPNAME}/dao/postgres/shard_migrations/ ./dao/postgres/shard_migrations/
COPY --from=builder /swagger-pkg /

# Expose port for service, (which should be defined in the ./files/.env file).
EXPOSE ${SERVICE_PORT}
# Explicitly calling shell to expand app binary name
CMD ["sh", "-c", "./${APPBINARY}"]

#
# Build service
# docker-build.sh
# docker build -t <service-name> .
#
# Run service:
# docker run -it -p 8091:8091 <service-name>
# If you need to override env var, for example, DB_PORT:
# docker run -e "DB_PORT=5434" -it -p 8091:8091 <service-name>
# if on mac, a workaround is needed for networking issues:
# docker run -e "DB_HOST=docker.for.mac.localhost" -it -p 8091:8091 <service-name>
#
# Debug service
#docker run -it -p 8091:8091 <service-name> sh
