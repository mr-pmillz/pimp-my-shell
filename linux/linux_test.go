package linux

import "testing"

func TestCustomTilixBookmarks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Test CustomTerminalBookmarks", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CustomTerminalBookmarks(); (err != nil) != tt.wantErr {
				t.Errorf("CustomTerminalBookmarks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
