# Peer-to-Peer Signaling Server

This is a signaling server that allows clients to establish a WebRTC connection. The server supports multiple clients in a room and multiple rooms.

## Features

- Supports multiple clients in a room.
- Supports multiple rooms.
- Allows two or more clients to establish a WebRTC connection.
- Supports reconnecting to the server after a connection is lost.

## Tech Stack

The application is built with **Golang** and uses the **Gorilla** framework **mux** and **websocket** packages for the server logic.

## Folder Structure

```
/
├── cmd # Main application entry point.
├── config # Configuration files.
├── domain # Business logic.
├── request # Request types.
├── response # Response types.
├── server # Server logic.
  ├── handler # Request handlers.
  ├── websocket # Websocket logic.
├── test # Load testing files.
```

## Demo

Link for the demo: [https://webrtc.xhuliodo.xyz/](https://webrtc.xhuliodo.xyz/)

## Installation

1. Clone the repository.
2. Install the dependencies by running `go mod tidy`.
3. Start the development server by running `go run cmd/main.go`.

## Contact

**Maintainer**: Xhulio Doda

**Email**: xhuliodo@gmail.com
