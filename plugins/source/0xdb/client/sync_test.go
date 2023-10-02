package client

import (
	"fmt"
	"testing"
)

const (
	pass  = "\u2713"
	fail  = "\u2717"
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

func Test_calcMissed(t *testing.T) {
	tests := []struct {
		name string
		prev uint64
		next uint64
		want string
	}{
		{
			name: "no missed",
			prev: 0,
			next: 1,
			want: "",
		},
		{
			name: "one missed",
			prev: 0,
			next: 2,
			want: "1",
		},
		{
			name: "few missed",
			prev: 0,
			next: 3,
			want: "1-2",
		},
		{
			name: "many missed",
			prev: 0,
			next: 1000,
			want: "1-999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcMissed(&tt.prev, tt.next); got != tt.want {
				vs := fmt.Sprintf("\t%s difference in got vs want missed"+
					"\nGot: "+red+" \n\n%s\n\n "+reset+"\nWant: "+green+"\n\n%s\n\n"+reset,
					fail, string(got), string(tt.want))
				t.Errorf(vs)
			}
		})
	}
}
