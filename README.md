# Header transformation plugin for Traefik

[![Build Status](https://travis-ci.com/tomMoulard/htransformation.svg?branch=main)](https://travis-ci.com/tomMoulard/htransformation)

This plugin allows to change, on the fly, the header's value of a request.

## Dev `traefik.yml` configuration file for traefik

```yml
pilot:
  token: [REDACTED]

experimental:
  devPlugin:
    goPath: /home/tm/go
    moduleName: github.com/tommoulard/htransformation

entryPoints:
  http:
    address: ":8000"
    forwardedHeaders:
      insecure: true

api:
  dashboard: true
  insecure: true

providers:
  file:
    filename: rules-htransformation.yaml
```

## How to dev

```bash
$ docker run -d --network host containous/whoami -port 5000
# traefik --config-file traefik.yml
```

## How to use

To choose a Rule you have to fill the `Type` field with one of the following:

- 'Rename'  : to rename a header
- 'Set'     : to Set a header
- 'Del'     : to Delete a header

Each Rule can be named with the `Name` field

### Rename

A Rule Rename needs two arguments and optionally the third.

- `Header`, the regex of the header you want to replace
- `Value`, the new header
- `HeaderPrefix`, the prefix to denote the new Header name is to be taken from another header value

```yaml
# Example Rename
- Rule:
      Name: 'Header rename'
      Header: 'Cache-Control'
      Value: 'NewHeader'
      Type: 'Rename'
```

```yaml
# Old header:
Cache-Control: gzip, deflate

# New header:
NewHeader: gzip, deflate
```

``` yaml
- Rule:
      Name: 'Header Renaming'
      Header: 'X-Traefik-*'
      Value: 'X-Traefik-merged'
      Type: 'Rename'
```

```yaml
# Old header:
X-Traefik-uuid: 0
X-Traefik-date: mer. 21 oct. 2020 11:57:39 CEST
# New header:
X-Traefik-merged: 0 # A value from old headers
```

### Set

A Set rule will either create or replace the header and value (if it already exists), appending multiple values with the separator if specified.

A Rule Set needs the first two arguments, and optionally the next three.
`Value` can be skipped if specifying `Values`.

- `Header`, the header you want to create
- `Value`, the value of the new header
- `Values`, a list of values to add
- `Sep`, the separator you want to use
- `HeaderPrefix`, the prefix to denote the Value is to be taken from another header

```yaml
# Example Set
- Rule:
      Name: 'Set Cache-Control'
      Header: 'Cache-Control'
      Value: 'Foo'
      Type: 'Set'
```

```yaml
# New header:
Cache-Control: Foo
```

```yaml
# Example Usage
- Rule:
  Name: 'Header set'
  Header: 'X-Forwarded-For'
  Value: '^CF-Connecting-IP'
  HeaderPrefix: "^"
  Type: 'Set'
```

```yaml
# Old header:
CF-Connecting-IP: 1.1.1.1

# New headers:
CF-Connecting-IP: 1.1.1.1
X-Forwarded-For: 1.1.1.1
```

```yaml
# Example Usage
- Rule:
  Name: 'Header XFF'
  Header: 'X-Forwarded-For'
  Value: '^CF-Connecting-IP'
  Values:
    - '^X-Forwarded-For'
    - '192.168.0.1'
  Sep: ', '
  HeaderPrefix: "^"
  Type: 'Set'
```

```yaml
# Old header:
CF-Connecting-IP: 1.1.1.1
X-Forwarded-For: 10.0.0.1, 10.10.10.1
# New headers:
CF-Connecting-IP: 1.1.1.1
X-Forwarded-For: 1.1.1.1, 10.0.0.1, 10.10.10.1, 192.168.0.1
```

```yaml
# Example Join
- Rule:
      Name: 'Header join'
      Header: 'Cache-Control'
      Sep: ','
      HeaderPrefix: "^"
      Values:
        - '^Cache-Control'
        - 'Foo'
        - 'Bar'
      Type: 'Set'
```

```yaml
# Old header:
Cache-Control: gzip, deflate

# Joined header:
Cache-Control: gzip, deflate,Foo,Bar
```

### Delete

A Rule Delete needs only one argument

- `Header`, the header you want to delete

```yaml
# Example Del
- Rule:
      Name: 'Delete Cache-Control'
      Header: 'Cache-Control'
      Type: 'Del'
```

### Point to note

The rules will be evaluated in the order of definition

```yaml
#Example
- Rule:
  Name: 'Header addition'
  Header: 'X-Custom-2'
  Value: 'True'
  Type: 'Set'
- Rule:
  Name: 'Header deletion'
  Header: 'X-Custom-2'
  Type: 'Del'
- Rule:
  Name: 'Header join'
  Header: 'X-Custom-2'
  Value: 'False'
  Type: 'Set'
```

This will firstly set the header `X-Custom-2` to 'True', then delete it and finally set it again but with `False`

# Authors
| Tom Moulard | Clément David | Martin Huvelle | Alexandre Bossut-Lasry |
|-------------|---------------|----------------|------------------------|
|[![](img/gopher-tom_moulard.png)](https://tom.moulard.org)|[![](img/gopher-clement_david.png)](https://github.com/cledavid)|[![](img/gopher-martin_huvelle.png)](https://github.com/nitra-mfs)|[![](img/gopher-alexandre_bossut-lasry.png)](https://www.linkedin.com/in/alexandre-bossut-lasry/)|
