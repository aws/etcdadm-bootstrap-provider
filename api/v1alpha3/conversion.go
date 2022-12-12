package v1alpha3

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	etcdv1beta1 "github.com/aws/etcdadm-bootstrap-provider/api/v1beta1"
)

// ConvertTo converts this EtcdadmConfig to the Hub version (v1beta1).
func (src *EtcdadmConfig) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*etcdv1beta1.EtcdadmConfig)
	if err := Convert_v1alpha3_EtcdadmConfig_To_v1beta1_EtcdadmConfig(src, dst, nil); err != nil {
		return err
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this EtcdadmConfig.
func (dst *EtcdadmConfig) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*etcdv1beta1.EtcdadmConfig)
	return Convert_v1beta1_EtcdadmConfig_To_v1alpha3_EtcdadmConfig(src, dst, nil)
}

// ConvertTo converts this EtcdadmConfigList to the Hub version (v1beta1).
func (src *EtcdadmConfigList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*etcdv1beta1.EtcdadmConfigList)
	if err := Convert_v1alpha3_EtcdadmConfigList_To_v1beta1_EtcdadmConfigList(src, dst, nil); err != nil {
		return err
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this EtcdadmConfigList.
func (dst *EtcdadmConfigList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*etcdv1beta1.EtcdadmConfigList)
	return Convert_v1beta1_EtcdadmConfigList_To_v1alpha3_EtcdadmConfigList(src, dst, nil)
}

func Convert_v1beta1_BottlerocketConfig_To_v1alpha3_BottlerocketConfig(in *etcdv1beta1.BottlerocketConfig, out *BottlerocketConfig, s apiconversion.Scope) error {
	return autoConvert_v1beta1_BottlerocketConfig_To_v1alpha3_BottlerocketConfig(in, out, s)
}
