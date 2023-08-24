# Kubernetes & Helm-style values

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
