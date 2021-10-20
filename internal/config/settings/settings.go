package settings

import (
	"time"

	ResourceNameUtils "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/resource_name_utils"
)

type Settings struct {
	ResourceName                          ResourceNameUtils.ResouceNameUtils
	NamespaceNamePrefix                   string                `yaml:"namespace_name_prefix"`
	ManagedApplication                    string                `yaml:"managed_application"`
	ControllerPort                        string                `yaml:"controller_port"`
	ControllerLogin                       string                `yaml:"controller_login"`
	ControllerPassword                    string                `yaml:"controller_password"`
	MaxCountProcessDeploymentJobs         int                   `yaml:"max_count_process_deployment_jobs"`
	CustomerDefaultReplicationFactor      int32                 `yaml:"customer_default_replication_factor"`
	CustomerProductionReplicationFactor   int32                 `yaml:"customer_production_replication_factor"`
	ServiceChecks                         ServiceChecksSettings `yaml:"service_checks"`
	ResourceRequestBackOffPeriod          time.Duration         `yaml:"resource_request_back_off_period"`
	CountDisplayedEntriesForLogsOutput    int64                 `yaml:"count_displayed_entries_for_logs_output"`
	PodCidr                               string                `yaml:"pod_cidr"`
	PoolMachineTotalCpuMillicores         int64                 `yaml:"pool_machine_total_cpu_millicores"`
	PoolMachineTotalMemoryBytes           int64                 `yaml:"pool_machine_total_memory_bytes"`
	DefaultAutoscalingCpuThreshold        int                   `yaml:"default_autoscaling_cpu_threshold"`
	DefaultAutoscalingCpuThresholdEpsilon int                   `yaml:"default_autoscaling_cpu_threshold_epsilon"`
	AutoscalingMaxPerformanceReplicas     int32                 `yaml:"autoscaling_max_performance_replicas"`
	AutoscalingMinPerformanceReplicas     int32                 `yaml:"autoscaling_min_performance_replicas"`
	AutoscalingMaxEnterpriseReplicas      int32                 `yaml:"autoscaling_max_enterprise_replicas"`
	AutoscalingMinEnterpriseReplicas      int32                 `yaml:"autoscaling_min_enterprise_replicas"`
	StartupProbeDelaySeconds              int32                 `yaml:"startup_probe_delay_seconds"`
	StartupProbeFailureThreshold          int32                 `yaml:"startup_probe_failure_threshold"`
	StartupProbePeriodSettings            int32                 `yaml:"startup_probe_period_seconds"`
	CertManagerClusterIssuer              string                `yaml:"cert_manager_cluster_issuer"`
	EphemeralStorageCoefficient           float64               `yaml:"ephemeral_storage_coefficient"`
	IngressDefaultPort                    int                   `yaml:"ingress_default_port"`
}

type ServiceChecksSettings struct {
	IPTimeout           time.Duration `yaml:"ip_timeout"`
	IPPingTimeout       time.Duration `yaml:"ip_ping_timeout"`
	AvailabilityTimeout time.Duration `yaml:"availability_timeout"`
	AwaitStatusTimeout  time.Duration `yaml:"await_status_timeout"`
	PerAddressAttempts  uint          `yaml:"per_address_attempts"`
	PerAddressTimeout   time.Duration `yaml:"per_address_timeout"`
}
