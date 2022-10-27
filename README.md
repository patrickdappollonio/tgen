# `tgen`: a tiny template tool

[![Tests passing](https://img.shields.io/github/workflow/status/patrickdappollonio/tgen/Continuous%20Integration/master?logo=github&style=flat-square)](https://github.com/patrickdappollonio/tgen/actions)
[![Downloads](https://img.shields.io/github/downloads/patrickdappollonio/tgen/total?color=blue&logo=github&style=flat-square)](https://github.com/patrickdappollonio/tgen/releases)


`tgen` is a simple CLI application that allows you to write a template file and then use the power of Go Templates to generate an output (which is) outputted to `stdout`. Besides the Go Template engine itself, `tgen` contains a few extra utility functions to assist you when writing templates. See below for a description of each.

You can also use the `--help` (or `-h`) to see the available set of options. The only flag required is the file to process, and everything else is optional.

```
tgen is a template generator with the power of Go Templates

Usage:
  tgen [flags]

Flags:
  -e, --environment string   an optional environment file to use (key=value formatted) to perform replacements
  -f, --file string          the template file to process
  -d, --delimiter string     template delimiter (default "{{}}")
  -x, --execute string       a raw template to execute directly, without providing --file
  -v, --values string        a file containing values to use for the template, a la Helm
      --with-values          automatically include a values.yaml file from the current working directory
  -s, --strict               strict mode: if an environment variable or value is used in the template but not set, it fails rendering
  -h, --help                 help for tgen
      --version              version for tgen
```

### Environment file

`tgen` supports an optional environment variable collection in a file but it's a pretty basic implementation of a simple key/value pair. The environment file works by finding lines that aren't empty or preceded by a pound `#` -- since they're treated as comments -- and then tries to find at least one equal (`=`) sign. If it can find at least one, all values on the left side of the equal sign become the key and the contents on the right side become the value. If the same line has more than one equal, only the first one is honored and all remaining ones become part of the value.

There's no support for Bash interpolation or multiline values. If this is needed, consider using a YAML values file instead.

#### Example

Consider the following template, named `template.txt`:

```handlebars
The dog licked the {{ env "element" }} and everyone laughed.
```

And the following environment file, named `contents.env`:

```bash
element=Oil
```

After being passed to `tgen`, the output becomes:

```bash
$ tgen -e contents.env -f template.txt
The dog licked the Oil and everyone laughed.
```

Using the inline mode to execute a template, you can also call the program as such (note the use of single-quotes since in Go, strings are always double-quoted) which will yield the same result:

```bash
$ tgen -x '{{ env "element" }}' -e contents.env
The dog licked the Oil and everyone laughed.
```

Do note as well that using single quotes for the template allows you to prevent any bash special parsing logic that your terminal might have.

### Helm-style values

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
    Rk9PQkFSQkFaCg==
```

This output can be then passed to Kubernetes as follows:

```
tgen -f secret.yaml | kubectl apply -f -
```

Do keep in mind though your DevOps requirements in terms of keeping a copy of your YAML files, rendered. Additionally, the `readfile` function is akin to `helm`'s `.Files`, with the exception that **you can read any file the `tgen` binary has access**, including potentially sensitive files such as `/etc/passwd`. If this is a concern, please run `tgen` in a CI/CD environment or where access to these resources is limited.

You can also use a `values.yaml` file like Helm. `tgen` will allow you to read values from the values file as `.variable` or `.Values.variable`. The latter is the same as Helm's `.Values.variable` and the former is a shortcut to `.Values.variable` for convenience. Consider the following YAML values file:

```yaml
name: Patrick
```

And the following template:

```handlebars
Hello, my name is {{ .name }}.
```

Running `tgen` with the values file will yield the following output:

```bash
$ tgen -f template.yaml -v values.yaml
Hello, my name is Patrick.
```

If your values file is called `values.yaml`, you also have the handy shortcut of simply specifying `--with-values` and `tgen` will automatically include the values file from the current working directory.

### Template functions

See [template functions](docs/functions.md) for a list of all the functions available.
