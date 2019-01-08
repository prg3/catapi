# catapi

## Build

Execute `make` and it will build a local native copy of the application as well as a docker container

## Execute

`make docker-run`or `./catapi`

## Environment

Three environment variables are avaliable:
 - APIKEY - Required, this is an API Key from https://thecatapi.com
 - REDIS - Optional - String reprenseing a host:port for a Redis instance for persistence
 - DEBUG - If anything is set in this variable, additional debugging output will be printed

 Create a file called .env with these set for `make docker-run` to execute properly