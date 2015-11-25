# pipboygo
An enhanced Pip-Boy second screen for Fallout 4

## Work in Progress
Currently, the tool searches for the first non-busy game it can find and connects to it. When connected it dumps messages that it doesn't have handlers for. If it encounters a message it doesn't understand yet, it'll crash. One second after connecting it send a local map request and the response is saved to disk as image.png. If the connection is lost it will go back to looking for a non-busy game and start over again.

## Prerequisites
* [Go](https://golang.org/), probably any version

## Download and Install
* `go get -u github.com/nkatsaros/pipboygo`
* `go install github.com/nkatsaros/pipboygo`

## Usage
`pipboygo`

Quit with Ctrl+C

## Thanks
* [rgbkrk](https://github.com/rgbkrk) for making his [blog post](https://getcarina.com/blog/fallout-4-service-discovery-and-relay/)
* [NimVek](https://github.com/NimVek/pipboy) for his excellent protocol analysis