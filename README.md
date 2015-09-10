# go-upstart
Go library to interact with upstart jobs

## Why?
Because you want to be able to control system services (via upstart) from your `go` code.

## What?

Upstart is the init system that several linux distributions have shipped with (most notably, Ubuntu 12.04 and 14.04).
You can read more about it [here](http://upstart.ubuntu.com/cookbook).

If you use one of these distributions in production, you've probably run something like `sudo service ssh restart` to control services.
This package provides a programmatic way of doing that from go code (without invoking the command in a subshell).

## How?
```
import github.com/amoghe/go-upstart

func main() {
  // Somewhere in your program
  _ = upstart.RestartJob("ssh")
}
```

This pkg also ships with a cmdline tool (`go-upstart-ctrl`) that allows you to exercise
this library without writing any go code.

## Gotchas!

This lib attempts to detect whether it is on a system where upstart is the init system in a very naive way -
by trying to communicate with the upstart manager object. If this fails, it determines that upstart is not
available and fails all subsequent operations.

This detection is done exactly once, at pkg initialization. Detection can be retriggered subsequently by using
the appropriate `Detect()` method (which returns a `bool`, indicating whether we could communicate with upstart).

Even if upstart is detected, and we're able to communicate with it, it is likely that your program may have
insufficient privileges to actually control jobs (stop/restart of system services typically needs root access).
In these cases the pkg functions return `error`s, which are quite obtuse. These arise in the underlying `dbus`
lib and this pkg doesn't attempt to differentiate between permission errors and other errors. Patches to improve
this are more than welcome!
