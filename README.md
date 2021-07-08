# Header transformation plugin for traefik

[![Build Status](https://travis-ci.com/tomMoulard/htransformation.svg?branch=main)](https://travis-ci.com/tomMoulard/htransformation)

This plugin allow to change on the fly header's value of a request.

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

To choose a Rule you have to fill the `Type` field with either

- 'Rename'  : to rename a header
- 'Set'     : to Set a header
- 'Del'     : to Delete a header
- 'Join'    : to Join values on a header

Each Rule can be named with the `Name` field

### Rename

A Rename rule need 2 arguments

- `Header`, the regex of the header you want to replace
- `Value`, the new header

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

A Set rule will either create or replace the header and value (if it already exist)

A rule Set need 2 arguments

- `Header`, the header you want to create
- `Value`, the value of the new header

```yaml
# Example Join
- Rule:
      Name: 'Set Cache-Control'
      Header: 'Cache-Control'
      Value: 'Foo'
      Type: 'Join'
```

```yaml
# New header:
Cache-Control: Foo
```

### Delete

A rule Delete need 1 arguments

- `Header`, the header you want to delete

```yaml
# Example Del
- Rule:
      Name: 'Delete Cache-Control'
      Header: 'Cache-Control'
      Type: 'Del'
```


### Join

A Join rule will concat the values of the existing header with the new one. If the header doesnt exist, it'll do nothing 

It needs 3 arguments
- `Header`, the header you want to join
- `Values`, a list of values to add to the existing header
- `Sep`, the separator you want to use

```yaml
# Example Join
- Rule:
      Name: 'Header join'
      Header: 'Cache-Control'
      Sep: ','
      Values:
        - 'Foo'
        - 'Bar'
      Type: 'Join'
```

```yaml
# Old header:
Cache-Control: gzip, deflate

# Joined header:
Cache-Control: gzip, deflate,Foo,Bar
```

### Careful

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

### Advanced: (Re)Using other headers

You can reuse other header values in `Value` or one of the `Values` by setting an additional argument `HeaderPrefix`.
Example:

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

Will firstly set the header `X-Custom-2` to 'True', then delete it and lastly set it again but with `False`

# Authors
| Tom Moulard | Clément David | Martin Huvelle | Alexandre Bossut-Lasry |
|-------------|---------------|----------------|------------------------|
|[![](img/gopher-tom_moulard.png)](https://tom.moulard.org)|[![](img/gopher-clement_david.png)](https://github.com/cledavid)|[![](img/gopher-martin_huvelle.png)](https://github.com/nitra-mfs)|[![](img/gopher-alexandre_bossut-lasry.png)](https://www.linkedin.com/in/alexandre-bossut-lasry/)|
