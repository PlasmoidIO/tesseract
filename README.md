# share

A peer-to-peer, consumer-facing filesharing solution built in Go.

## Developers

- `Development Lead: Saif Suleman <saif@visionituk.com>`
- `Development VP: Jordan LaPrise <jlaprise@blakwurm.com>`

## Architecture

### Authentication Server

The authentication server we design here will be used for all of our services. Later down the line, we plan on introducing other products such as card games, P2P messaging platforms, etc.

Designing a stable, powerful authentication server for all of our needs makes it easier to produce production-ready products on the Go, with minimizing the concern for security.

### Central Server

The central server will facilitate handshake communications between two users. Both users will send all handshake data through our central servers, in order to specify any information required before initiating a peer-to-peer connection.

### Peer (Microservice)

The peer microservice will handle all the technical aspects of sharing files. It will communicate with the central server, connect to other peers, and listen for peer-to-peer connections. The peer will use libp2p.io, and is compiled to a single executable file the UI will call.
