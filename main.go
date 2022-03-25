// BSD 3-Clause License

// Copyright (c) 2021, Michael Grigoryan
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.

// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	v1 "polygon.am/core/api/v1"
	"polygon.am/core/pkg/config"
	"polygon.am/core/pkg/types"
)

// Global, configuration variable for accessing and changing
// the configuration on demand.
var Configuration *types.Config

// The default path for looking for the default configuration
// file path, if the environment variable was not supplied.
const DefaultConfigurationFilePath string = "./.conf.yaml"

func init() {
	path, err := filepath.Abs(DefaultConfigurationFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.ParseConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	// Assigning parsed configuration to a global variable
	Configuration = config
}

func main() {
	r := chi.NewRouter()

	if "production" != os.Getenv("POLYGON_CORE_CONFIG_ENV") {
		// Only enabling route logging in development
		r.Use(middleware.Logger)
	}

	r.Use(middleware.GetHead)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/status"))
	r.Use(httprate.LimitAll(100, 1*time.Minute))

	r.Mount("/api/v1", v1.Router())
	log.Println("getpolygon/corexp started at http://" + Configuration.Polygon.Addr)

	// Binding to the address specified or defaulted to from the configuration
	// and attaching chi routes to the server.
	http.ListenAndServe(Configuration.Polygon.Addr, r)
}
