FROM golang AS gobuilder
WORKDIR /src/github.com/florianloch/spotistate
COPY . .
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o spotistate .

FROM node AS webuibuilder
WORKDIR /build
COPY . .
# Don't use dependencies that have been build on the host system...
RUN rm -rf node_modules
RUN npm install
RUN npm install -g grunt
RUN grunt

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /src/github.com/florianloch/spotistate/spotistate .
COPY --from=webuibuilder /build/webui ./webui
CMD ["./spotistate"]