# golangQL

golangQL allows you to use [graphQL](http://graphql.org/) syntactics to select fields from `json`-tagged structs in Golang. 

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
and 
```graphQL
{
  name {
      firstName 
      lastName
    } 
 }
```
requests for the relevant types

## Installation
To install the library, run:
```bash
go get github.com/ngalayko/golangQL
```

## Use
See [example](example) folder
