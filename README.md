# Header transformation plugin for traefik

This plugin allows changing on the fly, the header value of a request.

## How to dev

```bash
$ docker compose up
```

## How to use

To choose a Rule you have to fill the `Type` field with one of the following:

- 'Del'             : to Delete a header
- 'Join'            : to Join values on a header
- 'Rename'          : to rename a header
- 'RewriteValueRule': to rewrite header values
- 'Set'             : to Set a header

Each Rule can be named with the `Name` field.

Each Rule can also be configured to change headers on the request or the
response by using the `SetOnResponse` configuration.
If `SetOnResponse` is set to `true`, the header will be changed on the response.
Otherwise, it will be changed on the request.
Its default value is `false`.

### Rename

A Rule Rename needs two arguments.

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
# Example 
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

### Delete

A rule Delete need one arguments

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

You can reuse other header values in `Value` or one of the `Values` by setting an additional argument `HeaderPrefix`.
Example:

```yaml
# Example Usage
- Rule:
  Name: 'Header set'
  Header: 'X-Forwarded-For'
  HeaderPrefix: "^"
  Sep: ','
  Values:
      - 'Foo'
      - '^CF-Connecting-IP'
  Type: 'Join'
```

```yaml
# Old header:
X-Forwarded-For: 1.1.1.1
CF-Connecting-IP: 2.2.2.2
# New headers:
X-Forwarded-For: 1.1.1.1,Foo,2.2.2.2
CF-Connecting-IP: 2.2.2.2
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
