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

## Helm-style `--set` and `--set-string` flags

Similar to Helm, `tgen` supports the `--set` and `--set-string` flags to set values directly from the command line. These flags allow you to override values without needing a separate values file.

### Basic Usage

The `--set` flag allows you to set values with automatic type inference:

```bash
tgen -f template.yaml --set name=Patrick --set age=30 --set active=true
```

The `--set-string` flag forces all values to be treated as strings:

```bash
tgen -f template.yaml --set-string name=Patrick --set-string age=30 --set-string active=true
```

### Type Inference with `--set`

The `--set` flag automatically infers the type of values:

- **Booleans**: `true`, `false`, `yes`, `no`, `on`, `off` (case-insensitive)
- **Integers**: `42`, `-10`, `0`
- **Floats**: `3.14`, `-2.5`, `0.0`
- **Strings**: Everything else, including `"quoted strings"`

```bash
# These create different types
tgen -f template.yaml --set debug=true --set replicas=3 --set version=1.2.3
```

### YAML Boolean Edge Cases

The `--set` flag properly handles YAML boolean values, including country codes:

```bash
# These all create boolean values
tgen -f template.yaml --set user.admin=yes --set user.active=on --set user.country=no

# This creates a string value
tgen -f template.yaml --set-string user.country=no
```

### Nested Values

Use dot notation to create nested structures:

```bash
tgen -f template.yaml --set app.name=myapp --set app.version=1.0.0 --set database.host=localhost
```

This creates:
```yaml
app:
  name: myapp
  version: 1.0.0
database:
  host: localhost
```

### Comma-separated Values

Set multiple values in a single `--set` or `--set-string` flag:

```bash
tgen -f template.yaml --set 'app.name=myapp,app.version=1.0.0,replicas=3'
```

### Array Syntax

Create arrays using curly braces:

```bash
tgen -f template.yaml --set 'tags={web,api,database}' --set 'ports={80,443,8080}'
```

This creates:
```yaml
tags:
  - web
  - api
  - database
ports:
  - 80
  - 443
  - 8080
```

### Array Indexing

Set specific array elements using bracket notation:

```bash
tgen -f template.yaml --set 'servers[0].host=web1' --set 'servers[0].port=80' --set 'servers[1].host=web2'
```

This creates:
```yaml
servers:
  - host: web1
    port: 80
  - host: web2
```

### Special Values

Handle null values and empty arrays:

```bash
# Null values
tgen -f template.yaml --set 'database.password=null'

# Empty arrays
tgen -f template.yaml --set 'tags=[]'
```

### Escaping Special Characters

Escape commas, dots, and other special characters:

```bash
# Escape commas in values
tgen -f template.yaml --set 'message=Hello\, World!'

# Escape dots in keys
tgen -f template.yaml --set 'nodeSelector."kubernetes\.io/role"=master'
```

### Complex Example

Here's a comprehensive example combining multiple features:

```bash
tgen -f deployment.yaml \
  --set 'app.name=myapp,app.version=2.1.0,app.debug=true' \
  --set 'replicas=3' \
  --set 'image.repository=myregistry/myapp,image.tag=v2.1.0' \
  --set 'env={PROD,STAGING}' \
  --set 'servers[0].host=web1.example.com,servers[0].port=80' \
  --set 'servers[1].host=web2.example.com,servers[1].port=80' \
  --set 'database.host=db.example.com,database.port=5432,database.password=null' \
  --set-string 'annotations."deployment\.kubernetes\.io/revision"=1'
```

### Combining with Values Files

You can combine `--set` and `--set-string` with values files. Command-line values take precedence:

```bash
tgen -f template.yaml -v values.yaml --set 'app.debug=true' --set-string 'app.version=override'
```

### Template Usage

In your templates, access these values just like values from files:

```handlebars
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .app.name }}
  labels:
    version: {{ .app.version }}
spec:
  replicas: {{ .replicas }}
  selector:
    matchLabels:
      app: {{ .app.name }}
  template:
    metadata:
      labels:
        app: {{ .app.name }}
    spec:
      containers:
      - name: {{ .app.name }}
        image: {{ .image.repository }}:{{ .image.tag }}
        ports:
        {{- range .servers }}
        - containerPort: {{ .port }}
        {{- end }}
        env:
        {{- range .env }}
        - name: ENVIRONMENT
          value: {{ . }}
        {{- end }}
```

### Best Practices

1. **Use `--set` for configuration values** where type inference is important
2. **Use `--set-string` for version numbers, IDs, or other string-like values** that should remain strings
3. **Quote complex values** to avoid shell interpretation issues
4. **Combine multiple related values** using comma separation for readability
5. **Use values files for large configurations** and `--set` for overrides
