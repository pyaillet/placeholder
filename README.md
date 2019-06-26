# placeholder

[![Build Status](https://travis-ci.org/pyaillet/placeholder.svg?branch=master)](https://travis-ci.org/pyaillet/placeholder)
[![codecov](https://codecov.io/gh/pyaillet/placeholder/branch/master/graph/badge.svg)](https://codecov.io/gh/pyaillet/placeholder)

Minimal project used to handle placeholders in files

## Description

The goal of this project is to provide a simple tool to replace placeholders
in files with values either from json, yaml, property files or from the
environment.
A placeholder is made from a prefix, a suffix and a key identifier matching
this regex: `[A-Za-z][A-Za-z0-9_]*`.

It's also possible to list keys identified in one or more files by using the
command: `placeholder -s '{{' -e '}}' list *.templates`.

When a key is identified but no corresponding value is provided the program
exits with an error.

## Usage examples

### With json

```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is ${KEY}
```
```shell
$ cat values.json
{
  "KEY": "value from json"
}
```
```shell
$ placeholder replace -i values.json template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from json
```

### With yaml

```shell
$ cat values.yaml
KEY: value from yaml
```
```shell
$ placeholder replace -i values.yaml template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from yaml
```

### With properties

```shell
$ cat values.properties
KEY=value from properties
```
```shell
$ placeholder replace -i values.properties template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from properties
```

### From environment

```shell
$ KEY="value from env" placeholder replace template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from env
```

### Bloc separator
You can also change the default bloc separator by specifying the start and end bloc with the `-s` and `-e` arguments :

```shell
$ placeholder -s '%#' -e '#%' replace -i values.json template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'%#' and '#%' which will be replaced by the value contained in a value file.

The key value is value from env
```

## Build it

To build this project, you must have make and go >= 1.12 installed.
You can then just type:
`make build`

If you don't want to install go but use docker instead, type:
`make docker-build`
