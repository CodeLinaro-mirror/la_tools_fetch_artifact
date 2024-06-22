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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	target   = flag.String("target", "", "the target to fetch from")
	buildID  = flag.String("build_id", "", "the build id to fetch from, can use '-branch' to get the latest passed build ID")
	branch   = flag.String("branch", "", "the branch to fetch from, used when '-build_id' is not provided,\nit would fetch the latest successful build")
	artifact = flag.String("artifact", "", "the artifact to download")
	output   = flag.String("output", "", "the file name to save as")
)

var writeToStdout = false

type BuildResponse struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	BuildId string `json:"buildId"`
}

func errPrint(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func checkFlags() error {
	if len(*target) == 0 {
		return errors.New("missing target")
	}
	if len(*artifact) == 0 {
		return errors.New("missing artifact")
	}
	if len(*buildID) == 0 && len(*branch) == 0 {
		return errors.New("missing build_id or branch")
	}
	if len(*buildID) != 0 && len(*branch) != 0 {
		return errors.New("too many arguments. give only build ID or branch")
	}
	return nil
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

	if err := checkFlags(); err != nil {
		flag.Usage()
		errPrint(err.Error())
	}

	client := &http.Client{}

	if len(*buildID) == 0 {
		latestGoodBuildID, err := getLatestGoodBuild(client, *branch, *target)
		if err != nil {
			errPrint(fmt.Sprintf("Error fetching latest good buildID %s, err", err))
		}
		flag.Set("build_id", latestGoodBuildID)
	}
	err := fetchArtifact(client, *target, *buildID, *artifact, *output)
	if err != nil {
		errPrint(fmt.Sprintf("Fetch artifact error: %s", err))
	}
}

func sendRequest(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func fetchArtifact(client *http.Client, target string, buildID string, artifact string, output string) error {
	url := fmt.Sprintf("https://androidbuildinternal.googleapis.com/android/internal/build/v3/builds/%s/%s/attempts/latest/artifacts/%s/url", url.QueryEscape(buildID), url.QueryEscape(target), url.QueryEscape(artifact))
	res, err := sendRequest(client, url)
	if err != nil {
		return fmt.Errorf("error fetching artifact %w", err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		body, _ := ioutil.ReadAll(res.Body)
		errPrint(fmt.Sprintf("Unable to download artifact: %s\n %s.", res.Status, body))
	}

	if writeToStdout {
		io.Copy(os.Stdout, res.Body)
		return nil
	}

	fileName := artifact
	if len(output) > 0 {
		fileName = output
	}

	f, err := os.Create(path.Base(fileName))
	if err != nil {
		return fmt.Errorf("unable to create file %w", err)
	}
	defer f.Close()
	io.Copy(f, res.Body)
	fmt.Printf("File %s created.\n", f.Name())
	return nil
}

func getLatestGoodBuild(client *http.Client, branch string, target string) (string, error) {
	url := fmt.Sprintf("https://androidbuildinternal.googleapis.com/android/internal/build/v3/builds?branches=%s&buildAttemptStatus=complete&buildType=submitted&maxResults=1&successful=true&target=%s", url.QueryEscape(branch), url.QueryEscape(target))
	res, err := sendRequest(client, url)
	if err != nil {
		return "", fmt.Errorf("send request error: %w", err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("unable to get Build ID: %s\n %s", res.Status, body)
	}

	body, _ := io.ReadAll(res.Body)
	var buildData BuildResponse
	err = json.Unmarshal(body, &buildData)

	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}
	if len(buildData.Builds) == 0 {
		return "", errors.New("error no build ID is found")
	}

	return buildData.Builds[0].BuildId, nil
}
