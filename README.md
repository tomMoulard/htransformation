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

4 types of Rules are possibles:
- Rename
- Set
- Delete
- Join

To choose a Rule you have to fill the `Type` field with either
- 'Rename'
- 'Set'
- 'Del'
- 'Join'

Each Rule can be named with the `Name` field

### Rename

A rule Rename need 2 arguments
- `Header`, the header you want to replace
- `Value`, the new header

### Set

A Set rule will either create or replace the header and value (if it already exist)

A rule Set need 2 arguments
- `Header`, the header you want to create
- `Value`, the value of the new header

### Delete

A rule Delete need 1 arguments
- `Header`, the header you want to delete

### Join

A Join rule will concat the values of the existing header with the new one. If the header doesnt exist, it'll do nothing 

It needs 3 arguments
- `Header`, the header you want to join
- `Values`, a list of values to add to the existing header
- `Sep`, the separator you want to use