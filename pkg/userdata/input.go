package userdata

import (
	"fmt"
	"strings"

	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta1"
	bootstrapv2 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	"sigs.k8s.io/cluster-api/util/secret"
)

// EtcdPlaneInput defines the context to generate etcd instance user data for initializing etcd cluster.
type EtcdPlaneInput struct {
	BaseUserData
	secret.Certificates
	EtcdadmArgs

	EtcdadmInitCommand string
}

// EtcdPlaneJoinInput defines context to generate etcd instance user data for etcd plane node join.
type EtcdPlaneJoinInput struct {
	BaseUserData
	secret.Certificates
	EtcdadmArgs

	EtcdadmJoinCommand string
	JoinAddress        string
}

// BaseUserData is shared across all the various types of files written to disk.
type BaseUserData struct {
	Header              string
	PreEtcdadmCommands  []string
	PostEtcdadmCommands []string
	AdditionalFiles     []bootstrapv1.File
	WriteFiles          []bootstrapv1.File
	Users               []bootstrapv1.User
	NTP                 *bootstrapv1.NTP
	DiskSetup           *bootstrapv1.DiskSetup
	Mounts              []bootstrapv1.MountPoints
	ControlPlane        bool
	SentinelFileCommand string
	Hostname            string
	RegistryMirrorCredentials
}

type EtcdadmArgs struct {
	Version         string
	ImageRepository string
	EtcdReleaseURL  string
	InstallDir      string
	CipherSuites    string
}

type RegistryMirrorCredentials struct {
	Username string
	Password string
}

func (args *EtcdadmArgs) SystemdFlags() []string {
	flags := make([]string, 0, 3)
	flags = append(flags, "--init-system systemd")
	if args.Version != "" {
		flags = append(flags, fmt.Sprintf("--version %s", args.Version))
	}
	if args.EtcdReleaseURL != "" {
		flags = append(flags, fmt.Sprintf("--release-url %s", args.EtcdReleaseURL))
	}
	if args.InstallDir != "" {
		flags = append(flags, fmt.Sprintf("--install-dir %s", args.InstallDir))
	}
	if args.CipherSuites != "" {
		flags = append(flags, fmt.Sprintf("--cipher-suites %s", args.CipherSuites))
	}
	return flags
}

func AddSystemdArgsToCommand(cmd string, args *EtcdadmArgs) string {
	flags := args.SystemdFlags()
	fullCommand := make([]string, len(flags)+1)
	fullCommand = append(fullCommand, cmd)
	fullCommand = append(fullCommand, flags...)

	return strings.Join(fullCommand, " ")
}

// ConvertCertificateFiles converts v1beta2.File slice to v1beta1.File slice using cluster-api conversion
func ConvertCertificateFiles(v2Files []bootstrapv2.File) []bootstrapv1.File {
	v1Files := make([]bootstrapv1.File, len(v2Files))
	for i, v2File := range v2Files {
		// Use cluster-api's built-in conversion function
		if err := bootstrapv1.Convert_v1beta2_File_To_v1beta1_File(&v2File, &v1Files[i], nil); err != nil {
			// Fallback to manual conversion if the built-in conversion fails
			v1Files[i] = bootstrapv1.File{
				Path:        v2File.Path,
				Owner:       v2File.Owner,
				Permissions: v2File.Permissions,
				Encoding:    bootstrapv1.Encoding(v2File.Encoding),
				Content:     v2File.Content,
				// ContentFrom conversion is handled by the built-in conversion function
			}
		}
	}
	return v1Files
}
