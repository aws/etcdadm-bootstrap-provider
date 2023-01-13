package bottlerocket

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"
	"text/template"

	etcdbootstrapv1 "github.com/aws/etcdadm-bootstrap-provider/api/v1beta1"
	"github.com/aws/etcdadm-bootstrap-provider/pkg/userdata"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
)

const (
	hostContainersTemplate = `{{ define "hostContainersSettings" -}}
{{- range .HostContainers }}
[settings.host-containers.{{ .Name }}]
enabled = true
superpowered = {{ .Superpowered }}
{{- if .Image }}
source = "{{ .Image }}"
{{- end }}
{{- if .UserData }}
user-data = "{{ .UserData }}"
{{- end }}
{{- end }}
{{- end }}
`

	bootstrapContainersTemplate = `{{ define "bootstrapContainersSettings" -}}
{{- range .BootstrapContainers }}
[settings.bootstrap-containers.{{ .Name }}]
essential = {{ .Essential }}
mode = "{{ .Mode }}"
{{- if .Image }}
source = "{{ .Image }}"
{{- end }}
{{- if .UserData }}
user-data = "{{ .UserData }}"
{{- end }}
{{- end }}
{{- end }}
`

	kubernetesInitTemplate = `{{ define "kubernetesInitSettings" -}}
[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "{{.PauseContainerSource}}"
{{- end -}}
`

	networkInitTemplate = `{{ define "networkInitSettings" -}}
[settings.network]
https-proxy = "{{.HTTPSProxyEndpoint}}"
no-proxy = [{{stringsJoin .NoProxyEndpoints "," }}]
{{- end -}}
`
	registryMirrorTemplate = `{{ define "registryMirrorSettings" -}}
[settings.container-registry.mirrors]
"public.ecr.aws" = ["https://{{.RegistryMirrorEndpoint}}"]
{{- end -}}
`
	registryMirrorCACertTemplate = `{{ define "registryMirrorCACertSettings" -}}
[settings.pki.registry-mirror-ca]
data = "{{.RegistryMirrorCACert}}"
trusted=true
{{- end -}}
`
	registryMirrorCredentialsTemplate = `{{ define "registryMirrorCredentialsSettings" -}}
[[settings.container-registry.credentials]]
registry = "public.ecr.aws"
username = "{{.RegistryMirrorUsername}}"
password = "{{.RegistryMirrorPassword}}"
[[settings.container-registry.credentials]]
registry = "{{.RegistryMirrorEndpoint}}"
username = "{{.RegistryMirrorUsername}}"
password = "{{.RegistryMirrorPassword}}"
{{- end -}}
`
	bottlerocketNodeInitSettingsTemplate = `{{template "hostContainersSettings" .}}

{{template "kubernetesInitSettings" .}}

{{- if .BootstrapContainers }}
{{template "bootstrapContainersSettings" .}}
{{- end -}}

{{- if (ne .HTTPSProxyEndpoint "")}}
{{template "networkInitSettings" .}}
{{- end -}}

{{- if (ne .RegistryMirrorEndpoint "")}}
{{template "registryMirrorSettings" .}}
{{- end -}}

{{- if (ne .RegistryMirrorCACert "")}}
{{template "registryMirrorCACertSettings" .}}
{{- end -}}

{{- if and (ne .RegistryMirrorUsername "") (ne .RegistryMirrorPassword "")}}
{{template "registryMirrorCredentialsSettings" .}}
{{- end -}}
`
)

type bottlerocketSettingsInput struct {
	PauseContainerSource   string
	HTTPSProxyEndpoint     string
	NoProxyEndpoints       []string
	RegistryMirrorEndpoint string
	RegistryMirrorCACert   string
	RegistryMirrorUsername string
	RegistryMirrorPassword string
	HostContainers         []etcdbootstrapv1.BottlerocketHostContainer
	BootstrapContainers    []etcdbootstrapv1.BottlerocketBootstrapContainer
}

