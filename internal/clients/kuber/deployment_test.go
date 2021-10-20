package kuber

import (
	"fmt"
	"os"
	"testing"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	"k8s.io/apimachinery/pkg/api/resource"
)

// For every 8GB of memory, 1000 millicores of CPU per ticket #76
func TestCpuProportion(t *testing.T) {
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	global.Settings.PoolMachineTotalCpuMillicores = 2000
	global.Settings.PoolMachineTotalMemoryBytes = 17179869184

	tests := [][]string{
		{"8Gi", "1000m"}, {"4Gi", "500m"}, {"1Gi", "125m"},
		{"512Mi", "62m"}, {"256Mi", "31m"}, {"128Mi", "15m"},
		{"64Mi", "7m"}, {"32Mi", "3m"}, {"16Mi", "1m"},
	}
	for _, value := range tests {
		memoryQuantity := resource.MustParse(value[0])
		expectedQuantity := resource.MustParse(value[1])
		cpuQuantity := cpuProportion(memoryQuantity)

		if cpuQuantity.MilliValue() != expectedQuantity.MilliValue() {
			t.Error(fmt.Sprintf(
				"cpuProportion(%d) = %d; want %d",
				memoryQuantity.Value(),
				cpuQuantity.MilliValue(),
				expectedQuantity.MilliValue(),
			))
		}
	}
}

func TestPodCpuProportion(t *testing.T) {
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	global.Settings.DefaultAutoscalingCpuThreshold = 75
	global.Settings.DefaultAutoscalingCpuThresholdEpsilon = 2

	type testItem struct {
		totalMemory string
		expected    string
	}

	tests := []testItem{
		{
			totalMemory: "8Gi",
			expected:    "770m",
		},
		{
			totalMemory: "9Gi",
			expected:    "867m",
		},
		{
			totalMemory: "16Mi",
			expected:    "1m",
		},
	}

	for _, value := range tests {
		memoryQuantity := resource.MustParse(value.totalMemory)
		expectedQuantity := resource.MustParse(value.expected)
		cpuQuantity := podCpuProportion(&memoryQuantity)

		if cpuQuantity.MilliValue() != expectedQuantity.MilliValue() {
			t.Error(fmt.Sprintf("podCpuProportion(%d) = %d; want %d",
				memoryQuantity.Value(),
				cpuQuantity.MilliValue(),
				expectedQuantity.MilliValue(),
			))
		}
	}
}

func TestCreateClient(t *testing.T) {
	t.Skip("stub")

	// This doesn't work from this directory.
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	config, err := clients.InitializeKubeConfig()
	if err != nil {
		panic(err)
	}

	kuberClient, err := NewClient(config)
	if err != nil {
		panic(err)
	}

	if kuberClient == nil {
		panic(kuberClient)
	}
}

func Test_ephemeralStorageProportion(t *testing.T) {
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	global.Settings.EphemeralStorageCoefficient = 1.9

	type Test struct {
		input  string
		output string
	}

	tests := []Test{
		{
			input:  "0.5Gi",
			output: "0.95Gi",
		},
		{
			input:  "1Gi",
			output: "1.9Gi",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%v -> %v", test.input, test.output), func(t *testing.T) {
			input, _ := resource.ParseQuantity(test.input)
			output, _ := resource.ParseQuantity(test.output)
			if got := ephemeralStorageProportion(input); got.Value() != output.Value() {
				t.Errorf("ephemeralStorageProportion() = %v, want %v", got.Value(), output.Value())
			}
		})
	}
}
