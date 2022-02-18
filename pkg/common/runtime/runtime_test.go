package runtime

import "testing"

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
			name: "Get testing function name",
			args: args{fn: TestGetFuncModule},
			want: "runtime",
		},
		{
			name: "Get anonymous function name",
			args: args{fn: func() {}},
			want: "runtime",
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
