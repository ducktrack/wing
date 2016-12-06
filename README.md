# Wing

## Running

`go run main.go`

## Running tests

`go test $(go list ./... | grep -v /vendor/)`

## Development

__Wing__ uses `glide` for dependency management, this setup the project do the following:

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
