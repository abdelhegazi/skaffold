/*
Copyright 2018 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package docker

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

	cmdutil "github.com/GoogleContainerTools/skaffold/cmd/skaffold/app/cmd/util"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var output string

func NewCmdContext(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Outputs a minimal context tarball to stdout",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runContext(out, filename, context)
		},
	}
	AddDockerFlags(cmd)
	return cmd
}

func runContext(out io.Writer, filename, context string) error {
	config, err := cmdutil.ParseConfig(filename)
	if err != nil {
		return errors.Wrap(err, "parsing skaffold config")
	}
	dockerFilePath, err := filepath.Abs(filename)
	logrus.Info(filename)
	logrus.Info(dockerFilePath)
	if err != nil {
		return err
	}

	// Write everything to memory, then flush to disk at the end.
	// This prevents recursion problems, where the output file can end up
	// in the context itself during creation.
	var b bytes.Buffer
	if err := docker.CreateDockerTarGzContext(getBuildArgsForDockerfile(config, dockerfile), &b, context, dockerFilePath); err != nil {
		return err
	}
	return ioutil.WriteFile(output, b.Bytes(), 0644)
}
