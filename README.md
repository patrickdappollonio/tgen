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
  -x, --execute string       a raw template to execute directly, without providing --file
  -f, --file string          the template file to process (required)
  -h, --help                 help for tgen
  -s, --strict               enables strict mode: if an environment variable in the file is defined but not set, it'll fail
      --version              version for tgen
```

### Available functions

All available functions are in the top of [`template_functions.go`](template_functions.go). For convenience, here's a list of them:

```go
// Go built-ins
"lowercase":  strings.ToLower,    // "HELLO" → "hello"
"lower":      strings.ToLower,    // "HELLO" → "hello"
"uppercase":  strings.ToUpper,    // "hello" → "HELLO"
"upper":      strings.ToUpper,    // "hello" → "HELLO"
"title":      strings.Title,      // "hello" → "Hello"
"sprintf":    fmt.Sprintf,        // sprintf "Hello, %s" "world" → "Hello, world"
"printf":     fmt.Sprintf,        // printf "Hello, %s" "world" → "Hello, world"
"println":    fmt.Sprintln,       // println "Hello" "world!" → "Hello world!\n"
"trim":       strings.TrimSpace,  // trim "   hello   " → "hello"
"trimPrefix": strings.TrimPrefix, // trimPrefix "abcdef" "abc" → "def"
"trimSuffix": strings.TrimSuffix, // trimSuffix "abcdef" "def" → "abc"
"base":       filepath.Base,      // base "/foo/bar/baz" → "baz"
"dir":        filepath.Dir,       // dir "/foo/bar/baz" → "/foo/bar"
"clean":      filepath.Clean,     // clean "/foo/bar/../baz" → "/foo/baz"
"ext":        filepath.Ext,       // ext "/foo.zip" → ".zip"
"isAbs":      filepath.IsAbs,     // isAbs "foo.zip" → false

// Locally defined functions
"env":          envstrict(strict), // env "user" → "patrick"
"envdefault":   envdefault,        // env "SQL_HOST" "sql.example.com" → "sql.example.com"
"rndstring":    rndgen,            // rndstring 8 → "lFEqUUOJ"
"repeat":       repeat,            // repeat 3 "abc" → "abcabcabc"
"nospace":      nospace,           // nospace "hello world!" → "helloworld!"
"quote":        quote,             // quote "hey" → `"hey"`
"squote":       squote,            // squote "hey" → "'hey'"
"indent":       indent,            // indent 3 "abc" → "  abc"
"nindent":      nindent,           // nindent 3 "abc" → "\n   abc"
"b64enc":       base64encode,      // b64enc "abc" → "YWJj"
"base64encode": base64encode,      // base64encode "abc" → "YWJj"
"b64dec":       base64decode,      // b64dec "YWJj" → "abc"
"base64decode": base64decode,      // base64decode "YWJj" → "abc"
"sha1sum":      sha1sum,           // sha1sum "abc" → "a9993e364706816aba3e25717850c26c9cd0d89d"
"sha256sum":    sha256sum,         // sha256sum "abc" → "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
"replace":      replace,           // replace "World" "Patrick" "Hello, World!" → "Hello, Patrick!"
"readfile":     readfile,          // readfile "foobar.txt" → "Hello, world!"
"linebyline":   linebyline,        // linebyline "foo\nbar" → ["foo", "bar"]
"lbl":          linebyline,        // linebyline "foo\nbar" → ["foo", "bar"]
```

Some of them expose two functions with different names but same value output: for example, `b64enc` and `base64encode` both encode to `base64` the same way.

### Environment file

`tgen` supports an optional environment variable collection in a file but it's a pretty basic implementation of a simple key/value pair. The environment file works by finding lines that aren't empty or preceded by a pound `#` -- since they're treated as comments -- and then tries to find at least one equal (`=`) sign. If it can find at least one, all values on the left side of the equal sign become the key -- which is also uppercased so it's compatible with the `env` function defined above -- and the contents on the right side become the value. If the same line has more than one equal, only the first one is honored and all remaining ones become part of the value.

As an important note, environment variables found in the environment have preference over the environment file. That way, the environment file can define `A=1` but then the application can be run with `A=2 tgen [flags]` so it overrides `A` to the value of `2`.

## Example

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

## Template Generation _a la Helm_

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
