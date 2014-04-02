package steam

import "testing"

func TestUserCountGetter(t *testing.T) {
	samples := []int{0, 10, 20, 730, 550}
	for _, sample := range samples {
		_, err := GetUserCount(sample)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestAppsGetter(t *testing.T) {
	_, err := GetApps()
	if err != nil {
		t.Error(err)
	}
}

var result error

func BenchmarkUserCountGetter(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = GetUserCount(0)
	}
	result = err
}

func BenchmarkAppsGetter(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = GetApps()
	}
	result = err
}
