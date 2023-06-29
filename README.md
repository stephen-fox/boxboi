# boxboi

boxboi is a capture the flag challenge that demonstrates privilege escalation
and other fun vulnerabilities. Its security model may or may not be based on
a game console from the early 2000s...

## Playing

Before playing, please watch Michael Steil's ["Deconstructing the Xbox
security system"](https://www.youtube.com/watch?v=9NqLljaHc80) presentation.
This CTF is heavily based on one of the vulnerabilities discussed therein.

The [boxboi executable](boxboi) contains a vulnerability that allows players
to execute code of their choosing. To find the bug, review the source code of
the application and the other Go applications contained in this repository.
Pay close attention to terminology used in the source code. There are direct
references to Michael's presentation which will guide you along.

Players can connect to boxboi by first SSH'ing to a virtual machine.
From there, user's can reach boxboi's UI via loopback at TCP port 3249:

```sh
nc 127.0.0.1 3249
```

The goal of the challenge is to get code execution as the user that owns
the boxboi process and retrieve the flag file in the process owner's
home directory.

glhf :)

## Setup

boxboi can be built using the `builder` program. This will create a directory
named `build` which can then placed into a CTF environment:

```sh
go run builder/main.go
```

The `boxboi` executable is the entry point for the challenge. It should be
executed as a different user account from the player's (e.g., `root`).

Players should be provided a user account (with a shell) on the same computer
that boxboi runs on.
