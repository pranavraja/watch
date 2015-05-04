
Hey don't use this, it doesn't scale.
[fswatch](https://github.com/emcrisostomo/fswatch) is way better.

Watch the current directory, run a command.

    go get github.com/pranavraja/watch

Doesn't re-run the command if already running.

Processes simple .gitignore rules, but nothing too fancy

# Prerequisites

- [Go 1.2](http://golang.org/doc/install)

# Setup

    go get

# Running

    watch <cmd>

e.g.

    watch make

