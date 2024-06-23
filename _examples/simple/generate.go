// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:generate sh -c "mkdir -p ./schema"
//go:generate go run -mod=readonly generate.go

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/lrstanley/entrest"
	"github.com/ogen-go/ogen"
)

func main() {
	ex, err := entrest.NewExtension(&entrest.Config{
		Handler:           entrest.HandlerChi,
		EnableSpecHandler: true,
		GlobalRequestHeaders: map[string]*ogen.Header{
			"X-Request-Id": {
				Description: "A unique identifier for the request.",
				Required:    false,
				Schema:      &ogen.Schema{Type: "string"},
			},
		},
		GlobalResponseHeaders: map[string]*ogen.Header{
			"X-Ratelimit-Limit": {
				Description: "The maximum number of requests that the consumer is permitted to make in a given period.",
				Required:    true,
				Schema:      &ogen.Schema{Type: "integer"},
			},
			"X-Ratelimit-Remaining": {
				Description: "The number of requests remaining in the current rate limit window.",
				Required:    true,
				Schema:      &ogen.Schema{Type: "integer"},
			},
			"X-Ratelimit-Reset": {
				Description: "The time at which the current rate limit window resets in UTC epoch seconds.",
				Required:    true,
				Schema:      &ogen.Schema{Type: "integer"},
			},
		},
	})
	if err != nil {
		log.Fatalf("creating entrest extension: %v", err)
	}

	err = entc.Generate(
		"./schema",
		&gen.Config{
			Target:   "./ent",
			Schema:   "github.com/lrstanley/entrest/_examples/simple/schema",
			Package:  "github.com/lrstanley/entrest/_examples/simple/ent",
			Features: []gen.Feature{},
		},
		entc.Extensions(ex),
	)
	if err != nil {
		log.Fatalf("failed to run ent codegen: %v", err)
	}
}