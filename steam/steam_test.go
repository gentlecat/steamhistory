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
