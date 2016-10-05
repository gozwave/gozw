# gozw

Golang Z/IP Portal

## Getting Started

1. You need a Z-Wave USB bridge controller with Serial API >= 4.05
1. With the controller plugged in, if you do not see a device in `/dev` named
   `usbserial*` or `usbmodem*` (or similar), you need to install the FTDI VCP
   driver for your platform
1. OPTIONAL: Install [`interceptty`](http://www.suspectclass.com/sgifford/interceptty/),
   which is an extremely useful serial port / tty proxy (it will allow you to see
   the raw data being transmitted to/from the controller).
  - Installation is as simple as `./configure && make && sudo cp interceptty /usr/local/bin`
  - To run: `interceptty -s 'ispeed 115200 ospeed 115200' /dev/<serialdevice> /tmp/<serialdevice>`
  - Be sure to update your config.yaml to point to `/tmp/<serialdevice>`
1. `go get -u gitlab.com/helioslabs/zwgen`
1. `make install-deps` in the zwgen folder root
1. `make install` in the zwgen folder root
1. `make build` in the zwgen folder
1. `go get ./...` in the gozw folder root
1. `go generate ./...` in the gozw folder root
1. `go run cmd/portald/main.go`
1. `go run cmd/gatewayd/main.go`

## Summary
1. **Transport layer** - handles raw communication with the Z-Wave controller via serial port
1. **Frame layer** - handles frame encoding/decoding to/from binary
  1. **Frame parser** - parses frames by implementing the finite state machine defined in INS12350 section 6.6.1
1. **Session Layer** - handles frame sequencing and host media access control
1. **Serial API Layer** - defines Z-Wave API functions (as described in INS12308) and handles requests, responses, and callbacks
1. **Application Layer** - tracks network data, node information, security information, etc
  1. **Security Layer** - provides utilities for encrypting/decrypting secure messages

```
Application Layer   < ---- >   Security layer
       |
       V
 Serial API Layer
       |
       V
  Session Layer
       |
       V
   Frame Layer      < ---- >   Frame parser
       |
       V
 Transport Layer
```

## Layers

### Transport Layer
Handles reads/writes with the serial port. In the case of future implementations that do not use a USB-UART driver, it should be possible to substitute an implementation for some other I/O device (UART/SPI/etc).

It satisfies several interfaces from the `io` package. In the future, it should be possible to dump a long-running Z-Wave interaction (output from interceptty, for example) to a file, then use that file as input to reproduce crashes.

#### Responsibilities
 - Read/write bytes to/from an I/O source (serial port, file, network, etc.)

### Frame Layer
Parses incoming frames using the frame parser and determines how to proceed based on the parser result. Handles parse/receive timeouts, as well as transmission of ACKs/NAKs.

#### Responsibilities
 - Encode frames to raw bytes and write them to the transport layer
 - Read bytes from the transport layer and send them to the frame parser
 - Perform basic locking to prevent transmitting frames while in the process of receiving frames
   - **Note:** because this layer does not have any knowledge of Z-Wave API functions, it will not perform any locking with regard to the REQ/RES flow
 - ACK valid frames (based on the frame checksum)
 - NAK invalid frames (based on the frame checksum)
 - **TODO:** Handle ACK timeouts
 - **TODO:** Handle transmit/receive collisions
 - **TODO:** Handle CAN frames
 - **TODO:** Handle retransmission and backoff

### Session Layer
Facilitates the request/response flow by queueing requests when awaiting responses and callbacks. Implements the Host Request/Response Session state machine as described in INS12350 section 6.6.3.

#### Responsibilities
 - Locking to prevent request concurrency
 - Routing and matching of responses and callbacks to the appropriate handlers
 - Routing of unsolicited commands (usually from ApplicationControllerUpdate) to the application layer

### Serial API Layer
Exposes an API to make Z-Wave function calls (such as `AddNodeToNetwork` and `SendData`).

Whenever possible, methods exposed by this layer should block until their corresponding Z-Wave operation has completed (e.g. AddNodeToNetwork, which has a complex workflow consisting of multiple callback functions, should block until the entire process has concluded, and return the newly added node or an error).

#### Responsibilities
 - Implementation of Serial API functions

### Application Layer
The application layer abstracts the Z-Wave protocol so that user implementations do not need an in-depth understanding of Serial API functions or command classes. It manages the network at a high level by keeping track of nodes, acting as a proxy to the session layer, and receiving data frames (typically command classes whether GETs or REPORTs) from the session layer.

#### Responsibilities
 - Network management
 - Node information tracking
 - Handling of security command classes (via the Security Layer)

### Security Layer
For background, see SDS10865 (yep, the whole document; it's not that long).

Provides utilities for encrypting/decrypting messages, storing and timing out nonces, and security sequence counters.

#### Responsibilities
 - Generating internal nonces in response to nonce challenges from other nodes
 - Fetching nonces from other secure nodes
 - Mangaging nonce timers and usage (they are short-lived and can only be used once)
 - Encrypting outgoing message payloads (and creating HMACs)
 - Sequencing (splitting) of outgoing message payloads when necessary
 - Decrypting incoming message paylods (and verifying HMACs)
 - Reassembly of sequenced incoming messages

## Resources

1. INS12308 - Z-Wave 500 Series Application Programming Guide (v6.51.06)
1. INS12350 - Serial API Host Application Programming Guide
1. SDS10865 - Z-Wave Application Security Layer
1. http://pepper1.net/zwavedb/ Random Z-Wave device library (lots of technical info)
