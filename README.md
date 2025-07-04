# Gorillas

The GorillaStack maintained version of everyone's favourite DOS game.

Friday evening gorillas tournaments and beers form a cornerstone of GorillaStack's ethos, culture, ideology [insert further synonyms here].

![Screenshot_20250628_100528.png](docs/Screenshot_20250628_100528.png)

### Development Roadmap

* Optional wind fluctuations on each throw via `GORILLAS_WIND_FLUCT` setting
* Optional winner's throw first via `-winnerfirst` flag or `GORILLAS_WINNER_FIRST` setting
* BASIC-style wind each round via `GORILLAS_VARIABLE_WIND` setting
* Save your best shots and replay them later from the new Replays menu


## Go Ports

Vibe code ported

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

### Configuration

Certain options can also be toggled through environment variables. Set
`GORILLAS_FORCE_CGA=true` to restrict both frontends to the classic CGA
palette (black, cyan, magenta, white). This can be handy on limited
terminals or for nostalgia.

### Building and Running

#### Prerequisites

- Go toolchain installed (`go` 1.20 or newer).
- Module downloads require internet access on first build.
- `gorillia-ebiten` needs a graphical desktop environment and a working C toolchain.
- Install Ebitengine dependencies on Linux:
  `apt-get install gcc libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`
- `gorillia-tcell` requires a valid `$TERM` setting and the `infocmp` command.

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

### Controls

Use the arrow keys or a gamepad to adjust angle and power in 0.5-unit steps.
You can also type numbers directly. After a short pause the value is accepted,
or press `,` while typing to switch from angle entry to power entry (and again
to confirm and throw).

### Replays

Your best shots are now stored for later. Select "R - Replays" from the menu to watch any saved throws.

### Running Tests

The core library depends on Ebiten for sound effects, which requires system
libraries that may not be available in all environments. When running tests you
can use stub implementations to avoid these dependencies:

```bash
go test -tags test
```

### Known limitations

- The Ebiten version currently has no computer controlled opponent.
 - The tcell version allows arrow keys or typed numbers for input and requires a UTF-8 capable terminal.
- Sound support may vary across platforms.

