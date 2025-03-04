// Copyright 2020 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var target = flag.String("target", "", "the target to fetch from")
var buildID = flag.String("build_id", "", "the build id to fetch from")
var artifact = flag.String("artifact", "", "the artifact to download")
var output = flag.String("output", "", "the file name to save as")
var writeToStdout = false

func errPrint(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()
	// We only support passing 1 argument `-` so if we have more than
	// one argument this is an error state,
	if len(args) > 1 {
		errPrint("Error: Too many arguments passed to fetch_artifact.")
	}

	if len(args) > 0 {
		writeToStdout = args[len(args)-1] == "-"
		if !writeToStdout {
			errPrint(fmt.Sprintf(
				"Error: Only supported final argument to fetch_artifact is `-` but got `%s`.", args[len(args)-1]))
		}

		if len(*output) > 0 && writeToStdout {
			errPrint("Error: Both '-output' and '-' flags can not be used together.")
		}
	}

	url := fmt.Sprintf("https://androidbuildinternal.googleapis.com/android/internal/build/v3/builds/%s/%s/attempts/latest/artifacts/%s/url", *buildID, *target, *artifact)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errPrint(fmt.Sprintf("unable to build request %v", err))
	}
	req.Header.Set("Accept", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errPrint(fmt.Sprintf("Unable to make request %s.", err))
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		body, _ := ioutil.ReadAll(res.Body)
		errPrint(fmt.Sprintf("Unable to download artifact: %s\n %s.", res.Status, body))
	}

	if writeToStdout {
		io.Copy(os.Stdout, res.Body)
		return
	}

	fileName := *artifact
	if len(*output) > 0 {
		fileName = *output
	}

	f, err := os.Create(path.Base(fileName))
	if err != nil {
		errPrint(fmt.Sprintf("Unable to create file %s.", err))
	}
	defer f.Close()
	io.Copy(f, res.Body)
	fmt.Printf("File %s created.\n", f.Name())
}
