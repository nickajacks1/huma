![Huma Rest API Framework](https://user-images.githubusercontent.com/106826/78105564-51102780-73a6-11ea-99ff-84d6c1b3e8df.png)

[![HUMA Powered](https://img.shields.io/badge/Powered%20By-HUMA-f40273)](https://huma.rocks/) [![CI](https://github.com/danielgtaylor/huma/workflows/CI/badge.svg?branch=master)](https://github.com/danielgtaylor/huma/actions?query=workflow%3ACI+branch%3Amaster++) [![codecov](https://codecov.io/gh/danielgtaylor/huma/branch/master/graph/badge.svg)](https://codecov.io/gh/danielgtaylor/huma) [![Docs](https://godoc.org/github.com/danielgtaylor/huma?status.svg)](https://pkg.go.dev/github.com/danielgtaylor/huma?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/danielgtaylor/huma)](https://goreportcard.com/report/github.com/danielgtaylor/huma)

- [What is huma?](#intro)
- [Example](#example)
- [Install](#install)
- [Documentation](#documentation)
- [Design](#design)

<a name="intro"></a>
A modern, simple, fast & flexible micro framework for building REST/RPC APIs in Go backed by OpenAPI 3 and JSON Schema. Pronounced IPA: [/'hjuːmɑ/](https://en.wiktionary.org/wiki/Wiktionary:International_Phonetic_Alphabet). The goals of this project are to provide:

- Incremental adoption for teams with existing services
  - Bring your own router, middleware, and logging/metrics
  - Extensible OpenAPI & JSON Schema layer to document existing routes
- A modern REST or RPC API backend framework for Go developers
  - Described by [OpenAPI 3.1](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md) & [JSON Schema](https://json-schema.org/)
- Guard rails to prevent common mistakes
- Documentation that can't get out of date
- High-quality generated developer tooling

Features include:

- Declarative interface on top of your router of choice:
  - Operation & model documentation
  - Request params (path, query, or header)
  - Request body
  - Responses (including errors)
  - Response headers
- JSON Errors using [RFC7807](https://tools.ietf.org/html/rfc7807) and `application/problem+json` by default (but can be changed)
- Per-operation request size limits with sane defaults
- [Content negotiation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation) between server and client
  - Support for JSON ([RFC 8259](https://tools.ietf.org/html/rfc8259)) and CBOR ([RFC 7049](https://tools.ietf.org/html/rfc7049)) content types via the `Accept` header with the default config.
- Conditional requests support, e.g. `If-Match` or `If-Unmodified-Since` header utilities.
- Optional automatic generation of `PATCH` operations that support:
  - [RFC 7386](https://www.rfc-editor.org/rfc/rfc7386) JSON Merge Patch
  - [RFC 6902](https://www.rfc-editor.org/rfc/rfc6902) JSON Patch
  - [Shorthand](https://github.com/danielgtaylor/shorthand)
- Annotated Go types for input and output models
  - Generates JSON Schema from Go types
  - Automatic input model validation & error handling
- Documentation generation using [Stoplight Elements](https://stoplight.io/open-source/elements)
- Optional CLI built-in, configured via arguments or environment variables
  - Set via e.g. `-p 8000`, `--port=8000`, or `SERVICE_PORT=8000`
  - Startup actions & graceful shutdown built-in
- Generates OpenAPI for access to a rich ecosystem of tools
  - Mocks with [API Sprout](https://github.com/danielgtaylor/apisprout)
  - SDKs with [OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator)
  - CLI with [Restish](https://rest.sh/)
  - And [plenty](https://openapi.tools/) [more](https://apis.guru/awesome-openapi3/category.html)
- Generates JSON Schema for each resource using optional `describedby` link relation headers as well as optional `$schema` properties in returned objects that integrate into editors for validation & completion.

This project was inspired by [FastAPI](https://fastapi.tiangolo.com/). Logo & branding designed by Kari Taylor.

# Example

Here is a complete basic hello world example in Huma, that shows how to initialize a Huma app complete with CLI, declare a resource with an operation, and define its handler function.

```go
package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/adapters/humachi"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" default:"8888"`
}

// GreetingInput represents the greeting operation request.
type GreetingInput struct {
	Name string `path:"name"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func main() {
	// Create a CLI app which takes a port option.
	huma.NewCLI(func(cli huma.CLI, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Register GET /greeting/{name}
		huma.Register(api, huma.Operation{
			OperationID: "get-greeting",
			Method:      http.MethodGet,
			Path:        "/greeting/{name}",
		}, func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message: fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		// Tell the CLI how to start your router.
		cli.OnStart(func() {
			router.Listen(fmt.Sprintf(":%d", options.Port)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	app.Run()
}
```

You can test it with `go run hello.go` (optionally pass `--port` to change the default) and make a sample request using [Restish](https://rest.sh/) (or `curl`):

```sh
# Get the message from the server
$ restish :8888/greeting/world
HTTP/1.1 200 OK
...
{
	$schema: "http://localhost:8888/schemas/GreetingOutputBody.json",
	message: "Hello, world!"
}
```

Even though the example is tiny you can also see some generated documentation at http://localhost:8888/docs. The generated OpenAPI is available at http://localhost:8888/openapi.json or http://localhost:8888/openapi.yaml.

See the examples directory for more complete examples.

- [Minimal](./examples/minimal/minimal.go) (a minimal "hello world")
- [Echo](./examples/echo/echo.go) (echo input back to the user with validation)
- [Notes](./examples/notes/notes.go) (note-taking API)
- [Timeout](./examples/timeout/timeout.go) (show third-party request timing out)
- [Test](./examples/test/service.go) (how to write a test)

# Install

```sh
# After: go mod init ...
go get -u github.com/danielgtaylor/huma/v2
```

# Documentation

Official Go package documentation can always be found at https://pkg.go.dev/github.com/danielgtaylor/huma. Below is an introduction to the various features available in Huma.

> :whale: Hi there! I'm the happy Huma whale here to provide help. You'll see me leave helpful tips down below.

## BYOR (Bring Your Own Router)

Huma is designed to be router-agnostic to enable incremental adoption in existing and new services across a large number of organizations. This means you can use any router you want, or even write your own. The only requirement is that you implement a small `huma.Adapter` interface. This is how Huma integrates with your router.

Adapters are in the `adapters` directory and named after the router they support. Many common routers are supported out of the box:

- [chi](https://github.com/go-chi/chi) via `humachi`
- [gin](https://gin-gonic.com/) via `humagin`
- [Fiber](https://gofiber.io/) via `humafiber`

Writing your own wrapper is quick and simple, and PRs are accepted for additional adapters to be built-in.

Adapters are instantiated by wrapping your router and providing a Huma configuration object which describes the API. Here is a simple example using Fiber:

```go
// Create your router.
app := fiber.New()

// Wrap the router with Huma to create an API instance.
api := humafiber.New(app, huma.DefaultConfig("My API", "1.0.0"))

// Register your operations with the API.
// ...

// Start the server!
app.Listen(":3000")
```

## Open API Generation & Extensibility

Huma generates Open API 3.1.0 compatible JSON/YAML specs and provides rendered documentation automatically. Every operation that is registered with the API is included in the spec by default. The operation's inputs and outputs are used to generate the request and response parameters / schemas.

Sometimes you may need to customize the generated Open API spec. With Huma v2 you have full access and can modify it as needed. For example, to set up a security scheme:

```go
config := huma.DefaultConfig("My API", "1.0.0")
config.Components.SecuritySchemes.....
api := humachi.New(router, config)
```

---

# Old docs below!

These will get updated soon.

---

The Huma router is the entrypoint to your service or application. There are a couple of ways to create it, depending on what level of customization you need.

```go
// Simplest way to get started, which creats a router and a CLI with default
// middleware attached. Note that the CLI is a router.
app := cli.NewRouter("API Name", "1.0.0")

// Doing the same as above by hand:
router := huma.New("API Name", "1.0.0")
app := cli.New(router)
middleware.Defaults(app)

// Start the CLI after adding routes:
app.Run()
```

You can also skip using the built-in `cli` package:

```go
// Create and start a new router by hand:
router := huma.New("API Name", "1.0.0")
router.Middleware(middleware.DefaultChain)
router.Listen("127.0.0.1:8888")
```

## Resources

Huma APIs are composed of resources and sub-resources attached to a router. A resource refers to a unique URI on which operations can be performed. Huma resources can have middleware attached to them, which run before operation handlers.

```go
// Create a resource at a given path.
notes := app.Resource("/notes")

// Add a middleware to all operations under `/notes`.
notes.Middleware(MyMiddleware())

// Create another resource that includes a path parameter: /notes/{id}
// Paths look like URI templates and use wrap parameters in curly braces.
note := notes.SubResource("/{id}")

// Create a sub-resource at /notes/{id}/likes.
sub := note.SubResource("/likes")
```

> :whale: Resources should be nouns, and plural if they return more than one item. Good examples: `/notes`, `/likes`, `/users`, `/videos`, etc.

## Operations

Operations perform an action on a resource using an HTTP method verb. The following verbs are available:

- Head
- Get
- Post
- Put
- Patch
- Delete
- Options

Operations can take inputs in the form of path, query, and header parameters and/or request bodies. They must declare what response status codes, content types, and structures they return.

Every operation has a handler function and takes at least a `huma.Context`, described in further detail below:

```go
app.Resource("/op").Get("get-op", "Example operation",
	// Response declaration goes here!
).Run(func (ctx huma.Context) {
	// Handler implementation goes here!
})
```

> :whale: Operations map an HTTP action verb to a resource. You might `POST` a new note or `GET` a user. Sometimes the mapping is less obvious and you can consider using a sub-resource. For example, rather than unliking a post, maybe you `DELETE` the `/posts/{id}/likes` resource.

## Context

As seen above, every handler function gets at least a `huma.Context`, which combines an `http.ResponseWriter` for creating responses, a `context.Context` for cancellation/timeouts, and some convenience functions. Any library that can use either of these interfaces will work with a Huma context object. Some examples:

```go
// Calling third-party libraries that might take too long
results := mydb.Fetch(ctx, "some query")

// Write an HTTP response
ctx.Header().Set("Content-Type", "text/plain")
ctx.WriteHeader(http.StatusNotFound)
ctx.Write([]byte("Could not find foo"))
```

> :whale: Since you can write data to the response multiple times, the context also supports streaming responses. Just remember to set (or remove) the timeout.

## Responses

In order to keep the documentation & service specification up to date with the code, you **must** declare the responses that your handler may return. This includes declaring the content type, any headers it might return, and what model it returns (if any). The `responses` package helps with declaring well-known responses with the right code/docs/model and corresponds to the statuses in the `http` package, e.g. `resposes.OK()` will create a response with the `http.StatusOK` status code.

```go
// Response structures are just normal Go structs
type Thing struct {
	Name string `json:"name"`
}

// ... initialization code goes here ...

things := app.Resource("/things")
things.Get("list-things", "Get a list of things",
	// Declare a successful response that returns a slice of things
	responses.OK().Headers("Foo").Model([]Thing{}),
	// Errors automatically set the right status, content type, and model for you.
	responses.InternalServerError(),
).Run(func(ctx huma.Context) {
	// This works because the `Foo` header was declared above.
	ctx.Header().Set("Foo", "Some value")

	// The `WriteModel` convenience method handles content negotiation and
	// serializaing the response for you.
	ctx.WriteModel(http.StatusOK, []Thing{
		Thing{Name: "Test1"},
		Thing{Name: "Test2"},
	})

	// Alternatively, you can write an error
	ctx.WriteError(http.StatusInternalServerError, "Some message")
})
```

If you try to set a response status code or header that was not declared you will get a runtime error. If you try to call `WriteModel` or `WriteError` more than once then you will get an error because the writer is considered closed after those methods.

### Errors

Errors use [RFC 7807](https://tools.ietf.org/html/rfc7807) and return a structure that looks like:

```json
{
  "status": 504,
  "title": "Gateway Timeout",
  "detail": "Problem with HTTP request",
  "errors": [
    {
      "message": "Get \"https://httpstat.us/418?sleep=5000\": context deadline exceeded"
    }
  ]
}
```

The `errors` field is optional and may contain more details about which specific errors occurred.

It is recommended to return exhaustive errors whenever possible to prevent user frustration with having to keep retrying a bad request and getting back a different error. The context has `AddError` and `HasError()` functions for this:

```go
app.Resource("/exhaustive").Get("exhaustive", "Exhastive errors example",
	responses.OK(),
	responses.BadRequest(),
).Run(func(ctx huma.Context) {
	for i := 0; i < 5; i++ {
		// Use AddError to add multiple error details to the response.
		ctx.AddError(fmt.Errorf("Error %d", i))
	}

	// Check if the context has had any errors added yet.
	if ctx.HasError() {
		// Use WriteError to set the actual status code, top-level message, and
		// any additional errors. This sends the response.
		ctx.WriteError(http.StatusBadRequest, "Bad input")
		return
	}
})
```

While every attempt is made to return exhaustive errors within Huma, each individual response can only contain a single HTTP status code. The following chart describes which codes get returned and when:

```mermaid
flowchart TD
	Request[Request has errors?] -->|yes| Panic
	Request -->|no| Continue[Continue to handler]
	Panic[Panic?] -->|yes| 500
	Panic -->|no| RequestBody[Request body too large?]
	RequestBody -->|yes| 413
	RequestBody -->|no| RequestTimeout[Request took too long to read?]
	RequestTimeout -->|yes| 408
	RequestTimeout -->|no| ParseFailure[Cannot parse input?]
	ParseFailure -->|yes| 400
	ParseFailure -->|no| ValidationFailure[Validation failed?]
	ValidationFailure -->|yes| 422
	ValidationFailure -->|no| 400
```

This means it is possible to, for example, get an HTTP `408 Request Timeout` response that _also_ contains an error detail with a validation error for one of the input headers. Since request timeout has higher priority, that will be the response status code that is returned.

## WriteContent

Write contents allows you to write content in the provided ReadSeeker in the response. It will handle Range, If-Match, If-Unmodified-Since, If-None-Match, If-Modified-Since, and if-Range requests for caching and large object partial content responses.

example of a partial content response sending a large video file:

```go
package main

import (
	"net/http"
	"os"

	"github.com/danielgtaylor/huma"
	"github.com/danielgtaylor/huma/cli"
	"github.com/danielgtaylor/huma/responses"
)

func main() {
	// Create a new router & CLI with default middleware.
	app := cli.NewRouter("Video Content", "1.0.0")

	// Declare the /vid resource and a GET operation on it.
	app.Resource("/vid").Get("get-vid", "Get video content",
		// responses.WriteContent returns all the responses needed for context.WriteContent
		responses.WriteContent()...,
	).Run(func(ctx huma.Context) {
		// This is he handler function for the operation. Write the response.
		ctx.Header().Set("Content-Type", "video/mp4")
		f, err := os.Open("./vid.mp4")
		if err != nil {
			ctx.WriteError(http.StatusInternalServerError, "Error while opening video")
			return
		}
		defer f.Close()

		fStat, err := os.Stat("./vid.mp4")
		if err != nil {
			ctx.WriteError(http.StatusInternalServerError, "Error while attempting to get Last Modified time")
			return
		}

		// Note that name (ex: vid.mp4) and lastModified (ex: fState.ModTime()) are both optional
		// if name is "" Content-Type must be set by the handler
		// if lastModified is time.Time{} Modified request headers will not be respected by WriteContent
		ctx.WriteContent("vid.mp4", f, fStat.ModTime())
	})

	// Run the CLI. When passed no arguments, it starts the server.
	app.Run()
}
```

Note that `WriteContent` does not automatically set the mime type. You should set the `Content-Type` response header directly. Also in order for `WriteContent` to respect the `Modified` headers you must call `SetContentLastModified`. This is optional and if not set `WriteContent` will simply not respect the `Modified` request headers.

## Request Inputs

Requests can have parameters and/or a body as input to the handler function. Like responses, inputs use standard Go structs but the tags are different. Here are the available tags:

| Tag      | Description                        | Example                  |
| -------- | ---------------------------------- | ------------------------ |
| `path`   | Name of the path parameter         | `path:"thing-id"`        |
| `query`  | Name of the query string parameter | `query:"q"`              |
| `header` | Name of the header parameter       | `header:"Authorization"` |

The following types are supported out of the box:

| Type                | Example Inputs         |
| ------------------- | ---------------------- |
| `bool`              | `true`, `false`        |
| `[u]int[16/32/64]`  | `1234`, `5`, `-1`      |
| `float32/64`        | `1.234`, `1.0`         |
| `string`            | `hello`, `t`           |
| `time.Time`         | `2020-01-01T12:00:00Z` |
| slice, e.g. `[]int` | `1,2,3`, `tag1,tag2`   |

For example, if the parameter is a query param and the type is `[]string` it might look like `?tags=tag1,tag2` in the URI.

The special struct field `Body` will be treated as the input request body and can refer to another struct or you can embed a struct inline. `RawBody` can also be used to provide access to the `[]byte` used to validate & load `Body`.

Here is an example:

```go
type MyInputBody struct {
	Name string `json:"name"`
}

type MyInput struct {
	ThingID     string      `path:"thing-id" doc:"Example path parameter"`
	QueryParam  int         `query:"q" doc:"Example query string parameter"`
	HeaderParam string      `header:"Foo" doc:"Example header parameter"`
	Body        MyInputBody `doc:"Example request body"`
}

// ... Later you use the inputs

// Declare a resource with a path parameter that matches the input struct. This
// is needed because path parameter positions matter in the URL.
thing := app.Resource("/things/{thing-id}")

// Next, declare the handler with an input argument.
thing.Get("get-thing", "Get a single thing",
	responses.NoContent(),
).Run(func(ctx huma.Context, input MyInput) {
	fmt.Printf("Thing ID: %s\n", input.ThingID)
	fmt.Printf("Query param: %s\n", input.QueryParam)
	fmt.Printf("Header param: %s\n", input.HeaderParam)
	fmt.Printf("Body name: %s\n", input.Body.Name)
})
```

Try a request against the service like:

```sh
# Restish example
$ restish :8888/things/abc123?q=3 -H "Foo: bar" name: Kari
```

### Multiple Request Bodies

Request input structs can support multiple body types based on the content type of the request, with an unknown content type defaulting to the first-defined body. This can be used for things like versioned inputs or to support wildly different input types (e.g. JSON Merge Patch vs. JSON Patch). Example:

```go
type MyInput struct {
	BodyV2 *MyInputBodyV1 `body:"application/my-type-v2+json"`
	BodyV1 *MyInputBodyV1 `body:"application/my-type-v1+json"`
}
```

It's your responsibility to check which one is non-`nil` in the operation handler. If not using pointers, you'll need to check a known field to determine which was actually sent by the client.

### Parameter & Body Validation

All supported JSON Schema tags work for parameters and body fields. Validation happens before the request handler is called, and if needed an error response is returned. For example:

```go
type MyInput struct {
	ThingID    string `path:"thing-id" pattern:"^th-[0-9a-z]+$" doc:"..."`
	QueryParam int    `query:"q" minimum:"1" doc:"..."`
}
```

See "Validation" for more info.

### Input Composition

Because inputs are just Go structs, they are composable and reusable. For example:

```go
type AuthParam struct {
	Authorization string `header:"Authorization"`
}

type PaginationParams struct {
	Cursor string `query:"cursor"`
	Limit  int    `query:"limit"`
}

// ... Later in the code
app.Resource("/things").Get("list-things", "List things",
	responses.NoContent(),
).Run(func (ctx huma.Context, input struct {
	AuthParam
	PaginationParams
}) {
	fmt.Printf("Auth: %s, Cursor: %s, Limit: %d\n", input.Authorization, input.Cursor, input.Limit)
})
```

### Input Streaming

It's possible to support input body streaming for large inputs by declaring your body as an `io.Reader`:

```go
type StreamingBody struct {
	Body io.Reader
}
```

You probably want to combine this with custom timeouts, or removing them altogether.

```go
op := app.Resource("/streaming").Post("post-stream", "Write streamed data",
	responses.NoContent(),
)
op.NoBodyReadTimeout()
op.Run(...)
```

If you just need access to the input body bytes and still want to use the built-in JSON Schema validation, then you can instead use the `RawBody` input struct field.

```go
type MyBody struct {
	// This will generate JSON Schema, validate the input, and parse it.
	Body MyStruct

	// This will contain the raw bytes used to load the above.
	RawBody []byte
}
```

### Resolvers

Sometimes the built-in validation isn't sufficient for your use-case, or you want to do something more complex with the incoming request object. This is where resolvers come in.

Any input struct can be a resolver by implementing the `huma.Resolver` interface, including embedded structs. Each resolver takes the current context and the incoming request. For example:

```go
// MyInput demonstrates inputs/transformation
type MyInput struct {
	Host   string
	Name string `query:"name"`
}

func (m *MyInput) Resolve(ctx huma.Context, r *http.Request) {
	// Get request info you don't normally have access to.
	m.Host = r.Host

	// Transformations or other data validation
	m.Name = strings.Title(m.Name)
}

// Then use it like any other input struct:
app.Resource("/things").Get("list-things", "Get a filtered list of things",
	responses.NoContent(),
).Run(func(ctx huma.Context, input MyInput) {
	fmt.Printf("Host: %s\n", input.Host)
	fmt.Printf("Name: %s\n", input.Name)
})
```

It is recommended that you do not save the request. Whenever possible, use existing mechanisms for describing your input so that it becomes part of the OpenAPI description.

#### Resolver Errors

Resolvers can set errors as needed and Huma will automatically return a 400-level error response before calling your handler. This makes resolvers a good place to run additional complex validation steps so you can provide the user with a set of exhaustive errors.

```go
type MyInput struct {
	Host   string
}

func (m *MyInput) Resolve(ctx huma.Context, r *http.Request) {
	if m.Host = r.Hostname; m.Host == "localhost" {
		ctx.AddError(&huma.ErrorDetail{
			Message: "Invalid value!",
			Location: "request.host",
			Value: m.Host,
		})
	}
}
```

### Conditional Requests

There are built-in utilities for handling [conditional requests](https://developer.mozilla.org/en-US/docs/Web/HTTP/Conditional_requests), which serve two broad purposes:

1. Sparing bandwidth on reading a document that has not changed, i.e. "only send if the version is different from what I already have"
2. Preventing multiple writers from clobbering each other's changes, i.e. "only save if the version on the server matches what I saw last"

Adding support for handling conditional requests requires four steps:

1. Import the `github.com/danielgtaylor/huma/conditional` package.
2. Add the response definition (`304 Not Modified` for reads or `412 Precondition Failed` for writes)
3. Add `conditional.Params` to your input struct.
4. Check if conditional params were passed and handle them. The `HasConditionalParams()` and `PreconditionFailed(...)` methods can help with this.

Implementing a conditional read might look like:

```go
app.Resource("/resource").Get("get-resource", "Get a resource",
	responses.OK(),
	responses.NotModified(),
).Run(func(ctx huma.Context, input struct {
	conditional.Params
}) {
	if input.HasConditionalParams() {
		// TODO: Get the ETag and last modified time from the resource.
		etag := ""
		modified := time.Time{}

		// If preconditions fail, abort the request processing. Response status
		// codes are already set for you, but you can optionally provide a body.
		// Returns an HTTP 304 not modified.
		if input.PreconditionFailed(ctx, etag, modified) {
			return
		}
	}

	// Otherwise do the normal request processing here...
	// ...
})
```

Similarly a write operation may look like:

```go
app.Resource("/resource").Put("put-resource", "Put a resource",
	responses.OK(),
	responses.PreconditionFailed(),
).Run(func(ctx huma.Context, input struct {
	conditional.Params
}) {
	if input.HasConditionalParams() {
		// TODO: Get the ETag and last modified time from the resource.
		etag := ""
		modified := time.Time{}

		// If preconditions fail, abort the request processing. Response status and
		// errors have already been set. Returns an HTTP 412 Precondition Failed.
		if input.PreconditionFailed(ctx, etag, modified) {
			return
		}
	}

	// Otherwise do the normal request processing here...
	// ...
})
```

### Automatic PATCH Support

If a `GET` and a `PUT` exist for the same resource, but no `PATCH` exists at server start up, then by default a `PATCH` operation will be generated for you to make editing more convenient for clients. This behavior can be disabled via `app.DisableAutoPatch()`.

If the `GET` returns an `ETag` or `Last-Modified` header, then these will be used to make conditional requests on the `PUT` operation to prevent distributed write conflicts that might otherwise overwrite someone else's changes.

If the `PATCH` request has no `Content-Type` header, or uses `application/json` or a variant thereof, then JSON Merge Patch is assumed.

## Validation

Go struct tags are used to annotate inputs/output structs with information that gets turned into [JSON Schema](https://json-schema.org/) for documentation and validation.

The standard `json` tag is supported and can be used to rename a field and mark fields as optional using `omitempty`. The following additional tags are supported on model fields:

| Tag                | Description                               | Example                  |
| ------------------ | ----------------------------------------- | ------------------------ |
| `doc`              | Describe the field                        | `doc:"Who to greet"`     |
| `format`           | Format hint for the field                 | `format:"date-time"`     |
| `enum`             | A comma-separated list of possible values | `enum:"one,two,three"`   |
| `default`          | Default value                             | `default:"123"`          |
| `minimum`          | Minimum (inclusive)                       | `minimum:"1"`            |
| `exclusiveMinimum` | Minimum (exclusive)                       | `exclusiveMinimum:"0"`   |
| `maximum`          | Maximum (inclusive)                       | `maximum:"255"`          |
| `exclusiveMaximum` | Maximum (exclusive)                       | `exclusiveMaximum:"100"` |
| `multipleOf`       | Value must be a multiple of this value    | `multipleOf:"2"`         |
| `minLength`        | Minimum string length                     | `minLength:"1"`          |
| `maxLength`        | Maximum string length                     | `maxLength:"80"`         |
| `pattern`          | Regular expression pattern                | `pattern:"[a-z]+"`       |
| `minItems`         | Minimum number of array items             | `minItems:"1"`           |
| `maxItems`         | Maximum number of array items             | `maxItems:"20"`          |
| `uniqueItems`      | Array items must be unique                | `uniqueItems:"true"`     |
| `minProperties`    | Minimum number of object properties       | `minProperties:"1"`      |
| `maxProperties`    | Maximum number of object properties       | `maxProperties:"20"`     |
| `example`          | Example value                             | `example:"123"`          |
| `nullable`         | Whether `null` can be sent                | `nullable:"false"`       |
| `readOnly`         | Sent in the response only                 | `readOnly:"true"`        |
| `writeOnly`        | Sent in the request only                  | `writeOnly:"true"`       |
| `deprecated`       | This field is deprecated                  | `deprecated:"true"`      |

Parameters have some additional validation tags:

| Tag        | Description                    | Example           |
| ---------- | ------------------------------ | ----------------- |
| `internal` | Internal-only (not documented) | `internal:"true"` |

## Recursion

Recursive inputs or outputs are supported by setting a special field in the schema generator (this is to support backward-compatibility):

```go
	// Set to false to generate references to each struct that are shared between
	// the various operation inputs/outputs, enabling support for recursion.
	schema.GenerateInline = false
```

The default behavior is to generate everything inline. Note that switching between the two styles is likely to break generated SDKs, so plan for a major version bump if/when you make the change.

## Middleware

Standard [Go HTTP middleware](https://justinas.org/writing-http-middleware-in-go) is supported. It can be attached to the main router/app or to individual resources, but **must** be added _before_ operation handlers are added.

```go
// Middleware from some library
app.Middleware(somelibrary.New())

// Custom middleware
app.Middleware(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// Request phase, do whatever you want before next middleware or handler
		// gets called.
		fmt.Println("Request coming in")

		// Call the next middleware/handler
		next.ServeHTTP(w, r)

		// Response phase, after handler has run.
		fmt.Println("Response going out!")
	})
})
```

When using the `cli.NewRouter` convenience method, a set of default middleware is added for you. See `middleware.DefaultChain` for more info.

### Enabling OpenTracing

[OpenTracing](https://opentracing.io/) support is built-in, but you have to tell the global tracer where to send the information, otherwise it acts as a no-op. For example, if you use [DataDog APM](https://www.datadoghq.com/blog/opentracing-datadog-cncf/) and have the agent configured wherever you deploy your service:

```go
import (
	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	t := opentracer.New(tracer.WithAgentAddr("host:port"))
	defer tracer.Stop()

	// Set it as a Global Tracer
	opentracing.SetGlobalTracer(t)

	app := cli.NewRouter("My Cool Service", "1.0.0")
	// register routes here
	app.Run()
}
```

### Timeouts, Deadlines, Cancellation & Limits

Huma provides utilities to prevent long-running handlers and issues with huge request bodies and slow clients with sane defaults out of the box.

#### Context Timeouts

Set timeouts and deadlines on the request context and pass that along to libraries to prevent long-running handlers. For example:

```go
app.Resource("/timeout").Get("timeout", "Timeout example",
	responses.String(http.StatusOK),
	responses.GatewayTimeout(),
).Run(func(ctx huma.Context) {
	// Add a timeout to the context. No request should take longer than 2 seconds
	newCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Create a new request that will take 5 seconds to complete.
	req, _ := http.NewRequestWithContext(
		newCtx, http.MethodGet, "https://httpstat.us/418?sleep=5000", nil)

	// Make the request. This will return with an error because the context
	// deadline of 2 seconds is shorter than the request duration of 5 seconds.
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.WriteError(http.StatusGatewayTimeout, "Problem with HTTP request", err)
		return
	}

	ctx.Write([]byte("success!"))
})
```

#### Request Timeouts

By default, a `ReadHeaderTimeout` of _10 seconds_ and an `IdleTimeout` of _15 seconds_ are set at the server level and apply to every incoming request.

Each operation's individual read timeout defaults to _15 seconds_ and can be changed as needed. This enables large request and response bodies to be sent without fear of timing out, as well as the use of WebSockets, in an opt-in fashion with sane defaults.

When using the built-in model processing and the timeout is triggered, the server sends an error as JSON with a message containing the time waited.

```go
type Input struct {
	ID string `json:"id"`
}

app := cli.NewRouter("My API", "1.0.0")
foo := app.Resource("/foo")

// Limit to 5 seconds
create := foo.Post("create-item", "Create a new item",
	responses.NoContent(),
)
create.BodyReadTimeout(5 * time.Second)
create.Run(func (ctx huma.Context, input Input) {
	// Do something here.
})
```

You can also access the underlying TCP connection and set deadlines manually:

```go
create.Run(func (ctx huma.Context, input struct {
	Body io.Reader
}) {
	// Get the connection.
	conn := huma.GetConn(ctx)

	// Set a new deadline on connection reads.
	conn.SetReadDeadline(time.Now().Add(600 * time.Second))

	// Read all the data from the request.
	data, err := ioutil.ReadAll(input.Body)
	if err != nil {
		// If a timeout occurred, this will be a net.Error with `err.Timeout()`
		// returning true.
		panic(err)
	}

	// Do something with data here...
})
```

> :whale: Use `NoBodyReadTimeout()` to disable the default.

#### Request Body Size Limits

By default each operation has a 1 MiB reqeuest body size limit.

When using the built-in model processing and the timeout is triggered, the server sends an error as JSON with a message containing the maximum body size for this operation.

```go
app := cli.NewRouter("My API", "1.0.0")

create := app.Resource("/foo").Post("create-item", "Create a new item",
	responses.NoContent(),
)
// Limit set to 10 MiB
create.MaxBodyBytes(10 * 1024 * 1024)
create.Run(func (ctx huma.Context, input Input) {
	// Body is guaranteed to be 10MiB or less here.
})
```

> :whale: Use `NoMaxBodyBytes()` to disable the default.

## Logging

Huma provides a Zap-based contextual structured logger as part of the default middleware stack. You can access it via the `middleware.GetLogger(ctx)` which returns a `*zap.SugaredLogger`. It requires the use of the `middleware.Logger`, which is included by default when using either `cli.NewRouter` or `middleware.Defaults`.

```go
app := cli.NewRouter("Logging Example", "1.0.0")

app.Resource("/log").Get("log", "Log example",
	responses.NoContent(),
).Run(func (ctx huma.Context) {
	logger := middleware.GetLogger(ctx)
	logger.Info("Hello, world!")
})
```

Manual setup:

```go
router := huma.New("Loggin Example", "1.0.0")
app := cli.New(router)

app.Middleware(middleware.Logger)
middleware.AddLoggerOptions(app)

// Rest is same as above...
```

You can also modify the base logger as needed. Set this up _before_ adding any routes. Note that the function returns a low-level `Logger`, not a `SugaredLogger`.

```go
middleware.NewLogger = func() (*zap.Logger, error) {
	l, err := middleware.NewDefaultLogger()
	if err != nil {
		return nil, err
	}

	// Add your own global tags.
	l = l.With(zap.String("env", "prod"))

	return l, nil
}
```

You can also modify the logger in the current request's context from resolvers or operation handlers. This modifies the context in-place for the lifetime of the request.

```go
original := middleware.GetLogger(ctx)
modified := original.With("my-value", 123)
middleware.SetLoggerInContext(ctx, modified)
```

### Getting Operation Info

When setting up logging (or metrics, or auditing) you may want to have access to some additional information like the ID of the current operation. You can fetch this from the context **after** the handler has run.

```go
app.Middleware(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First, make sure the handler function runs!
		next.ServeHTTP(w, r)

		// After that, you can get the operation info.
		opInfo := GetOperationInfo(r.Context())
		fmt.Println(opInfo.ID)
		fmt.Println(opInfo.URITemplate)
	})
})
```

## Changing the Documentation Renderer

You can choose between [RapiDoc](https://mrin9.github.io/RapiDoc/), [ReDoc](https://github.com/Redocly/redoc), or [SwaggerUI](https://swagger.io/tools/swagger-ui/) to auto-generate documentation. Simply set the documentation handler on the router:

```go
app := cli.NewRouter("My API", "1.0.0")
app.DocsHandler(huma.ReDocHandler(app.Router))
```

> :whale: Pass a custom handler function to have even more control for branding or browser authentication.

## OpenAPI

By default, the generated OpenAPI spec, schemas, and autogenerated documentation are served in the root at `/openapi.json`, `/schemas`, and `/docs` respectively. The default prefix for all, and the suffix for each individual route can be modified:

```go
// Serve `/public/openapi.json`, `/public/schemas`, and `/public/docs`:
app.DocsPrefix("/public")
```

```go
// Serve `/internal/app/myService/model/openapi.json`, `/internal/app/myService/model/schemas`, and `/internal/app/myService/documentation`:
app.DocsPrefix("/internal/app/myService")
app.DocsSuffix("documentation")
app.SchemasSuffix("model/schemas")
app.SpecSuffix("model/openapi")
```

### Custom OpenAPI Fields

Use the OpenAPI hook for OpenAPI customization. It gives you a `*gabs.Container` instance that represents the root of the OpenAPI document.

```go
func modify(openapi *gabs.Container) {
	openapi.Set("value", "paths", "/test", "get", "x-foo")
}

app := cli.NewRouter("My API", "1.0.0")
app.OpenAPIHook(modify)
```

> :whale: See the [OpenAPI 3 spec](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md) for everything that can be set.

### JSON Schema

Each resource operation also returns a `describedby` HTTP link relation which references a JSON-Schema file. These schemas re-use the `DocsPrefix` and `SchemasSuffix` described above and default to the server root. For example:

```http
Link: </schemas/Note.json>; rel="describedby"
```

Object resources (i.e. not arrays) can also optionally return a `$schema` property with such a link, which enables the described-by relationship to outlive the HTTP request (i.e. saving the body to a file for later editing) and enables some editors like [VSCode](https://code.visualstudio.com/docs/languages/json#_mapping-in-the-json) to provide code completion and validation as you type.

```json
{
  "$schema": "http://localhost:8888/schemas/Note.json",
  "title": "I am a note title",
  "contents": "Example note contents",
  "labels": ["todo"]
}
```

Operations which accept objects as input will ignore the `$schema` property, so it is safe to submit back to the API.

This feature can be disabled if desired by using `app.DisableSchemaProperty`.

## GraphQL

Huma includes an optional, built-in, read-only GraphQL interface that can be enabled via `app.EnableGraphQL(config)`. It is mostly automatic and will re-use all your defined resources, read operations, and their params, headers, and models. By default it is available at `/graphql`.

If you want your resources to automatically fill in params, such as an item's ID from a list result, you must tell Huma how to map fields of the response to the correct parameter name. This is accomplished via the `graphParam` struct field tag. For example, given the following resources:

```go
app.Resource("/notes").Get("list-notes", "docs",
	responses.OK().Headers("Link").Model([]NoteSummary{}),
).Run(func(ctx huma.Context, input struct {
	Cursor string `query:"cursor" doc:"Pagination cursor"`
	Limit  int    `query:"limit" doc:"Number of items to return"`
}) {
	// Handler implementation goes here...
})

app.Resource("/notes/{note-id}").Get("get-note", "docs",
	responses.OK().Model(Note{}),
).Run(func(ctx huma.Context, input struct {
	NodeID string `path:"note-id"`
}) {
	// Handler implementation goes here...
})
```

You would map the `/notes` response to the `/notes/{note-id}` request with a `graphParam` tag on the response struct's field that tells Huma that the `note-id` parameter in URLs can be loaded directly from the `id` field of the response object.

```go
type NoteSummary struct {
	ID string `json:"id" graphParam:"note-id"`
}
```

Whenever a list of items is returned, you can access the detailed item via the name+"Item", e.g. `notesItem` would return the `get-note` response.

Then you can make requests against the service like `http://localhost:8888/graphql?query={notes{edges{id%20notesItem{contents}}}}`.

See the `graphql_test.go` file for a full-fledged example.

> :whale: Note that because Huma knows nothing about your database, there is no way to make efficient queries to only select the fields that were requested. This GraphQL layer works by making normal HTTP requests to your service as needed to fulfill the query. Even with that caveat it can greatly simplify and speed up frontend requests.

### GraphQL List Responses

HTTP responses may be lists, such as the `list-notes` example operation above. Since GraphQL responses need to account for more than just the response body (i.e. headers), Huma returns this as a wrapper object similar to but as a more general form of [Relay's Cursor Connections](https://relay.dev/graphql/connections.htm) pattern. The structure knows how to parse link relationship headers and looks like:

```
{
	"edges": [... your responses here...],
	"links": {
		"next": [
			{"key": "param1", "value": "value1"},
			{"key": "param2", "value": "value2"},
			...
		]
	}
	"headers": {
		"headerName": "headerValue"
	}
}
```

If you want a different paginator then this can be configured by creating your own struct which includes a field of `huma.GraphQLItems` and which implements the `huma.GraphQLPaginator` interface. For example:

```go
// First, define the custom paginator. This does nothing but return the list
// of items and ignores the headers.
type MySimplePaginator struct {
	Items huma.GraphQLItems `json:"items"`
}

func (m *MySimplePaginator) Load(headers map[string]string, body []interface{}) error {
	// Huma creates a new instance of your paginator before calling `Load`, so
	// here you populate the instance with the response data as needed.
	m.Items = body
	return nil
}

// Then, tell your app to use it when enabling GraphQL.
app.EnableGraphQL(&huma.GraphQLConfig{
	Paginator: &MySimplePaginator{},
})
```

Using the same mechanism above you can support Relay Collections or any other pagination spec as long as your underlying HTTP API supports the inputs/outputs required for populating the paginator structs.

### Custom GraphQL Path

You can set a custom path for the GraphQL endpoint:

```go
app.EnableGraphQL(&huma.GraphQLConfig{
	Path: "/graphql",
})
```

### Enabling the GraphiQL UI

You can turn on a UI for writing and making queries with schema documentation via the GraphQL config:

```go
app.EnableGraphQL(&huma.GraphQLConfig{
	GraphiQL: true,
})
```

It is [recommended](https://graphql.org/learn/serving-over-http/#graphiql) to turn GraphiQL off in production. Instead a tool like [graphqurl](https://github.com/hasura/graphqurl) can be useful for using GraphiQL in production on the client side, and it supports custom headers for e.g. auth. Don't forget to enable CORS via e.g. [`rs/cors`](https://github.com/rs/cors) so browsers allow access.

### GraphQL Query Complexity Limits

You can limit the maximum query complexity your server allows:

```go
app.EnableGraphQL(&huma.GraphQLConfig{
	ComplexityLimit: 250,
})
```

Complexity is a rough measure of the request load against your service and is calculated as the following:

| Field Type                       |                         Complexity |
| -------------------------------- | ---------------------------------: |
| Enum                             |                                  0 |
| Scalar (e.g. int, float, string) |                                  0 |
| Plain array / object             |                                  0 |
| Resource object                  |                                  1 |
| Array of resources               | count + (childComplexity \* count) |

`childComplexity` is the total complexity of any child selectors and the `count` is determined by passed in parameters like `first`, `last`, `count`, `limit`, `records`, or `pageSize` with a built-in default multiplier of `10`.

If a single resource is a child of a list, then the resource's complexity is also multiplied by the number of resources. This means nested queries that make list calls get very expensive fast. For example:

```
{
	categories(first: 10) {
		edges {
			catgoriesItem {
				products(first: 10) {
					edges {
						productsItem {
							id
							price
						}
					}
				}
			}
		}
	}
}
```

Because you are fetching up to 10 categories, and for each of those fetching a `categoriesItem` object and up to 10 products within each category, then a `productsItem` for each product, this results in:

```
Calculation:
(((1 producstItem * 10 products) + 10 products) + 1 categoriesItem) * 10 categories + 10 categories

Result:
220 complexity
```

## CLI

The `cli` package provides a convenience layer to create a simple CLI for your server, which lets a user set the host, port, TLS settings, etc when running your service.

```go
app := cli.NewRouter("My API", "1.0.0")

// Do resource/operation setup here...

app.Run()
```

Then run the service:

```sh
$ go run yourservice.go --help
```

## CLI Runtime Arguments & Configuration

The CLI can be configured in multiple ways. In order of decreasing precedence:

1. Commandline arguments, e.g. `-p 8000` or `--port=8000`
2. Environment variables prefixed with `SERVICE_`, e.g. `SERVICE_PORT=8000`

It's also possible to load configured flags from config files. JSON/YAML/TOML are supported. For example, to load `some/path/my-app.json` you can do the following before calling `app.Run()`:

```go
viper.AddConfigPath("some/path")
viper.SetConfigName("my-app")
viper.ReadInConfig()
```

## Custom CLI Arguments

You can add additional CLI arguments, e.g. for additional logging tags. Use the `Flag` method along with the `viper` module to get the parsed value.

```go
app := cli.NewRouter("My API", "1.0.0")

// Add a long arg (--env), short (-e), description & default
app.Flag("env", "e", "Environment", "local")

r.Resource("/current_env").Get("get-env", "Get current env",
	responses.String(http.StatusOK),
).Run(func(ctx huma.Context) {
	// The flag is automatically bound to viper settings using the same name.
	ctx.Write([]byte(viper.GetString("env")))
})
```

Then run the service:

```sh
$ go run yourservice.go --env=prod
```

Note that passed flags are not parsed during application setup. They only get parsed after calling `app.Run()`, so if you need their value for some setup code you can use the `ArgsParsed` handler:

```go
app.ArgsParsed(func() {
	fmt.Printf("Env is %s\n", viper.GetString("env"))
})
```

See lazy loading below for more details.

> :whale: Combine custom arguments with [customized logger setup](#customizing-logging) and you can easily log your cloud provider, environment, region, pod, etc with every message.

## Custom CLI Commands

You can access the root `cobra.Command` via `app.Root()` and add new custom commands via `app.Root().AddCommand(...)`. The `openapi` sub-command is one such example in the default setup.

> :whale: You can also overwite `app.Root().Run` to completely customize how you run the server. Or just ditch the `cli` package completely.

## Lazy-loading at Server Startup

You can register functions to run before any command handler or before the server starts, allowing for things like lazy-loading dependencies. It is safe to call these methods multiple times.

```go
var db *mongo.Client

app := cli.NewRouter("My API", "1.0.0")

// Add a long arg (--env), short (-e), description & default
app.Flag("env", "e", "Environment", "local")

app.ArgsParsed(func() {
	// Arguments have been parsed now. This runs before *any* command including
	// custom commands, not just server-startup.
	fmt.Println(viper.GetString("env"))
})

app.PreStart(func() {
	// Server is starting up, so connect to the datastore. This runs only
	// before server start.
	var err error
	db, err = mongo.Connect(context.Background(),
		options.Client().ApplyURI("..."))
})
```

> :whale: This is especially useful for external dependencies and if any custom CLI commands are set up. For example, you may not want to require a database to run `my-service openapi my-api.json`.

## Testing

The Go standard library provides useful testing utilities and Huma routers implement the [`http.Handler`](https://golang.org/pkg/net/http/#Handler) interface they expect. Huma also provides a `humatest` package with utilities for creating test routers capable of e.g. capturing logs.

You can see an example in the [`examples/test`](https://github.com/danielgtaylor/huma/tree/master/examples/test) directory:

```go
package main

import (
	"github.com/danielgtaylor/huma"
	"github.com/danielgtaylor/huma/cli"
	"github.com/danielgtaylor/huma/responses"
)

func routes(r *huma.Router) {
	// Register a single test route that returns a text/plain response.
	r.Resource("/test").Get("test", "Test route",
		responses.OK().ContentType("text/plain"),
	).Run(func(ctx huma.Context) {
		ctx.Write([]byte("Hello, test!"))
	})
}

func main() {
	// Create the router.
	app := cli.NewRouter("Test", "1.0.0")

	// Register routes.
	routes(app.Router)

	// Run the service.
	app.Run()
}
```

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/humatest"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// Set up the test router and register the routes.
	r := humatest.NewRouter(t)
	routes(r)

	// Make a request against the service.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Assert the response is as expected.
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, test!", w.Body.String())
}
```

# Design

General Huma design principles:

- HTTP/2 and streaming out of the box
- Describe inputs/outputs and keep docs up to date
- Generate OpenAPI for automated tooling
- Re-use idiomatic Go concepts whenever possible
- Encourage good behavior, e.g. exhaustive errors

## High-level design

The high-level design centers around a `Router` object.

- CLI (optional)
  - Router
    - []Middleware
    - []Resource
      - URI path
      - []Middleware
      - []Operations
        - HTTP method
        - Inputs / outputs
          - Go structs with tags
        - Handler function

## Router Selection

- Why not Gin? Lots of stars on GitHub, but... Overkill, non-standard handlers & middlware, weird debug mode.
- Why not fasthttp? Fiber? Not fully HTTP compliant, no HTTP/2, no streaming request/response support.
- Why not httprouter? Non-standard handlers, no middleware.
- HTTP/2 means HTTP pipelining benchmarks don't really matter.

Ultimately using Chi because:

- Fast router with support for parameterized paths & middleware
- Standard HTTP handlers
- Standard HTTP middleware

### Compatibility

Huma tries to be compatible with as many Go libraries as possible by using standard interfaces and idiomatic concepts.

- Standard middleware `func(next http.Handler) http.Handler`
- Standard context `huma.Context` is a `context.Context`
- Standard HTTP writer `huma.Context` is an `http.ResponseWriter` that can check against declared response codes and models.
- Standard streaming support via the `io.Reader` and `io.Writer` interfaces

## Compromises

Given the features of Go, the desire to strictly keep the code and docs/tools in sync, and a desire to be developer-friendly and quick to start using, Huma makes some necessary compromises.

- Struct tags are used as metadata for fields to support things like JSON Schema-style validation. There are no compile-time checks for these, but basic linter support.
- Handler functions registration uses `interface{}` to support any kind of input struct.
- Response registration takes an _instance_ of your type since you can't pass types in Go.
- Many checks happen at service startup rather than compile-time. Luckily the most basic unit test that creates a router should catch these.
- `ctx.WriteModel` and `ctx.WriteError` do checks at runtime and can be at least partially bypassed with `ctx.Write` by design. We trade looser checks for a nicer interface and more compatibility.

> :whale: Thanks for reading!
