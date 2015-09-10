# go-upstart

Go library to interact with upstart jobs

## Why?

Because you want to be able to control system services (via upstart) from your `go` code.

## What?

Upstart is the init system that several linux distributions have shipped with (most notably, Ubuntu 12.04 and 14.04).
You can read more about it [here](http://upstart.ubuntu.com/cookbook).

If you use one of these distributions in production, you've probably run something like `sudo service ssh restart` to
control services. This package provides a programmatic way of doing that from go code (without invoking the command in
a subshell - because thats feels dirty).

## How?

```
import github.com/amoghe/go-upstart

func main() {
  // Somewhere in your program
  _ = upstart.RestartJob("ssh")
}
```

This pkg also ships with a cmdline tool (`go-upstart-tool` is in a subdir - so you won't get it on `go install`, you'll
need to build it by running `go build` in the appropriate subdir under your`$GOPATH`).

The tool is useful for testing that the pkg actually works when running on a system that has upstart because it allows
you to exercise the functions in this library without writing any go code.

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

Also, this package lacks tests because tests would need superuser privileges to actually make requests from the upstart
daemon. Any tests written without these privileges would just be testing whether we can exchange messages over dbus,
and thats not very interesting.

## Contributing

If you feel like this package is missing something or if you find bugs, file an issue (and feel free to throw in a PR)!

## License

This pkg is available under the simplified BSD license. See the LICENSE.txt file.
