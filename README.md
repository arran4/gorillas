# Gorillas

The GorillaStack maintained version of everyone's favourite DOS game.

Friday evening gorillas tournaments and beers form a cornerstone of GorillaStack's ethos, culture, ideology [insert further synonyms here].

### Development Roadmap

* Optional wind fluctuations on each throw
* Optional winner's throw first via `-winnerfirst` flag or `GORILLAS_WINNER_FIRST` setting
* BASIC-style wind each round via `GORILLAS_VARIABLE_WIND` setting
* Option to save throw and replay 'Greatest Hits'


## Go Ports

A minimal Go implementation of Gorillas is provided under `cmd/gorillia-ebiten` using the Ebiten game library. An alternative terminal based version lives in `cmd/gorillia-tcell` built with the tcell library.

Both ports are simplified remakes, sharing core gameplay logic in the
`github.com/arran4/gorillas` module. The tcell version uses ASCII graphics
reminiscent of the original QBasic game and keeps score between rounds.
When launching either port you can override some gameplay settings using
command-line flags:

```
  -wind       starting wind speed
  -gravity    gravitational constant
  -rounds     number of rounds to play
  -buildings  how many buildings appear in the skyline
  -winnerfirst winner of a round starts next
```

### Building and Running

#### Prerequisites

- Go toolchain installed (`go` 1.20 or newer).
- Module downloads require internet access on first build.
- `gorillia-ebiten` needs a graphical desktop environment.

#### Build commands

```bash
# Build the Ebiten GUI version
go build -o gorillia-ebiten ./cmd/gorillia-ebiten

# Build the terminal version
go build -o gorillia-tcell ./cmd/gorillia-tcell
```

#### Example usage

```bash
# Start the graphical port with 10 rounds
./gorillia-ebiten -rounds 10

# Play in the terminal with a computer opponent
./gorillia-tcell -ai
```

### Running Tests

The core library depends on Ebiten for sound effects, which requires system
libraries that may not be available in all environments. When running tests you
can use stub implementations to avoid these dependencies:

```bash
go test -tags test ./...
```

### Known limitations

- The Ebiten version currently has no computer controlled opponent.
- The tcell version uses arrow keys and Enter for input and requires a UTF-8 capable terminal.
- Sound support may vary across platforms.

