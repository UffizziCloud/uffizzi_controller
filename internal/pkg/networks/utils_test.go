package networks

import (
	"testing"

	networkingV1 "k8s.io/api/networking/v1"
)

func TestGetIngresEntrypoint(t *testing.T) {
	type args struct {
		ingress networkingV1.IngressLoadBalancerIngress
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get hostname",
			args: args{
				ingress: networkingV1.IngressLoadBalancerIngress{
					Hostname: "example.com",
				},
			},
			want: "example.com",
		},
		{
			name: "get IP",
			args: args{
				ingress: networkingV1.IngressLoadBalancerIngress{
					IP: "8.8.8.8",
				},
			},
			want: "8.8.8.8",
		},
		{
			name: "get entrypoint",
			args: args{
				ingress: networkingV1.IngressLoadBalancerIngress{
					Hostname: "example.com",
					IP:       "8.8.8.8",
				},
			},
			want: "example.com",
		},
		{
			name: "get empty entrypoint",
			args: args{
				ingress: networkingV1.IngressLoadBalancerIngress{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if got := GetIngresEntrypoint(tt.args.ingress); got != tt.want {
				t.Errorf("GetIngresEntrypoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
