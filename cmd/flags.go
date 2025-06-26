// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// AccessToken is used to access repo stars and gain non-authorized
// rate limits.

// AccessTokenDesc describes usage.
const AccessTokenDesc = "GitHub access token for authorized rate limits"

func getAccessToken() (string, error) {
	// First check if token is provided via command line
	if len(AccessToken) > 0 {
		return AccessToken, nil
	}

	// Try to load from config (which includes environment variables)
	cfg := LoadConfigOrNil()
	if cfg != nil {
		if token := cfg.GetToken(); token != "" {
			return token, nil
		}
	}

	// No token found, return error
	return "", errors.New(`An access token must be specified via --token, config.yaml, or GH_TOKEN environment variable.

To generate an access token for accessing repo stars and gaining authorized
rate limits, see:

https://help.github.com/articles/creating-an-access-token-for-command-line-use/

When creating a token, ensure that only the public_repo permission is enabled.

You can provide the token in several ways:
1. Command line: --token=YOUR_TOKEN
2. Config file: Add 'token: "YOUR_TOKEN"' to the github section in config.yaml
3. Environment variable: Set GH_TOKEN=YOUR_TOKEN
4. .env file: Create a .env file with GH_TOKEN=YOUR_TOKEN
`)
}

// CacheDirDesc describes usage.
const CacheDirDesc = "directory for storing cached GitHub API responses"

// RepoDesc describes usage.
const RepoDesc = "GitHub owner and repository, formatted as :owner/:repo"

func mustUsage(cmd *cobra.Command) {
	if err := cmd.Usage(); err != nil {
		panic(err)
	}
}
