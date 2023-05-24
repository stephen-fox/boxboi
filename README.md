# boxboi

boxboi is a capture the flag (CTF)-style challenge that demonstrates
privilege escalation and other fun vulnerabilities. Its security model
may or may not be based on a game console from the early 2000s...

## Playing

The boxboi executable contains a vulnerability that allows players to
execute code of their choosing. Players should be provided a user
account (with a shell) on the same computer that boxboi runs on.

The goal of the challenge is to get code execution as the user that
owns the boxboi process.

Players can connect to boxboi on loopback at TCP port 3249:

```sh
nc 127.0.0.1 3249
```

glhf :)

## Setup

boxboi can be built using the `builder` program. This will create a directory
named `build` which can then placed into a CTF environment:

```sh
go run builder/main.go
```

The `boxboi` executable is the entry point for the challenge. It should be
executed as a different user account from the player's (e.g., `root`).
