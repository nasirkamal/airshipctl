/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package phase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/phase"
)

const (
	testFactoryErr   = "test config error"
	testNewHelperErr = "Missing configuration"
	testNoBundlePath = "no such file or directory"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name        string
		errContains string
		runFlags    phase.RunFlags
		factory     config.Factory
	}{
		{
			name: "Error config factory",
			factory: func() (*config.Config, error) {
				return nil, fmt.Errorf(testFactoryErr)
			},
			errContains: testFactoryErr,
		},
		{
			name: "Error new helper",
			factory: func() (*config.Config, error) {
				return &config.Config{
					CurrentContext: "does not exist",
					Contexts:       make(map[string]*config.Context),
				}, nil
			},
			errContains: testNewHelperErr,
		},
		{
			name: "Error phase by id",
			factory: func() (*config.Config, error) {
				conf := config.NewConfig()
				conf.Manifests = map[string]*config.Manifest{
					"manifest": {
						MetadataPath: "broken_metadata.yaml",
						TargetPath:   "testdata",
					},
				}
				conf.CurrentContext = "context"
				conf.Contexts = map[string]*config.Context{
					"context": {
						Manifest: "manifest",
					},
				}
				return conf, nil
			},
			errContains: testNoBundlePath,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			command := phase.RunCommand{
				Options: tt.runFlags,
				Factory: tt.factory,
			}
			err := command.RunE()
			if tt.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlanCommand(t *testing.T) {
	tests := []struct {
		name        string
		errContains string
		runFlags    phase.RunFlags
		factory     config.Factory
	}{
		{
			name: "Error config factory",
			factory: func() (*config.Config, error) {
				return nil, fmt.Errorf(testFactoryErr)
			},
			errContains: testFactoryErr,
		},
		{
			name: "Error new helper",
			factory: func() (*config.Config, error) {
				return &config.Config{
					CurrentContext: "does not exist",
					Contexts:       make(map[string]*config.Context),
				}, nil
			},
			errContains: testNewHelperErr,
		},
		{
			name: "Error phase by id",
			factory: func() (*config.Config, error) {
				conf := config.NewConfig()
				conf.Manifests = map[string]*config.Manifest{
					"manifest": {
						MetadataPath: "broken_metadata.yaml",
						TargetPath:   "testdata",
					},
				}
				conf.CurrentContext = "context"
				conf.Contexts = map[string]*config.Context{
					"context": {
						Manifest: "manifest",
					},
				}
				return conf, nil
			},
			errContains: testNoBundlePath,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			command := phase.PlanCommand{
				Factory: tt.factory,
			}
			err := command.RunE()
			if tt.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
