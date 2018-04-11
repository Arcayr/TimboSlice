# TimboSlice

TimboSlice is an IRC bot generating sentences of arbitrary length based on the Markov stochastic model, with input fed from the channels in which it resides.

## Requirements

Nil. TimboSlice uses a SQLite backend as it is only one reader/writer at a time.

## Building

`go build`.

## Running

There must exist a valid `timboslice.yml`, `timboslice.json`, or `timboslice.toml` file from the directory which you run the application.

## License

TimHortons is BSD licensed. See [LICENSE](./LICENSE).
