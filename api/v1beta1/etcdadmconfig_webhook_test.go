/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"testing"

	"github.com/onsi/gomega"
)

func TestEtcdadmConfigDefault(t *testing.T) {
	g := gomega.NewWithT(t)

	config := &EtcdadmConfig{}
	defaulter := &EtcdadmConfigDefaulter{}

	err := defaulter.Default(context.TODO(), config)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestEtcdadmConfigValidateCreate(t *testing.T) {
	g := gomega.NewWithT(t)

	config := &EtcdadmConfig{}
	validator := &EtcdadmConfigValidator{}

	warnings, err := validator.ValidateCreate(context.TODO(), config)
	g.Expect(warnings).To(gomega.BeNil())
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestEtcdadmConfigValidateUpdate(t *testing.T) {
	g := gomega.NewWithT(t)

	oldConfig := &EtcdadmConfig{}
	newConfig := &EtcdadmConfig{}
	validator := &EtcdadmConfigValidator{}

	warnings, err := validator.ValidateUpdate(context.TODO(), oldConfig, newConfig)
	g.Expect(warnings).To(gomega.BeNil())
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestEtcdadmConfigValidateDelete(t *testing.T) {
	g := gomega.NewWithT(t)

	config := &EtcdadmConfig{}
	validator := &EtcdadmConfigValidator{}

	warnings, err := validator.ValidateDelete(context.TODO(), config)
	g.Expect(warnings).To(gomega.BeNil())
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
