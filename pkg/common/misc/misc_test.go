package runtime

import (
	"testing"
)

func TestGetFuncModule(t *testing.T) {
	type args struct {
		fn interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get testing function module name",
			args: args{fn: TestGetFuncModule},
			want: "misc",
		},
		{
			name: "Get anonymous function module name",
			args: args{fn: func() {}},
			want: "misc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFuncModule(tt.args.fn); got != tt.want {
				t.Errorf("GetFuncModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetParentFuncName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Get testing function name",
			want: "testing.tRunner",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetParentFuncName(); got != tt.want {
				t.Errorf("GetParentFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}