// generateBottlerocketNodeUserData returns the userdata for the host bottlerocket in toml format
func generateBottlerocketNodeUserData(kubeadmBootstrapContainerUserData []byte, users []bootstrapv1.User, registryMirrorCredentials userdata.RegistryMirrorCredentials, config etcdbootstrapv1.EtcdadmConfigSpec, log logr.Logger) ([]byte, error) {
	// base64 encode the kubeadm bootstrapContainer's user data
	b64KubeadmBootstrapContainerUserData := base64.StdEncoding.EncodeToString(kubeadmBootstrapContainerUserData)

	// Parse out all the ssh authorized keys
	sshAuthorizedKeys := getAllAuthorizedKeys(users)

	// generate the userdata for the admin container
	adminContainerUserData, err := generateAdminContainerUserData("InitAdminContainer", usersTemplate, sshAuthorizedKeys)
	if err != nil {
		return nil, err
	}
	b64AdminContainerUserData := base64.StdEncoding.EncodeToString(adminContainerUserData)

	hostContainers := []etcdbootstrapv1.BottlerocketHostContainer{
		{
			Name:         "admin",
			Superpowered: true,
			Image:        config.BottlerocketConfig.AdminImage,
			UserData:     b64AdminContainerUserData,
		},
		{
			Name:         "kubeadm-bootstrap",
			Superpowered: true,
			Image:        config.BottlerocketConfig.BootstrapImage,
			UserData:     b64KubeadmBootstrapContainerUserData,
		},
	}

	if config.BottlerocketConfig.ControlImage != "" {
		hostContainers = append(hostContainers, etcdbootstrapv1.BottlerocketHostContainer{
			Name:         "control",
			Superpowered: false,
			Image:        config.BottlerocketConfig.ControlImage,
		})
	}

	bottlerocketInput := &bottlerocketSettingsInput{
		PauseContainerSource: config.BottlerocketConfig.PauseImage,
		HostContainers:       hostContainers,
		BootstrapContainers:  config.BottlerocketConfig.CustomBootstrapContainers,
	}

	if config.Proxy != nil {
		bottlerocketInput.HTTPSProxyEndpoint = config.Proxy.HTTPSProxy
		for _, noProxy := range config.Proxy.NoProxy {
			bottlerocketInput.NoProxyEndpoints = append(bottlerocketInput.NoProxyEndpoints, strconv.Quote(noProxy))
		}
	}

	if config.RegistryMirror != nil {
		bottlerocketInput.RegistryMirrorEndpoint = config.RegistryMirror.Endpoint
		if config.RegistryMirror.CACert != "" {
			bottlerocketInput.RegistryMirrorCACert = base64.StdEncoding.EncodeToString([]byte(config.RegistryMirror.CACert))
		}
		bottlerocketInput.RegistryMirrorUsername = registryMirrorCredentials.Username
		bottlerocketInput.RegistryMirrorPassword = registryMirrorCredentials.Password
	}

	bottlerocketNodeUserData, err := generateNodeUserData("InitBottlerocketNode", bottlerocketNodeInitSettingsTemplate, bottlerocketInput)
	if err != nil {
		return nil, err
	}
	log.Info("Generated bottlerocket bootstrap userdata", "bootstrapContainerImage", config.BottlerocketConfig.BootstrapImage)
	return bottlerocketNodeUserData, nil
}

// getAllAuthorizedKeys parses through all the users and return list of all user's authorized ssh keys
func getAllAuthorizedKeys(users []bootstrapv1.User) string {
	var sshAuthorizedKeys []string
	for _, user := range users {
		if len(user.SSHAuthorizedKeys) != 0 {
			for _, key := range user.SSHAuthorizedKeys {
				quotedKey := "\"" + key + "\""
				sshAuthorizedKeys = append(sshAuthorizedKeys, quotedKey)
			}
		}
	}
	return strings.Join(sshAuthorizedKeys, ",")
}

func generateAdminContainerUserData(kind string, tpl string, data interface{}) ([]byte, error) {
	tm := template.New(kind)
	if _, err := tm.Parse(usersTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse users - %s template", kind)
	}
	t, err := tm.Parse(tpl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s template", kind)
	}
	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return nil, errors.Wrapf(err, "failed to generate %s template", kind)
	}
	return out.Bytes(), nil
}

func generateNodeUserData(kind string, tpl string, data interface{}) ([]byte, error) {
	tm := template.New(kind).Funcs(template.FuncMap{"stringsJoin": strings.Join})
	if _, err := tm.Parse(hostContainersTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse hostContainers %s template", kind)
	}
	if _, err := tm.Parse(bootstrapContainersTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse bootstrapContainers %s template", kind)
	}
	if _, err := tm.Parse(kubernetesInitTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse kubernetes %s template", kind)
	}
	if _, err := tm.Parse(networkInitTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse networks %s template", kind)
	}
	if _, err := tm.Parse(registryMirrorTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse registry mirror %s template", kind)
	}
	if _, err := tm.Parse(registryMirrorCACertTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse registry mirror ca cert %s template", kind)
	}
	if _, err := tm.Parse(registryMirrorCredentialsTemplate); err != nil {
		return nil, errors.Wrapf(err, "failed to parse registry mirror credentials %s template", kind)
	}

	t, err := tm.Parse(tpl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s template", kind)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return nil, errors.Wrapf(err, "failed to generate %s template", kind)
	}
	return out.Bytes(), nil
}
