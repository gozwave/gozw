## Summary

1. **Transport layer** - handles raw communication with the Z-Wave controller
1. **Frame layer** - handles frame encoding/decoding
  1. **Frame parser** - parses frames by implementing the state machine defined in INS12350 section 6.6.1
1. **Session Layer** - handles Z-Wave API function calls (as described in INS12308). Handles the Serial API request/response flow, as well as the association of callbacks to their requests.
  1. **Security Layer** - uses the session layer to handle functions related to the security command class.
1. **Application Layer** - interaction point for userland code.

```
Application Layer
       |
       V
 Session Layer   < ---- >   Security Layer
       |
       V
  Frame Layer    < ---- >   Frame parser (plus state machine)
       |
       V
 Transport Layer
```

## Layers

### Transport Layer

Handles reads/writes with the serial port. In the case of future implementations that do not use a USB-UART driver, it should be possible to substitute an implementation for some other I/O device (UART/SPI/etc).

#### Responsibilities

 - Write Serial API frames (represented as a byte slice) to an I/O device
 - Read Serial API frames from the I/O device and emit (one byte at a time) each byte on a channel that can be accepted by the frame layer

### Frame Layer

Parses incoming frames using the Frame Parser and determines how to proceed based on the parser result. Implements (through code, rather than an actual FSM, unlike the Frame Parser) the Host Media Access Control state machine as described in INS12350 section 6.6.2.

Handles parse/receive timeouts, as well as transmission of ACKs/NAKs

#### Responsibilities

 - Encode frames to raw bytes and write them to the transport layer
 - Read bytes from the transport layer and send them to the Frame Parser
 - Perform basic locking to prevent transmitting frames while in the process of receiving frames
   - **Note:** because this layer does not have any knowledge of Z-Wave API functions, it will not perform any locking with regard to the REQ/RES flow (INS12350 section 6.5.2)
 - ACK valid frames (based on the frame checksum)
 - NAK invalid frames (based on the frame checksum)
 - **TODO:** Handle ACK timeouts
 - **TODO:** Handle transmit/receive collisions
 - **TODO:** Handle CAN frames
 - **TODO:** Handle retransmission and backoff

### Session Layer

Handles incoming frames and exposes an API to make Z-Wave function calls (such as `AddNodeToNetwork` and `SendData`). Implements the Host Request/Response Session state machine as described in INS12350 section 6.6.3.

This layer uses a mutex to prevent concurrent requests and the interruption of blocking processes (such as network management)

Whenever possible, methods exposed by this layer should block until their corresponding Z-Wave operation has completed (e.g. AddNodeToNetwork, which has a complex workflow consisting of multiple callback functions, should block until the entire process has concluded, and return the newly added node or an error).

#### Responsibilities
 - Implementation of Serial API functions
 - Locking to prevent request concurrency and interruption of blocking processes
 - Routing of responses and callbacks to their registered handlers
 - Routing of unsolicited commands (usually from ApplicationControllerUpdate) to the application layer
 - Handling of security command classes (via the Security Layer)

### Security Layer

For background, see SDS10865 (yep, the whole document; it's not that long).

Intercepts appropriate commands in the security command class in order to make Z-Wave security functions transparent to the application (the application layer must decide whether to call SendData or SendDataSecure in the session layer, but it does not need to handle payload encryption/verification, nonce management, or sequencing).

This layer is also responsible for encrypting/decrypting message payloads, fetching nonces from other nodes, generating nonces for other nodes to use, and reassembling command classes that have been sequenced into multiple frames due to encryption overhead.

#### Responsibilities
 - Generating internal nonces in response to nonce challenges from other nodes
 - Fetching nonces from other secure nodes
 - Mangaging nonce timers and usage (they are short-lived and can only be used once)
 - Encrypting outgoing message payloads (and creating HMACs)
 - Sequencing (splitting) of outgoing message payloads when necessary
 - Decrypting incoming message paylods (and verifying HMACs)
 - Reassembly of sequenced incoming messages

### Application Layer

The application layer abstracts the Z-Wave protocol so that user implementations do not need an in-depth understanding of Serial API functions or command classes. It manages the network at a high level by keeping track of nodes, acting as a proxy to the session layer, and receiving data frames (typically command classes whether GETs or REPORTs) from the session layer.

## Resources

1. INS12308 - Z-Wave 500 Series Application Programming Guide (v6.51.06)
1. INS12350 - Serial API Host Application Programming Guide
1. SDS10865 - Z-Wave Application Security Layer
