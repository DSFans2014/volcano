/*
Copyright 2023 The Volcano Authors.

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

/*
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

package ascend

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
        //"volcano.sh/volcano/pkg/scheduler/api/devices"
	"volcano.sh/volcano/pkg/scheduler/api/devices/config"
)

var config_yaml = `
vnpus:
- chipName: 310P3

  commonWord: Ascend310P
  resourceName: huawei.com/Ascend310P
  resourceMemoryName: huawei.com/Ascend310P-memory
  memoryAllocatable: 21527
  memoryCapacity: 24576
  aiCore: 8
  aiCPU: 7
  templates:
    - name: vir01
      memory: 3072
      aiCore: 1
      aiCPU: 1
    - name: vir02
      memory: 6144
      aiCore: 2
      aiCPU: 2
    - name: vir04
      memory: 12288
      aiCore: 4
      aiCPU: 4
- chipName: 910B3
  commonWord: Ascend910B3
  resourceName: huawei.com/Ascend910B3
  resourceMemoryName: huawei.com/Ascend910B3-memory
  memoryAllocatable: 65536
  memoryCapacity: 65536
  aiCore: 20
  aiCPU: 7
  templates:
    - name: vir05_1c_16g
      memory: 16384
      aiCore: 5
      aiCPU: 1
    - name: vir10_3c_32g
      memory: 32768
      aiCore: 10
      aiCPU: 3
nvidia:
  resourceCountName: volcano.sh/vgpu-number
  resourceMemoryName: volcano.sh/vgpu-memory
  resourceMemoryPercentageName: volcano.sh/vgpu-memory-percentage
  resourceCoreName: volcano.sh/vgpu-cores
  overwriteEnv: false
  defaultMemory: 0
  defaultCores: 0
  defaultGPUNum: 1
  deviceSplitCount: 10
  deviceMemoryScaling: 1
  deviceCoreScaling: 1
  gpuMemoryFactor: 1
  knownMigGeometries:
  - models: [ "A30" ]
    allowedGeometries:
      - group: group1
        geometries: 
        - name: 1g.6gb
          memory: 6144
          count: 4
      - group: group2
        geometries: 
        - name: 2g.12gb
          memory: 12288
          count: 2
      - group: group3
        geometries: 
        - name: 4g.24gb
          memory: 24576
          count: 1
`

func yamlStringToConfig(yamlStr string) (*config.Config, error) {
    var config config.Config
    err := yaml.Unmarshal([]byte(yamlStr), &config)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
    }
    return &config, nil
}

func Test_trimMemory(t *testing.T) {
	_, err := yamlStringToConfig(config_yaml)
	conf, err := yamlStringToConfig(config_yaml)
	assert.Nil(t, err)
	dev := AscendDevice{
		config: conf.VNPUs[0],
	}
	tests := []struct {
		name      string
		inputMem  int64
		wantMem   int64
	}{
		{"test1", 0, 3072},
		{"test2", 1, 3072},
		{"test3", 100, 3072},
		{"test4", 3071, 3072},
		{"test5", 3072, 3072},
		{"test6", 3073, 6144},
		{"test7", 6143, 6144},
		{"test8", 6144, 6144},
		{"test9", 6144, 6144},
		{"test10", 6145, 12288},
		{"test11", 12288, 12288},
		{"test12", 12289, 21527},
		{"test13", 21527, 21527},
		{"test14", 21528, 21527},
		{"test15", 24576, 21527},
		{"test16", 24577, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		    got, _ := dev.trimMemory(tt.inputMem)
		    assert.Equal(t, tt.wantMem, got)
		})
	}
}
