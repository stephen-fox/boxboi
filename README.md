# boxboi

boxboi is a capture the flag (CTF)-style type challenge that demonstrates
privilege escalation and other fun vulnerabilities. Its security model may
or may not be based on a game console from the early 2000s...

## Setup

boxboi can be built using the `builder` program. This will create a directory
named `build` which can then placed into a CTF environment:

```sh
go run builder/main.go
```

The `boxboi` executable is the entry point for the challenge. It should be
executed as a different user account from the player's (e.g., `root`).

Players can connect to boxboi via TCP using a programs such as netcat:

```sh
nc 127.0.0.1 3249
```

glhf :)
