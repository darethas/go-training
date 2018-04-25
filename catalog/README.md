# catalog

catalog microservice.

keeps track of quantities and pricing information for the items in our catalog.

# dependencies

1. github.com/go-sql-driver/mysql
2. github.com/gorilla/mux

# build

`cd cmd/server`
`go build`

# run

./server

# deploy

`scp server your-remote-server:`

# env variables

`DB_USERNAME`
`DB_PASSWORD`
`DB_HOST`
`DB_PORT`