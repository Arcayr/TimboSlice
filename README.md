# TimHortons

TimHortons is an IRC bot generating sentences of arbitrary length based on the Markov stochastic model, with input fed from the channels in which it resides.

TimHortons the bot is utterly unrelated to Tim Hortons the Canadian fast food restaurant. Despite this, it is still time for Tims.

## Requirements

- PostgreSQL 9.x

## Building

`go build`.

## Running

`/path/to/timhortons /path/to/timhortons.json`

TimHortons requires a PostgreSQL user and database in order to store any chains it collects. The particulars of the username, password, host, and database name can be configured in the supplied `timhortons.json` file.

## License

TimHortons is ISC licensed. See [LICENSE](./LICENSE).
