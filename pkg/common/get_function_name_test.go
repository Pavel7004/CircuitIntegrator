package common

import "testing"

func TestGetFuncName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Get this func name",
			want: "common.TestGetFuncName.func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFuncName(); got != tt.want {
				t.Errorf("GetFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}
