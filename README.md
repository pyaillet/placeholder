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

$ cat values.json
{
  "KEY": "value from json"
}

$ placeholder -s \${ -e \} replace -i values.json template.file

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
$ placeholder -s \${ -e \} replace -i values.yaml template.file
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
$ placeholder -s \${ -e \} replace -i values.properties template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from properties
```

### From environment

```shell
$ KEY="value from env" placeholder -s \${ -e \} replace template.file
```
```shell
$ cat template.file
This is a simple template file containing a <KEY> placeholder surrounded with
'${' and '}' which will be replaced by the value contained in a value file.

The key value is value from env
```
