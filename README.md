# 📜 go-problemdetails

[![License](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![Build status](https://github.com/rgbrota/go-problemdetails/actions/workflows/ci.yml/badge.svg)](https://github.com/rgbrota/go-problemdetails/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/report/github.com/rgbrota/go-problemdetails)](https://goreportcard.com/report/github.com/rgbrota/go-problemdetails)

Problem details specification [RFC-7807] implementation library written in Go. 

The objective of problem details is to provide a way to carry machine-readable details of errors in a HTTP response to avoid the need to define new error response formats for HTTP APIs. For more information see [RFC-7807](https://www.rfc-editor.org/rfc/rfc7807).

## Installation

```go get github.com/rgbrota/go-problemdetails```

## Usage

This repository contains a struct named ```ProblemDetails```, which should be used as a response in HTTP APIs when dealing with errors. It is prepared to be used by the JSON (encoding/json) and XML (encoding/xml) marshaller provided by the standard library. 

A couple of things you will have to pay attention to in order to follow the specification when using ```ProblemDetails``` in your APIs:
- Set the proper ```content-type``` header value (```application/problem+json``` or ```application/problem+xml```) in your responses
- The HTTP response status code has to match the status code used in the ```ProblemDetails``` struct

For more information on how to be compliant with the specification, please see [RFC-7807](https://www.rfc-editor.org/rfc/rfc7807).

There are two ways to create a new instance: from scratch providing information for all the fields and from a HTTP status code.

### Creating a new ProblemDetails from a status code

The function ```FromHTTPStatus``` is the easiest way to create a new instance of the struct. It only requires the HTTP status code, the rest of the fields will have a default value, and some fields will not be marshalled.

```
pd := FromHTTPStatus(http.StatusInternalServerError)
```

This will create a new ```ProblemDetails``` struct that would be marshalled to the following structure:

```json
{
  "type": "about:blank",
  "title": "Internal Server Error",
  "status": 500
}
```

```xml
<problem xmlns="urn:ietf:rfc:7807">
    <type>about:blank</type>
    <title>Internal Server Error</title>
    <status>500</status>
</problem>
```

It doesn't provide much context, but depending on the error it can be a convenient solution.

### Creating a new ProblemDetails from scratch

The function ```New``` allows you to create an instance by passing all the fields. By using this method of creating a ```ProblemDetails``` we can provide much more context about what happened.

```
pd := New("https://some-domain.com/not_found", "Not Found", http.StatusNotFound, "The object with id 5 was not found", "https://some-domain.com/objects/5")
```

This will create a new ```ProblemDetails``` struct that would be marshalled to the following structure:

```json
{
  "type": "https://some-domain.com/not_found",
  "title": "Not Found",
  "status": 404,
  "detail": "The object with id 5 was not found",
  "instance": "https://some-domain.com/objects/5"
}
```

```xml
<problem xmlns="urn:ietf:rfc:7807">
    <type>https://some-domain.com/not_found</type>
    <title>Not Found</title>
    <status>404</status>
    <detail>The object with id 5 was not found</detail>
    <instance>https://some-domain.com/objects/5</instance>
</problem>
```

### Extension fields

The specification includes the possibility to add custom fields to the ```ProblemDetails``` definition which are known as extension fields and can be used to provide further context about the error. 

For example, we can take look at how ```ProblemDetails``` is used to provide validation errors in the ASP.NET world [in the following documentation](https://learn.microsoft.com/en-us/dotnet/api/microsoft.aspnetcore.mvc.validationproblemdetails?view=aspnetcore-6.0). The gist is that an ```Errors``` dictionary is added to the base struct definition as an extension field to provide error context by using a ```map[string][]string``` where the key is the field name and the value is a list of errors associated with it. 

The implementation of something similar would look like the following:

```go
type ValidationProblemDetails struct {
  ProblemDetails

  Errors map[string][]string
}
```

Having extended the base struct with the new errors map, an example ```ValidationProblemDetails``` would look like the following when marshalled to JSON:

```json
{
  "type": "https://some-domain.com/bad_request",
  "title": "Bad Request",
  "status": 400,
  "detail": "The were some validation errors while creating the object",
  "instance": "https://some-domain.com/objects",
  "errors": {
    "field1": [
      "The maximum length is 8 characters",
      "Does not meet the charset criteria"
    ],
    "field2": [
      "Date cannot be higher than current date"
    ]
  }
}
```
