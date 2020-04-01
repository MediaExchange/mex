# Build the UI
FROM node:alpine AS node
WORKDIR /app
COPY ./ui /app
RUN npm update && ./node_modules/.bin/ng build --prod=true

# Build mex
FROM golang:alpine as golang
WORKDIR /go/src/mex
COPY . /go/src/mex
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/mex

# Create the container
FROM scratch
COPY --from=node /app/dist/mex /ui/dist/mex
COPY --from=golang /go/bin/mex /
COPY --from=golang /go/src/mex/mex_config.yaml /
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 9000
CMD ["/mex"]
