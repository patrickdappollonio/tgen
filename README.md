# `tgen`

[![Tests passing](https://img.shields.io/github/workflow/status/patrickdappollonio/tgen/Continuous%20Integration/master?logo=github&style=flat-square)](https://github.com/patrickdappollonio/tgen/actions)
[![Downloads](https://img.shields.io/github/downloads/patrickdappollonio/tgen/total?color=blue&logo=github&style=flat-square)](https://github.com/patrickdappollonio/tgen/releases)


`tgen` is a simple CLI application that allows you to write a template file and then use the power of Go Templates to generate an output (which is) outputted to `stdout`. Besides the Go Template engine itself, `tgen` contains a few extra utility functions to assist you when writing templates. See below for a description of each.

You can also use the `--help` (or `-h`) to see the available set of options. The only flag required is the file to process, and everything else is optional.

```
tgen is a template generator with the power of Go Templates

Usage:
  tgen [flags]

Flags:
  -d, --delimiter string     delimiter (default "{{}}")
  -e, --environment string   an optional environment file to use (key=value formatted) to perform replacements
  -x, --execute string       a raw template to execute directly, without providing --file
  -f, --file string          the template file to process (required)
  -h, --help                 help for tgen
  -s, --strict               enables strict mode: if an environment variable in the file is defined but not set, it'll fail
      --version              version for tgen
```

### Environment file

`tgen` supports an optional environment variable collection in a file but it's a pretty basic implementation of a simple key/value pair. The environment file works by finding lines that aren't empty or preceded by a pound `#` -- since they're treated as comments -- and then tries to find at least one equal (`=`) sign. If it can find at least one, all values on the left side of the equal sign become the key -- which is also uppercased so it's compatible with the `env` function defined above -- and the contents on the right side become the value. If the same line has more than one equal, only the first one is honored and all remaining ones become part of the value.

As an important note, environment variables found in the environment have preference over the environment file. That way, the environment file can define `A=1` but then the application can be run with `A=2 tgen [flags]` so it overrides `A` to the value of `2`.

### Example

Consider the following template, named `template.txt`:

```go
The dog licked the {{ env "element" }} and everyone laughed.
```

And the following environment file, named `contents.env`:

```bash
element=Oil
```

After being passed to `tgen` by executing `tgen -e contents.env -f template.txt`, the output becomes:

```bash
$ tgen -e contents.env -f template.txt
The dog licked the Oil and everyone laughed.
```

Using the inline mode to execute a template, you can also call `tgen -x '{{ env "element" }}' -e contents.env` (note the use of single-quotes since in Go, strings are always double-quoted) which will yield the same result:

```bash
$ tgen -x '{{ env "element" }}' -e contents.env
The dog licked the Oil and everyone laughed.
```

Do note as well that using single quotes for the template allows you to prevent any bash special parsing logic that your terminal might have.

### Template Generation _a la Helm_

`tgen` can be used to generate templates, in a very similar way as `helm` can be used. However, do note that `tgen`'s intention is not to replace `helm` since it can't handle application lifecycle the way `helm` does, however, it can do a great job generating resources with very similar code.

Consider the following example of creating a Kubernetes secret for a `tls.crt` file -- in real environments, you'll also need the key, but for the sake of this example, it has been omitted.

Checking the files in the folder:

```bash
tree .
```

```text
.
├── secret.yaml
└── tls.crt

0 directories, 2 files
```

We have a `secret.yaml` which includes `tgen` templating notation:

```bash
cat secret.yaml
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secret-tls
type: kubernetes.io/tls
data:
  tls.crt: | {{ readfile "tls.crt" | b64enc | nindent 4 }}
```

The last line includes the following logic:

* Reads the `tls.crt` file from the current directory where `tgen` is run
* Takes the contents of the file and converts it to `base64` -- required by Kubernetes secrets
* Then indents with 4 spaces, starting with a new line

To generate the output, we can now run `tgen`:

```bash
tgen -f secret.yaml
```

And the output looks like this:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secret-tls
type: kubernetes.io/tls
data:
  tls.crt: |
    TmV2ZXIgZ29ubmEgZ2l2ZSB5b3UgdXAsIG5ldmVyIGdvbm5hIGxldCB5b3UgZG93bgpOZXZlciBnb25uYSBydW4gYXJvdW5kIGFuZCBkZXNlcnQgeW91Ck5ldmVyIGdvbm5hIG1ha2UgeW91IGNyeSwgbmV2ZXIgZ29ubmEgc2F5IGdvb2RieWUKTmV2ZXIgZ29ubmEgdGVsbCBhIGxpZSBhbmQgaHVydCB5b3UK
```

This output can be then passed to Kubernetes as follows:

```
tgen -f secret.yaml | kubectl apply -f -
```

Do keep in mind though your DevOps requirements in terms of keeping a copy of your YAML files, rendered. Additionally, the `readfile` function is akin to `helm`'s `.Files`, with the exception that **you can read any file the `tgen` binary has access**, including potentially sensitive files such as `/etc/passwd`. If this is a concern, please run `tgen` in a CI/CD environment or where access to these resources is limited.

### Template functions

See [template functions](docs/functions.md) for a list of all the functions available.
