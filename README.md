# ProtonHub

### What is it?

This is a web client for umu-launcher.

### How to use it?

For now you need to build it from source.
Clone this repository somewhere on your disk. `cd` into `protonhub`.
Make sure you have `go` installed on your system.

Run `go build -v` for building process. And then run `go install`

This will install `protonhub` binary to `$HOME/go/bin/` so you should have added this path to your `$PATH` env variable.

After installation you can open your localhost on port 8080 in your web browser.
From there you can create, edit or delete your games. Everything is running via `Proton`, via `umu-launcher`.
Make sure you have installed `umu-launcher` and have `umu-run` in your $PATH.
