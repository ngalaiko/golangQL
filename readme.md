# golangQL [![Build Status](https://travis-ci.org/ngalayko/golangQL.svg?branch=master)](https://travis-ci.org/ngalayko/golangQL)

golangQL allows you to use [graphQL](http://graphql.org/) syntactics to select fields from `json`-tagged Golang structs.

## Description
```go
func Filter(v interface{}, query string) (interface{}, error)
``` 
returns an `interface{}` value which actually is a `map[string]interface{}`. Keys of this map are `json`-tags and values are values of struct fields. So result of `json.Marshal(v interface{})` will contain only fields that were described in `graphQL query`. Null values are skipped.

See [example](example/main.go).

## Features
All filter functions are cached, for field structs also, so if you request 
```graphQL
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
```graphQL
{
  nephews { 
      name { 
        firstName 
      } 
}
```
```graphQL
{
  name {
      firstName 
      lastName
    } 
 }
```
```graphQL
{
  name { 
    firstName 
  } 
}
```
requests for the relevant types.

## Installation
To install the library, run:
```bash
go get github.com/ngalayko/golangQL
```
