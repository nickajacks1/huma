package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/autopatch"
	"github.com/danielgtaylor/huma/v2/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Options struct {
	// cli.DefaultOptions?
	Debug bool
	Host  string `doc:"Hostname to listen on."`
	Port  int    `doc:"Port to listen on." short:"p" default:"3001"`
}

type GreetingInputBody struct {
	Suffix   string `json:"suffix" default:"!" maxLength:"5"`
	Optional int    `json:"optional,omitempty" default:"2"`
}

func (b *GreetingInputBody) Resolve(ctx huma.Context) []error {
	if strings.Contains(b.Suffix, "err") {
		return []error{&huma.ErrorDetail{
			Location: "body.suffix",
			Message:  "Nested resolver works",
			Value:    b.Suffix,
		}}
	}
	return nil
}

type GreetingInput struct {
	ID   string `path:"id" example:"abc123" maxLength:"10"`
	Num  int    `query:"num" minimum:"0" default:"7"`
	Body GreetingInputBody
	// Body struct {
	// 	Suffix string `json:"suffix" default:"!" maxLength:"5"`
	// }
}

func (i *GreetingInput) Resolve(huma.Context) []error {
	if i.Body.Suffix == "reserr" {
		return []error{&huma.ErrorDetail{
			Location: "body.suffix",
			Message:  "Suffix weird and I don't like it input",
			Value:    i.Body.Suffix,
		}}
	}
	return nil
}

type GreetingOutputSub struct {
	Foo  int `json:"foo"`
	Sub2 struct {
		ThisFails string `json:"this_fails"`
	}
}

type GreetingOutput struct {
	ETag string `header:"ETag"`
	Body struct {
		Greeting string            `json:"greeting"`
		Suffix   string            `json:"suffix"`
		Total    int               `json:"total"`
		Sub      GreetingOutputSub `json:"sub"`
		Sub2     struct {
			Bar string `json:"bar"`
		} `json:"sub2"`
	}
}

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "greet",
		Method:      http.MethodPost,
		Path:        "/foo/{id}",
		Tags:        []string{"Greetings"},
		Errors: []int{
			http.StatusBadRequest,
			http.StatusNotFound,
		},
		// Example showing headers merged with the default response...
		// function composition should make this easier... provide example of that?
		Responses: map[string]*huma.Response{
			"200": {
				Headers: map[string]*huma.Header{
					"Other": {
						Example: "abc123",
					},
				},
			},
		},
		Extensions: map[string]any{
			"x-my-extension": "value",
		},
	}, func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
		if input.ID == "error" {
			return nil, huma.Error404NotFound("can't find greeting", &huma.ErrorDetail{
				Location: "path.id",
				Message:  "ID not found",
				Value:    input.ID,
			}, fmt.Errorf("plain error"))
		}

		if input.ID == "plain" {
			return nil, fmt.Errorf("plain error")
		}

		// fmt.Println("optional", input.Body.Optional)

		resp := &GreetingOutput{}
		resp.ETag = "abc123"
		resp.Body.Greeting = "Hello, " + input.ID + input.Body.Suffix
		resp.Body.Suffix = input.Body.Suffix
		resp.Body.Total = len(resp.Body.Greeting)
		resp.Body.Sub.Foo = input.Num
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "no-body",
		Method:      http.MethodGet,
		Path:        "/no-body",
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		return nil, nil
	})

	type Patch struct {
		Body struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}
	}

	huma.Register(api, huma.Operation{
		OperationID: "patch-get",
		Method:      http.MethodGet,
		Path:        "/patch-test",
	}, func(ctx context.Context, input *struct{}) (*Patch, error) {
		resp := &Patch{}
		resp.Body.Foo = "foo"
		resp.Body.Bar = "bar"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "patch-put",
		Method:      http.MethodPut,
		Path:        "/patch-test",
	}, func(ctx context.Context, input *Patch) (*Patch, error) {
		return input, nil
	})
}

func main() {
	var api huma.API

	cli := huma.NewCLI(func(cli huma.CLI, opts *Options) {
		// r := chi.NewMux()
		// api := huma.NewChi(r, huma.Config{
		// 	OpenAPI: &huma.OpenAPI{
		// 		Info: &huma.Info{
		// 			Title:   "My API",
		// 			Version: "1.0.0",
		// 		},
		// 		Servers: []*huma.Server{
		// 			{URL: "http://localhost:3001"},
		// 		},
		// 	},
		// })
		r := fiber.New()
		// r.Use(logger.New())
		// r.Use(recover.New())
		r.Use(compress.New())
		r.Use(requestid.New())

		// Add a custom health check
		r.Get("/health", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		config := huma.DefaultConfig("My API", "1.0.0")
		config.Transformers = append(config.Transformers, huma.FieldSelectTransform)
		api = humafiber.New(r, config)
		// api = humafiber.New(r, huma.Config{
		// 	OpenAPI: &huma.OpenAPI{
		// 		Info: &huma.Info{
		// 			Title:   "My API",
		// 			Version: "1.0.0",
		// 		},
		// 		Servers: []*huma.Server{
		// 			{URL: "http://localhost:3001"},
		// 		},
		// 		Extensions: map[string]interface{}{
		// 			"x-foo": "bar",
		// 		},
		// 	},
		// })

		// huma.AutoRegister(api, Things{...})
		// huma.Register(api, http.MethodGet, "/foo/{id}", func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
		// 	return &GreetingOutput{"Hello, " + input.ID}, nil
		// })
		RegisterRoutes(api)

		autopatch.AutoPatch(api)

		cli.OnStart(func() {
			// Connect dependencies here...
			// things.Init(...)
			r.Listen(fmt.Sprintf("%s:%d", opts.Host, opts.Port))
			// http.ListenAndServe(":3001", r)
		})

		cli.OnStop(func() {
			r.ShutdownWithTimeout(5 * time.Second)
		})
	})

	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: huma.WithOptions(
			func(cmd *cobra.Command, args []string, options *Options) {
				if options.Debug {
					fmt.Println("Debug mode enabled")
				}
				b, _ := yaml.Marshal(api.OpenAPI())
				fmt.Println(string(b))
			},
		),
	})

	cli.Run()
}