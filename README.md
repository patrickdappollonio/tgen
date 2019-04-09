# `tgen`

[![Build Status](https://travis-ci.org/patrickdappollonio/tgen.svg?branch=master)](https://travis-ci.org/patrickdappollonio/tgen)

`tgen` is a simple CLI application that allows you to write a template file and then use the power of Go Templates to generate an output (which is) outputted to `stdout`. Besides the Go Template engine itself, `tgen` contains three extra utility functions to assist you when writing templates. See below for a description of each.

You can also use the `--help` (or `-h`) to see the available set of options. The only flag required is the file to process, and everything else is optional.

```
tgen is a template generator with the power of Go Templates

Usage:
  tgen [flags]

Flags:
  -d, --delimiter string     delimiter (default "{{}}")
  -e, --environment string   an optional environment file to use (key=value formatted) to perform replacements
  -f, --file string          the template file to process (required)
  -h, --help                 help for tgen
      --version              version for tgen
```

### Available functions

* `{{ env "NAME" }}`: `env` allows you to fetch an environment variable. It is case insensitive, so either `NAME` or `name` will work.
* `{{ envdefault "NAME" "my name" }}`: `envdefault` allows you to fetch the value of an environment variable, and if it wasn't found, then set the value to `"my name"`. It's also case insensitive.
* `{{ raw "ABC123" }}`: `raw` will print a raw output. While this may seem cumbersome, some functions may output a buffer that needs to be rendered here, so this function can take it and spit it out.
* `{{ sprintf "format" }}`: `sprintf` works the same way as Go's `fmt.Sprintf()`, it allows you to pass a string format, for example, `"Hello, %s!"` and then zero or more arguments of any type, then uses `fmt.Sprintf()` under the hood to print it out. For example, calling `{{ sprintf "Hello, %s!" "Peter" }}` will print `Hello, Peter!`

### Environment file

`tgen` supports an optional environment variable collection in a file but it's a pretty basic implementation of a simple key/value pair. The environment file works by finding lines that aren't empty or preceded by a pound `#` -- since they're treated as comments -- and then tries to find at least one equal (`=`) sign. If it can find at least one, all values on the left side of the equal sign become the key -- which is also uppercased so it's compatible with the `env` function defined above -- and the contents on the right side become the value. If the same line has more than one equal, only the first one is honored and all remaining ones become part of the value.

As an important note, environment variables found in the environment have preference over the environment file. That way, the environment file can define `A=1` but then the application can be run with `A=2 tgen [flags]` so it overrides `A` to the value of `2`.

## Example

Consider the following template, named `template.txt`:

```
The dog licked the {{ env "element" }} and everyone laughed.
```

And the following environment file, named `contents.env`:

```
element=Oil
```

After being passed to `tgen` by executing `tgen -e contents.env -f template.txt`, the output becomes:

```
The dog licked the Oil and everyone laughed.
```
