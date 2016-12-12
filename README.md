# Wing

## Running

`go run main.go`

## Running tests

`make test`

## Development

__Wing__ uses `glide` for dependency management, to setup the project do the following:

```sh
brew install glide
```

__Note__: Read here (https://github.com/Masterminds/glide) for a different OS

```sh
glide install
```

To add new dependecies do:

```sh
glide get <dependency>
```