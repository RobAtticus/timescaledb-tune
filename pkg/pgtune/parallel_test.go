package pgtune

import (
	"math/rand"
	"testing"
)

// parallelSettingsMatrix stores the test cases for ParallelRecommender along
// with the expected values for its keys
var parallelSettingsMatrix = map[int]map[string]string{
	2: map[string]string{
		MaxWorkerProcessesKey:       "2",
		MaxParallelWorkersGatherKey: "1",
		MaxParallelWorkers:          "2",
	},
	4: map[string]string{
		MaxWorkerProcessesKey:       "4",
		MaxParallelWorkersGatherKey: "2",
		MaxParallelWorkers:          "4",
	},
	5: map[string]string{
		MaxWorkerProcessesKey:       "5",
		MaxParallelWorkersGatherKey: "3",
		MaxParallelWorkers:          "5",
	},
}

func TestNewParallelRecommender(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		cpus := rand.Intn(128)
		r := NewParallelRecommender(cpus)
		if r == nil {
			t.Errorf("unexpected nil recommender")
		}
		if got := r.cpus; got != cpus {
			t.Errorf("recommender has incorrect cpus: got %d want %d", got, cpus)
		}
	}
}

func TestParallelRecommenderIsAvailable(t *testing.T) {
	if r := NewParallelRecommender(0); r.IsAvailable() {
		t.Errorf("unexpectedly available for 0 cpus")
	}
	if r := NewParallelRecommender(1); r.IsAvailable() {
		t.Errorf("unexpectedly available for 1 cpus")
	}

	for i := 2; i < 1000; i++ {
		if r := NewParallelRecommender(i); !r.IsAvailable() {
			t.Errorf("unexpected UNavailable for %d cpus", i)
		}
	}
}

func TestParallelRecommenderRecommend(t *testing.T) {
	for cpus, matrix := range parallelSettingsMatrix {
		r := &ParallelRecommender{cpus}
		testRecommender(t, r, matrix)
	}
}

func TestParallelRecommenderRecommendPanics(t *testing.T) {
	func() {
		r := &ParallelRecommender{5}
		defer func() {
			if re := recover(); re == nil {
				t.Errorf("did not panic when should")
			}
		}()
		r.Recommend("foo")
	}()

	func() {
		r := &ParallelRecommender{1}
		defer func() {
			if re := recover(); re == nil {
				t.Errorf("did not panic when should")
			}
		}()
		r.Recommend("foo")
	}()
}

func TestParallelSettingsGroup(t *testing.T) {
	for cpus, matrix := range parallelSettingsMatrix {
		config := NewSystemConfig(1024, cpus, "9.6")
		sg := GetSettingsGroup(ParallelLabel, config)
		testSettingGroup(t, sg, matrix, ParallelLabel, ParallelKeys)

		// PG10 adds a key
		config.PGMajorVersion = "10"
		sg = GetSettingsGroup(ParallelLabel, config)
		testSettingGroup(t, sg, matrix, ParallelLabel, ParallelKeys)
	}

}
