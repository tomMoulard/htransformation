# Header transformation plugin for traefik

[![Build Status](https://travis-ci.com/tomMoulard/htransformation.svg?branch=main)](https://travis-ci.com/tomMoulard/htransformation)

This plugin allows changing on the fly header's value of a request.

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

To choose a Rule, you have to fill the `Type` field with either
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

A Set rule will either create or replace the header and value (if it already exists)

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

A Join rule will concatenate the values of the existing header with the new one. If the header doesn't exist, it'll do nothing

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

### RewriteValue Rule

A RewriteValue Rule will replace the values of the headers identified by a matching regex with the provided value.

It needs 2 arguments
- `Header`, the header or regex identifying the headers you want to change
- `Value`, the new value of the headers

```yaml
# Example RewriteValueRule
- Rule:
      Name: 'Header rewriteValue'
      Header: 'Foo'
      Value: 'X-(.*)'
      ValueReplace: 'Y-$1'
      Type: 'RewriteValueRule'
```
```yaml
# Old header:
Foo: X-Test

# Modified header:
Foo: Y-Test
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
Will set the header `X-Custom-2` to 'True', then delete it and set it again but with `False`

# Authors
| Tom Moulard | Clément David | Martin Huvelle | Alexandre Bossut-Lasry |
|-------------|---------------|----------------|------------------------|
|[![](img/gopher-tom_moulard.png)](https://tom.moulard.org)|[![](img/gopher-clement_david.png)](https://github.com/cledavid)|[![](img/gopher-martin_huvelle.png)](https://github.com/nitra-mfs)|[![](img/gopher-alexandre_bossut-lasry.png)](https://www.linkedin.com/in/alexandre-bossut-lasry/)|
