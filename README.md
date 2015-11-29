# pipboygo
An enhanced Pip-Boy second screen for Fallout 4

## Work in Progress
The tool waits for the first websocket connection it receives (http://localhost:8000) and streams local map and game info to it. In the browser that visited the page, you will get a live updating local map and a dump of the game state. The color of the map updates when you change your pipboy color in the game. Refreshing is currently broken, so if you want to change something in the javascript you have to stop the pipboygo application and restart it.

## Prerequisites
* [Go](https://golang.org/), probably any version
* node

## Download and Install
* `go get -u github.com/nkatsaros/pipboygo`
* `go install github.com/nkatsaros/pipboygo`
* `npm install`
* `npm run build`

## Usage
`pipboygo -public path/to/build`

Go to http://localhost:8000

Quit with Ctrl+C

## Things to do
* This repo should probably be split into the frontend application and the client meant for other Go apps to use.
* Track state of connected clients in the go application
* Build assets into the Go binary in production mode
* ...

## Thanks
* @rgbkrk for making his [blog post](https://getcarina.com/blog/fallout-4-service-discovery-and-relay/)
* @NimVek for his excellent [protocol analysis](https://github.com/NimVek/pipboy)