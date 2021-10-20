package networks

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestGetIngresEntrypoint(t *testing.T) {
	type args struct {
		ingress corev1.LoadBalancerIngress
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get hostname",
			args: args{
				ingress: corev1.LoadBalancerIngress{
					Hostname: "example.com",
				},
			},
			want: "example.com",
		},
		{
			name: "get IP",
			args: args{
				ingress: corev1.LoadBalancerIngress{
					IP: "8.8.8.8",
				},
			},
			want: "8.8.8.8",
		},
		{
			name: "get entrypoint",
			args: args{
				ingress: corev1.LoadBalancerIngress{
					Hostname: "example.com",
					IP:       "8.8.8.8",
				},
			},
			want: "example.com",
		},
		{
			name: "get empty entrypoint",
			args: args{
				ingress: corev1.LoadBalancerIngress{},
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
