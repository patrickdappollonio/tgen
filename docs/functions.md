# Template functions

- [Template functions](#template-functions)
  - [Sprig functions](#sprig-functions)
  - [Additional `tgen` specific functions](#additional-tgen-specific-functions)
    - [`raw`](#raw)
    - [`lowercase`](#lowercase)
    - [`uppercase`](#uppercase)
    - [`sprintf`, `printf`, `println`](#sprintf-printf-println)
    - [`env`, `envdefault`](#env-envdefault)
    - [`rndstring`](#rndstring)
    - [`base64encode`, `base64decode`](#base64encode-base64decode)
    - [`readfile`, `readlocalfile`](#readfile-readlocalfile)
    - [`readdir`, `readlocaldir`, `readdirrecursive`, `readlocaldirrecursive`](#readdir-readlocaldir-readdirrecursive-readlocaldirrecursive)
    - [`linebyline`, `lbl`](#linebyline-lbl)
    - [`after`, `skip`](#after-skip)
    - [`required`](#required)

All examples below have been generated using `-x` -- or `--execute`, which allows passing a template as argument rather than reading a file. In either case, whether the template file -- with `-f` or `--file` -- or the template argument is used, all functions are available.

Each function includes a set of examples. The lines prepended with a `$` are bash commands you can try by running them on your terminal.

## Sprig functions

`tgen` includes all functions from [Sprig](https://masterminds.github.io/sprig/), which are the same functions you could be used to if you have ever used Helm. This includes:

* [String Functions](https://masterminds.github.io/sprig/strings.html): `trim`, `wrap`, `randAlpha`, `plural`, etc.
  * [String List Functions](https://masterminds.github.io/sprig/string_slice.html): `splitList`, `sortAlpha`, etc.
* [Integer Math Functions](https://masterminds.github.io/sprig/math.html): `add`, `max`, `mul`, etc.
    * [Integer Slice Functions](https://masterminds.github.io/sprig/integer_slice.html): `until`, `untilStep`
* [Float Math Functions](https://masterminds.github.io/sprig/mathf.html): `addf`, `maxf`, `mulf`, etc.
* [Date Functions](https://masterminds.github.io/sprig/date.html): `now`, `date`, etc.
* [Defaults Functions](https://masterminds.github.io/sprig/defaults.html): `default`, `empty`, `coalesce`, `fromJson`, `toJson`, `toPrettyJson`, `toRawJson`, `ternary`
* [Encoding Functions](https://masterminds.github.io/sprig/encoding.html): `b64enc`, `b64dec`, etc.
* [Lists and List Functions](https://masterminds.github.io/sprig/lists.html): `list`, `first`, `uniq`, etc.
* [Dictionaries and Dict Functions](https://masterminds.github.io/sprig/dicts.html): `get`, `set`, `dict`, `hasKey`, `pluck`, `dig`, `deepCopy`, etc.
* [Type Conversion Functions](https://masterminds.github.io/sprig/conversion.html): `atoi`, `int64`, `toString`, etc.
* [Path and Filepath Functions](https://masterminds.github.io/sprig/paths.html): `base`, `dir`, `ext`, `clean`, `isAbs`, `osBase`, `osDir`, `osExt`, `osClean`, `osIsAbs`
* [Flow Control Functions](https://masterminds.github.io/sprig/flow_control.html): `fail`
* Advanced Functions
    * [UUID Functions](https://masterminds.github.io/sprig/uuid.html): `uuidv4`
    * [OS Functions](https://masterminds.github.io/sprig/os.html): `env`, `expandenv`
    * [Version Comparison Functions](https://masterminds.github.io/sprig/semver.html): `semver`, `semverCompare`
    * [Reflection](https://masterminds.github.io/sprig/reflection.html): `typeOf`, `kindIs`, `typeIsLike`, etc.
    * [Cryptographic and Security Functions](https://masterminds.github.io/sprig/crypto.html): `derivePassword`, `sha256sum`, `genPrivateKey`, etc.
    * [Network](https://masterminds.github.io/sprig/network.html): `getHostByName`

## Additional `tgen` specific functions

These are functions that are not part of Sprig, but are included in `tgen` for convenience.

### `raw`

Raw returns the value provided as a string. It's kept for backwards compatibility and non-breaking old resources:

```bash
$ tgen -x '{{ "hello" | raw }}'
hello
```

### `lowercase`

Converts the string to a lowercase value:

```bash
$ tgen -x '{{ "HELLO" | lowercase }}'
hello
```

### `uppercase`

Converts the string to a uppercase value:

```bash
$ tgen -x '{{ "hello" | uppercase }}'
HELLO
```

### `sprintf`, `printf`, `println`

Functions akin to Go's own `fmt.Sprintf` and `fmt.Sprintln`. `printf` is an alias of `sprintf`:

```bash
$ tgen -x '{{ sprintf "Hello, %s!" "World" }}'
Hello, World!
```

### `env`, `envdefault`

Functions to grab environment variable values. For `env`, the value will be printed out or be empty if the environment variable is not set. For `envdefault`, the value will be the value retrieved from the environment variable or the default value specified.

Both `env` and `envdefault` are case insensitive -- either `"home"` or `"HOME"` will work.

When `--strict` mode is enabled, if `env` is called with a environment variable name with no value set or set to empty, the application will exit with error. Useful if you must receive a value or fail a CI build, for example.

Consider the following example reading these environment variables:

```bash
$ tgen -x '{{ env "user" }}'
patrick

$ tgen -x '{{ env "USER" }}'
patrick
```

Then trying to read a nonexistent environment variable with `--strict` mode enabled:

```bash
$ tgen -x '{{ env "foobar" }}' --strict
Error: evaluating /dev/stdin:1:3: strict mode on: environment variable not found: $FOOBAR
```

And bypassing strict mode by setting a default value:

```bash
$ tgen -x '{{ envdefault "SQL_HOST" "sql.example.com" }}' --strict
sql.example.com
```

For custom messages, [consider using `required` instead](#required).

### `rndstring`

Generates a random string of a given length:

```bash
$ tgen -x '{{ rndstring 8 }}'
mHNmtrbf
```

### `base64encode`, `base64decode`

Functions to encode and decode from `base64`. These are also available from Sprig as `b64enc` and `b64dec`.

```bash
$ tgen -x '{{ base64encode "hello" }}'
aGVsbG8=
```

```bash
$ tgen -x '{{ base64decode "aGVsbG8=" }}'
hello
```

### `readfile`, `readlocalfile`

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

### `readdir`, `readlocaldir`, `readdirrecursive`, `readlocaldirrecursive`

Read a directory from a local path -- either relative or absolute -- and returns it as an array of strings, which can be used to iterate over the files in the directory.

`readdir` and `readlocaldir` do not recurse into subdirectories, while `readdirrecursive` and `readlocaldirrecursive` do.

```bash
$ tree testdata
testdata
├── file1.txt
└── file2.txt

$ tgen -x '{{ readdir "testdata" }}'
[file1.txt file2.txt]
```

```bash
$ tgen -x '{{ readlocaldir "testdata" }}'
[file1.txt file2.txt]
```

Attempting to read a directory that does not exist will return an error:

```bash
$ tgen -x '{{ readdir "doesnotexist" }}'
Error: template: tgen:1:3: executing "tgen" at <readdir "doesnotexist">: error calling readdir: open doesnotexist: no such file or directory
```

And attempting to read a directory outside the current working directory with `readlocaldir` will return an error:

```bash
$ tgen -x '{{ readlocaldir "../testdata" }}'
Error: template: tgen:1:3: executing "tgen" at <readlocaldir "../testdata">: error calling readlocaldir: unable to open local directory "../testdata": directory is not under current working directory
```

With the recursive functions, the same rules apply but they will also include the files and folders in the subdirectories. Folders will be returned as strings with a trailing `/`.

```bash
$ tree testdata
testdata
├── file1.txt
├── file2.txt
└── subdir
    ├── subfile1.txt
    └── subfile2.txt

$ tgen -x '{{ readdirrecursive "testdata" }}'
[file1.txt file2.txt subdir/ subdir/subfile1.txt subdir/subfile2.txt]
```

```bash
$ tgen -x '{{ readlocaldirrecursive "testdata" }}'
[file1.txt file2.txt subdir/ subdir/subfile1.txt subdir/subfile2.txt]
```

Some considerations:

* Symbolic links are not followed.
* If a relative path is provided, all paths must be relative to the current working directory.
* For `readdir` and `readdirrecursive`, the path can be either relative or absolute:
  * Any directory can be read through `readdir`, and yes, that includes `/etc` and other sensitive files. If this level of security is important to you, consider running `tgen` in trusted environments. This is by design to allow embedding files from other folders external to the current working directory and its subdirectories.
  * If reading any directory is a problem, consider using `readlocaldir` or `readlocaldirrecursive`.
* For `readlocaldir` and `readlocaldirrecursive`, the path can only be relative:
  * Absolute paths will return in an error.
  * The current working directory will be prepended to the path provided.
  * Only directories within the current working directory and its subdirectories can be read through this function.

### `linebyline`, `lbl`

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

### `after`, `skip`

Returns a Go slice to only the items after the `n`th item. Negative numbers for `after` are not supported and will result in an error.

```bash
# Creates a sequence from 1 to 5, then
# returns all values after 2
$ tgen -x '{{ after 2 (seq 5) }}'
[3 4 5]

# Alternate way of writing it
$ tgen -x '{{ seq 5 | after 2 }}'
[3 4 5]
```

### `required`

Returns an error if the value is empty. Useful to ensure a value is provided, and if not, fail the template generation.

```bash
$ tgen -x '{{ env "foo" | required "environment variable \"foo\" is required" }}'
Error: evaluating /dev/stdin:1:15: environment variable "foo" is required
```

Note that you can also use `--strict` mode to achieve a similar result. The difference between `--strict` and `required` is that `required` works anywhere: not just on missing YAML value keys or environment variables. Here's another example:

```bash
$ tgen -x '{{ "" | required "Value must be set" }}'
Error: evaluating /dev/stdin:1:8: Value must be set
```
