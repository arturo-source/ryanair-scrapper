# Get Ryanair flight prices with just 1 yaml

Set a `config.yml` like the [example.yml](blob/main/example.yml). Simply copy it and replace it with desired values. You can delete the Telegram secrets, if you only want the output in the console.

## Build the tool

It is programmed in Go, so you will need a Golang compiler.

```sh
git clone https://github.com/arturo-source/ryanair-scrapper.git
cd ryanair-scrapper
go build
```

After that you can run it `./ryanair-scrapper`.

## Extra usage

You can set also the values from the command line, see help `./ryanair-scrapper --help`:

```sh
Usage of ./ryanair-scrapper:
  -config-file string
        Set the yaml with the configuration (default "config.yml")
  -dates string
        Comma-separated dates
  -destinations string
        Comma-separated destinations
  -origins string
        Comma-separated origins
  -telegram-chat-id string
        Telegram chat id to send info
  -telegram-token string
        Telegram bot token
```

Example of usage (origins are set in the config file): `./ryanair-scrapper -destinations BER,BOD -dates 2024-12-10,2024-12-11`.

Values can be set also from the env variables. `TELEGRAM_TOKEN`, `TELEGRAM_CHAT_ID`, `DATES`, `ORIGINS`, `DESTINATIONS` are the env variables.

The order of preference that overrides the values is (from most important to less):

1. Command line
2. Environment variable
3. Config file
