# Template functions

- [Template functions](#template-functions)
  - [`raw`](#raw)
  - [`lowercase`, `lower`](#lowercase-lower)
  - [`uppercase`, `upper`](#uppercase-upper)
  - [`title`](#title)
  - [`sprintf`, `printf`, `println`](#sprintf-printf-println)
  - [`trim`, `trimPrefix`, `trimSuffix`](#trim-trimprefix-trimsuffix)
  - [`split`](#split)
  - [`base`, `dir`, `clean`, `ext`, `isAbs`](#base-dir-clean-ext-isabs)
  - [`env`, `envdefault`](#env-envdefault)
  - [`rndstring`](#rndstring)
  - [`repeat`](#repeat)
  - [`nospace`](#nospace)
  - [`quote`, `squote`](#quote-squote)
  - [`replace`](#replace)
  - [`indent`, `nindent`](#indent-nindent)
  - [`b64enc`, `base64encode`, `b64dec`, `base64decode`](#b64enc-base64encode-b64dec-base64decode)
  - [`sha1sum`, `sha256sum`](#sha1sum-sha256sum)
  - [`readfile`, `readlocalfile`](#readfile-readlocalfile)
  - [`linebyline`, `lbl`](#linebyline-lbl)
  - [`seq`](#seq)
  - [`slice`, `list`](#slice-list)

All examples below have been generated using `-x` -- or `--execute`, which allows passing a template as argument rather than reading a file. In either case, whether the template file -- with `-f` or `--file` -- or the template argument is used, all functions are available.

Each function includes a set of examples. The lines prepended with a `$` are bash commands you can try by running them on your terminal.

## `raw`

Raw returns the value provided as a string. It's kept for backwards compatibility and non-breaking old resources:

```bash
$ tgen -x '{{ "hello" | raw }}'
hello
```

## `lowercase`, `lower`

Converts the string to a lowercase value:

```bash
$ tgen -x '{{ "HELLO" | lowercase }}'
hello
```

## `uppercase`, `upper`

Converts the string to a uppercase value:

```bash
$ tgen -x '{{ "hello" | uppercase }}'
HELLO
```

## `title`

Converts the first letter of each word to uppercase:

```bash
$ tgen -x '{{ "hello world" | title }}'
Hello World
```

## `sprintf`, `printf`, `println`

Functions akin to Go's own `fmt.Sprintf` and `fmt.Sprintln`. `printf` is an alias of `sprintf`:

```bash
$ tgen -x '{{ sprintf "Hello, %s!" "World" }}'
Hello, World!
```

## `trim`, `trimPrefix`, `trimSuffix`

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

## `split`

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

## `base`, `dir`, `clean`, `ext`, `isAbs`

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

## `env`, `envdefault`

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

## `rndstring`

Generates a random string of a given length:

```bash
$ tgen -x '{{ rndstring 8 }}'
mHNmtrbf
```

## `repeat`

Repeats a string a given amount of times:

```bash
$ tgen -x '{{ repeat 3 "abc" }}'
abcabcabc
```

## `nospace`

Removes all spaces from a string:

```bash
$ tgen -x '{{ nospace "Lorem ipsum dolor sit amet" }}'
Loremipsumdolorsitamet
```

## `quote`, `squote`

Wrap a string in single quotes -- with `squote` -- or double quotes -- with `quote`:

```bash
$ tgen -x '{{ quote "Hello" }}'
"Hello"
```

```bash
$ tgen -x '{{ squote "Hello" }}'
'Hello'
```

## `replace`

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

## `indent`, `nindent`

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

## `b64enc`, `base64encode`, `b64dec`, `base64decode`

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

## `sha1sum`, `sha256sum`

Generates a `SHA1` hash or a `SHA256` hash of a given string:

```bash
$ tgen -x '{{ sha1sum "hello" }}'
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

```bash
$ tgen -x '{{ sha256sum "hello" }}'
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```

## `readfile`, `readlocalfile`

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

## `linebyline`, `lbl`

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

## `seq`

Generates a sequence of numbers based on how Unix's `seq` works, using a start, end, and step parameters: providing only an end parameter; providing starting and ending parameters; and providing starting, step, and ending parameters:

```bash
# Providing only an end number
$ tgen -x '{{ seq 10 }}'
[1 2 3 4 5 6 7 8 9 10]

# Providing a start and end number
$ tgen -x '{{ seq 10 20 }}'
[10 11 12 13 14 15 16 17 18 19 20]

# Providing a start, end, and step number
$ tgen -x '{{ seq 2 2 20 }}'
[2 4 6 8 10 12 14 16 18 20]
```

It's also possible to use reverse:

```bash
$ tgen -x '{{ seq 10 -1 1 }}'
[10 9 8 7 6 5 4 3 2 1]
```

## `slice`, `list`

Creates a Go slice -- essentially an array -- from a list of values:

```bash
$ tgen -x '{{ slice 1 2 3 }}'
[1 2 3]

$ tgen -x '{{ slice 1 2 3 "a" "b" "c" }}'
[1 2 3 a b c]
```
