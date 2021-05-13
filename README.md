# mountmond

## About

`mountmond` is a simple service that monitors your mounts. If any
mount is missing, `mountmond` runs a command that you specify to
restore the mount. This allows you to create custom alerting behavior,
etc. 

## Building

Run `make`.

## Configuring

You can edit `/etc/mountmond.yaml`, or use the `-config-path` flag and
specify the path to your configuration. See `examples/` for
configuration examples.
