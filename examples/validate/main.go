package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

type MyThing struct {
	ID string `json:"id" format:"uuid"`
}

var thingType = reflect.TypeOf(MyThing{})

func validateThing(b []byte, reg huma.Registry) error {
	s := reg.Map()[thingType.Name()]

	var parsed any
	json.Unmarshal(b, &parsed)

	res := huma.ValidateResult{}
	huma.Validate(reg, s, huma.NewPathBuffer([]byte(""), 0), huma.ModeReadFromServer, parsed, &res)

	if len(res.Errors) > 0 {
		return errors.Join(res.Errors...)
	}

	return nil
}

func main() {
	r := fiber.New()
	api := humafiber.New(r, huma.DefaultConfig("Validate", "1.0.0"))
	api.OpenAPI().Components.Schemas.Schema(
		thingType,
		true,
		thingType.Name(),
	)

	thing1, err := os.ReadFile("valid.json")
	if err != nil {
		panic(err)
	}
	thing2, err := os.ReadFile("invalid.json")
	if err != nil {
		panic(err)
	}

	err = validateThing(thing1,
		api.OpenAPI().Components.Schemas)
	if err != nil {
		panic(err)
	}
	fmt.Println("valid.json is valid.")

	err = validateThing(thing2,
		api.OpenAPI().Components.Schemas)
	if err == nil {
		panic("expected validation errors from invalid.json")
	}
	fmt.Printf("invalid.json is invalid:\n\t%v\n", err)
}
