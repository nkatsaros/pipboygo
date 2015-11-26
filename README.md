# pipboygo
An enhanced Pip-Boy second screen for Fallout 4

## Work in Progress
Currently, the tool searches for the first non-busy game it can find and connects to it. When connected, it prints that it has received messages. If the game is busy or an error occurs it will try to connect to the next available game.

## Prerequisites
* [Go](https://golang.org/), probably any version

## Download and Install
* `go get -u github.com/nkatsaros/pipboygo`
* `go install github.com/nkatsaros/pipboygo`

## Usage
`pipboygo`

Quit with Ctrl+C

## Thanks
* @rgbkrk for making his [blog post](https://getcarina.com/blog/fallout-4-service-discovery-and-relay/)
* @NimVek for his excellent [protocol analysis](https://github.com/NimVek/pipboy)