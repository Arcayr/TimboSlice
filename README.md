# TimHortons

TimHortons is an IRC bot generating sentences of arbitrary length based on the Markov stochastic model, with input fed from the channels in which it resides.

TimHortons the bot is utterly unrelated to Tim Hortons the Canadian fast food restaurant. Despite this, it is still time for Tims.

## Requirements

- PostgreSQL 9.x

## Building

`go build`.

## Running

There must exist a valid `timhortons.json` file from the directory which you run TimHortons.

Run: `/path/to/timhortons`

TimHortons requires a PostgreSQL user and database in order to store any chains it collects. The particulars of the username, password, host, and database name can be configured in the supplied `timhortons.json` file.

## Running in Docker

In order to build and run TimHortons in a docker container, run the following:

Build:  `docker build -t tim .`

Assuming you aren't using compose, and your postgresql server is running locally

Run local build:

`docker run --rm --net="host" -v $PWD/timhortons.json:/app/timhortons.json tim` 

Run from dockerhub:
 
`docker run --rm --net="host" -v $PWD/timhortons.json:/app/timhortons.json ceruleis/timhortons`

## Running with docker-compose

Coming soon.

## License

TimHortons is ISC licensed. See [LICENSE](./LICENSE).
