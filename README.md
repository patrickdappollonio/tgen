# `tgen`

[![Release](https://github.com/patrickdappollonio/tgen/actions/workflows/go-release.yml/badge.svg)](https://github.com/patrickdappollonio/tgen/releases)


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

All examples below have been generated using `-x` -- or `--execute`, which allows passing a template as argument rather than reading a file. In either case, whether the template file -- with `-f` or `--file` -- or the template argument is used, all functions are available.

Each function includes a set of examples. The lines prepended with a `$` are bash commands you can try by running them on your terminal.

#### `raw`

Raw returns the value provided as a string. It's kept for backwards compatibility and non-breaking old resources:

```bash
$ tgen -x '{{ "hello" | raw }}'
hello
```

#### `lowercase`, `lower`

Converts the string to a lowercase value:

```bash
$ tgen -x '{{ "HELLO" | lowercase }}'
hello
```

#### `uppercase`, `upper`

Converts the string to a uppercase value:

```bash
$ tgen -x '{{ "hello" | uppercase }}'
HELLO
```

#### `title`

Converts the first letter of each word to uppercase:

```bash
$ tgen -x '{{ "hello world" | title }}'
Hello World
```

#### `sprintf`, `printf`, `println`

Functions akin to Go's own `fmt.Sprintf` and `fmt.Sprintln`. `printf` is an alias of `sprintf`:

```bash
$ tgen -x '{{ sprintf "Hello, %s!" "World" }}'
Hello, World!
```

#### `trim`, `trimPrefix`, `trimSuffix`

Trim empty spaces, a prefix or a suffix:

```bash
$ tgen -x '{{ trim "   hello   " }}'
hello
```

```bash
$ tgen -x '{{ trimPrefix "hello" "h" }}'
ello
```

```bash
$ tgen -x '{{ trimSuffix "hello" "o" }}'
hell
```

#### `split`

Splits a string on a given character:

```bash
$ tgen -x '{{ split "Hello World" " " }}'
[Hello World]
```

```bash
$ tgen -x '{{ range split "Hello World" " "  }}{{ printf "%s\n" . }}{{ end }}'
Hello
World

```

#### `base`, `dir`, `clean`, `ext`, `isAbs`

Functions to use when handling directories:

```bash
$ tgen -x '{{ base "/foo/bar/baz" }}'
baz
```

```bash
$ tgen -x '{{ dir "/foo/bar/baz" }}'
/foo/bar
```

```bash
$ tgen -x '{{ clean "/foo/bar/../baz" }}'
/foo/baz
```

```bash
$ tgen -x '{{ ext "/foo.zip" }}'
.zip
```

```bash
$ tgen -x '{{ isAbs "foo.zip" }}'
false
```

#### `env`, `envdefault`

Functions to grab environment variable values. For `env`, the value will be printed out or be empty if the environment variable is not set. For `envdefault`, the value will be the value retrieved from the environment variable or the default value specified.

Both `env` and `envdefault` are case insensitive -- either `"home"` or `"HOME"` will work.

When `--strict` mode is enabled, if `env` is called with a environment variable name with no value set or set to empty, the application will exit with error. Useful if you must receive a value or fail a CI build, for example.

```bash
$ tgen -x '{{ env "user" }}'
patrick

$ tgen -x '{{ env "USER" }}'
patrick
```

```bash
$ tgen -x '{{ env "foobar" }}' -s
Error: strict mode on: environment variable not found: $FOOBAR
```

```bash
$ tgen -x '{{ envdefault "SQL_HOST" "sql.example.com" }}'
sql.example.com
```

#### `rndstring`

Generates a random string of a given length:

```bash
$ tgen -x '{{ rndstring 8 }}'
mHNmtrbf
```

#### `repeat`

Repeats a string a given amount of times:

```bash
$ tgen -x '{{ repeat 3 "abc" }}'
abcabcabc
```

#### `nospace`

Removes all spaces from a string:

```bash
$ tgen -x '{{ nospace "Lorem ipsum dolor sit amet" }}'
Loremipsumdolorsitamet
```

#### `quote`, `squote`

Wrap a string in single quotes -- with `squote` -- or double quotes -- with `quote`:

```bash
$ tgen -x '{{ quote "Hello" }}'
"Hello"
```

```bash
$ tgen -x '{{ squote "Hello" }}'
'Hello'
```

#### `replace`

Replaces a substring of a string. The first parameter is the old value, the second parameter is the new value, and the third parameter is the complete string you want to perform the replacement.

```bash
$ tgen -x '{{ replace "World" "Patrick" "Hello, World!" }}'
Hello, Patrick!
```

Go developers might realize the parameters are in different positions compared to `strings.Replace()`: this is intentional, to allow the Go template engine to have the last parameter to be piped:

```bash
$ tgen -x '{{ raw "Hello, World!" | replace "World" "Patrick" }}'
Hello, Patrick!
```

#### `indent`, `nindent`

Adds spaces at the beginning of each line for indentation. `nindent` will also add a starting line break before the string. Both functions are especially useful when dealing with YAML files.

```bash
$ tgen -x '{{ "name: patrick" | indent 4 }}'
    name: patrick
```

```bash
$ tgen -x '{{ "name: patrick" | nindent 4 }}'

    name: patrick
```

A more complete example, in YAML:

```bash
$ cat tests/basic.yaml
user_details:
  username: {{ env "user" }}
  home_directory: | {{ env "home" | nindent 4 }}

$ tgen -f tests/basic.yaml
user_details:
  username: patrick
  home_directory: |
    /home/patrick
```

#### `b64enc`, `base64encode`, `b64dec`, `base64decode`

Functions to encode and decode from `base64`.

```bash
$ tgen -x '{{ b64enc "hello" }}'
aGVsbG8=
```

```bash
$ tgen -x '{{ base64encode "hello" }}'
aGVsbG8=
```

```bash
$ tgen -x '{{ b64dec "aGVsbG8=" }}'
hello
```

```bash
$ tgen -x '{{ base64decode "aGVsbG8=" }}'
hello
```

#### `sha1sum`, `sha256sum`

Generates a `SHA1` hash or a `SHA256` hash of a given string:

```bash
$ tgen -x '{{ sha1sum "hello" }}'
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

```bash
$ tgen -x '{{ sha256sum "hello" }}'
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```

#### `readfile`, `readlocalfile`

Read a file from a local path -- either relative or absolute -- and print it as a string. Useful to embed files from your local machine or CI environment into your template:

```bash
$ tgen -x '{{ readfile "/etc/hostname" }}'
localhost
```

```bash
$ tgen -x '{{ readlocalfile "go.mod" }}'
module github.com/patrickdappollonio/tgen

require github.com/spf13/cobra v1.2.1

go 1.16
```

```bash
$ tgen -x '{{ readlocalfile "../etc/hosts" }}'
Error: template: tgen:1:3: executing "tgen" at <readlocalfile "../etc/hosts">: error calling readlocalfile: unable to open local file "/etc/hosts": file is not under current working directory
```

Some considerations:

* If a relative path is provided, all paths must be relative to the current working directory.
  * If a template is inside a subfolder from the current working directory, the path you must provide in `readfile` has to be starting from the current working directory, not from the location where the template file is.
* For `readfile`, the path can be eithe relative or absolute:
  * Any file can be read through `readfile`, and yes, that includes `/etc/passwd` and other sensitive files. If this level of security is important to you, consider running `tgen` in trusted environments. This is by design to allow embedding files from other folders external to the current working directory and its subdirectories.
  * If reading any file is a problem, consider using `readlocalfile`.
* For `readlocalfile`, the path can only be relative:
  * Absolute paths will return in an error.
  * The current working directory will be prepended to the path provided.
  * Only files within the current working directory and its subdirectories can be read through this function.

For a more complete example, see [Template Generation _a la Helm_](#template-generation-a-la-helm).

#### `linebyline`, `lbl`

Parses the input and splits on line breaks. `linebyline` is a shorcut for [`split`](#split) with a split character of `\n`. `lbl` is an alias of `linebyline`:

```bash
$ tgen -x '{{ linebyline "foo\nbar" }}'
[foo bar]
```

```bash
$ tgen -x '{{ range linebyline "foo\nbar" }}{{ . | nindent 2 }}{{ end }}'

  foo
  bar
```
