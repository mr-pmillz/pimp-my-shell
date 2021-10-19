package linux

import "testing"

func TestCustomTilixBookmarks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Test CustomTilixBookmarks", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CustomTilixBookmarks(); (err != nil) != tt.wantErr {
				t.Errorf("CustomTilixBookmarks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
