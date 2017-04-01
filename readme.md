# golangQL

golangQL allows you to use [graphQL](http://graphql.org/) syntactics to select fields from `json`-tagged structs in Golang. 

## Features
All filter functions are cached, for field structs also, so if you request 
```json
{
  name {
    firstName 
    lastName
  } 
  nephews { 
    name { 
      firstName 
    } 
  hat 
} 
```
it also will cache
```json
nephews { 
    name { 
      firstName 
    } 
```
and 
```
name {
    firstName 
    lastName
  } 
```

## Installation
To install the library, run:
```bash
go get github.com/ngalayko/golangQL
```

## Use
See [example](example) folder

## Third Party Libraries
| Name          | Author        | Description  |
|:-------------:|:-------------:|:------------:|
| [testify](github.com/stretchr/testify/assert) | [Stretchr, Inc. ](https://github.com/sogko) | A sacred extension to the standard go testing package
