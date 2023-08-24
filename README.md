# `tgen`: a tiny template tool

[![Downloads](https://img.shields.io/github/downloads/patrickdappollonio/tgen/total?color=blue&logo=github&style=flat-square)](https://github.com/patrickdappollonio/tgen/releases)

`tgen` is a simple CLI application that allows you to write a template file and then use the power of Go Templates to generate an output (which is) outputted to `stdout`. Besides the Go Template engine itself, `tgen` contains a few extra utility functions to assist you when writing templates. See below for a description of each.

You can also use the `--help` (or `-h`) to see the available set of options. The only flag required is the file to process, and everything else is optional.

```
tgen is a template generator with the power of Go Templates

Usage:
  tgen [flags]

Flags:
  -e, --environment string   an optional environment file to use (key=value formatted) to perform replacements
  -f, --file string          the template file to process, or "-" to read from stdin
  -d, --delimiter string     template delimiter (default "{{}}")
  -x, --execute string       a raw template to execute directly, without providing --file
  -v, --values string        a file containing values to use for the template, a la Helm
      --with-values          automatically include a values.yaml file from the current working directory
  -s, --strict               strict mode: if an environment variable or value is used in the template but not set, it fails rendering
  -h, --help                 help for tgen
      --version              version for tgen
```

## Usage

You can use `tgen`:

* By reading environment variables from the environment or a key-value file
* By reading variables from a YAML values file

While working with it, `tgen` supports a "strict" mode, where if a variable (either environment or from a values file) is used in the template but not set, it will fail the template generation.

## Examples

### Simple template

Using a template file and an environment file, you can generate a template as follows:

```bash
$ cat template.txt
The dog licked the {{ env "element" }} and everyone laughed.

$ cat contents.env
element=Oil

$ tgen -e contents.env -f template.txt
The dog licked the Oil and everyone laughed.
```

### Inline mode

You can skip the template file altogether and use the inline mode to execute a template directly:

```bash
$ cat contents.env
element=Oil

$ tgen -e contents.env -x 'The dog licked the {{ env "element" }} and everyone laughed.'
The dog licked the Oil and everyone laughed.
```

### Helm-style values

While `tgen` Helm-like support is currently limited to values, it allows for a powerful way to generate templates. Consider the following example:

```bash
$ cat template.txt
The dog licked the {{ .element }} and everyone laughed.

$ cat values.yaml
element: Oil

$ tgen -v values.yaml -f template.txt
The dog licked the Oil and everyone laughed.
```

In the last function call, if your file is named `values.yaml`, you can omit it calling it directly and instead use:

```bash
$ tgen --with-values -f template.txt
The dog licked the Oil and everyone laughed.
```

For more details, see the ["Kubernetes and Helm-style values" documentation page](docs/helm-style-values.md).

## Template functions

See [template functions](docs/functions.md) for a list of all the functions available.
