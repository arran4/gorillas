# Gorillas

The GorillaStack maintained version of everyone's favourite DOS game.

Friday evening gorillas tournaments and beers form a cornerstone of GorillaStack's ethos, culture, ideology [insert further synonyms here].

### Development Roadmap

* Optional wind fluctuations on each throw
* Optional winner's throw first
* Option to save throw and replay 'Greatest Hits'


## Go Ports

A minimal Go implementation of Gorillas is provided under `cmd/gorillia-ebiten` using the Ebiten game library. An alternative terminal based version lives in `cmd/gorillia-tcell` built with the tcell library.

Both ports are simplified remakes, sharing core gameplay logic in the
`github.com/arran4/gorillas` module. The tcell version uses ASCII graphics
reminiscent of the original QBasic game and keeps score between rounds.

