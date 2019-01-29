
# Instructor Notes for Solving Problems with NATS

These are the notes for instructor setup for the _Solving Problems with NATS_ course. These are included
in the solutions folder in case students want to repeat the setup.

## Running the main cluster

1. Open 2 terminals.
2. In the first terminal, run `gnatsd -p 4222 -cluster nats://localhost:5222 -routes nats://localhost:5223`.
3. In the second terminal, run `gnatsd -p 4223 -cluster nats://localhost:5223 -routes nats://localhost:5222`.

## Run the local TLS Server

1. Run `gnatsd -c resources/tls/tlsverify.conf`, port is 4443

## Run the streaming server

1. Run `go run nats-streaming-server.go -p 4223 -cid nats_course`

## Exercise Comments

* In ex01 the go client will not get discovered servers on an initial connect, you can make that happen by closing 1 server and starting it after people run their solution.
