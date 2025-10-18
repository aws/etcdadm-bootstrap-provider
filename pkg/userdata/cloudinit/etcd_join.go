package cloudinit

import (
	"fmt"

	etcdbootstrapv1 "github.com/aws/etcdadm-bootstrap-provider/api/v1beta1"
	"github.com/aws/etcdadm-bootstrap-provider/pkg/userdata"
	"github.com/pkg/errors"
)

const (
	etcdPlaneJoinCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
-   path: /run/cluster-api/placeholder
    owner: root:root
    permissions: '0640'
    content: "This placeholder file is used to create the /run/cluster-api sub directory in a way that is compatible with both Linux and Windows (mkdir -p /run/cluster-api does not work with Windows)"
runcmd:
{{- template "commands" .PreEtcdadmCommands }}
  - {{ .EtcdadmJoinCommand }} && {{ .SentinelFileCommand }}
{{- template "commands" .PostEtcdadmCommands }}
{{- template "ntp" .NTP }}
{{- template "users" .Users }}
{{- template "disk_setup" .DiskSetup}}
{{- template "fs_setup" .DiskSetup}}
{{- template "mounts" .Mounts}}
`
)

// NewJoinControlPlane returns the user data string to be used on a new control plane instance.
func NewJoinEtcdPlane(input *userdata.EtcdPlaneJoinInput, config etcdbootstrapv1.EtcdadmConfigSpec) ([]byte, error) {
	input.WriteFiles = userdata.ConvertCertificateFiles(input.Certificates.AsFiles())
	input.EtcdadmArgs = buildEtcdadmArgs(config)
	input.EtcdadmJoinCommand = userdata.AddSystemdArgsToCommand(fmt.Sprintf(standardJoinCommand, input.JoinAddress), &input.EtcdadmArgs)
	if err := setProxy(config.Proxy, &input.BaseUserData); err != nil {
		return nil, err
	}
	if err := setRegistryMirror(config.RegistryMirror, &input.BaseUserData); err != nil {
		return nil, err
	}
	if err := prepare(&input.BaseUserData); err != nil {
		return nil, err
	}
	userData, err := generate("JoinEtcdCluster", etcdPlaneJoinCloudInit, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate user data for machine joining etcd cluster")
	}

	return userData, err
}
