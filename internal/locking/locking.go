package locking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apicorev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	semaphoreInformationKey = "lock-information"
	eksaSystemNamespace     = "eksa-system"
)

// EtcdadmInitMutex uses a ConfigMap to synchronize cluster initialization.
type EtcdadmInitMutex struct {
	log    logr.Logger
	client client.Client
}

// NewControlPlaneInitMutex returns a lock that can be held by a control plane node before init.
func NewEtcdadmInitMutex(log logr.Logger, client client.Client) *EtcdadmInitMutex {
	return &EtcdadmInitMutex{
		log:    log,
		client: client,
	}
}

// Lock allows a control plane node to be the first and only node to run kubeadm init
func (c *EtcdadmInitMutex) Lock(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) bool {
	sema := newSemaphore()
	cmName := configMapName(cluster.Name)
	log := c.log.WithValues("namespace", cluster.Namespace, "cluster-name", cluster.Name, "configmap-name", cmName, "machine-name", machine.Name)
	err := c.client.Get(ctx, client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cmName,
	}, sema.ConfigMap)
	switch {
	case apierrors.IsNotFound(err):
		break
	case err != nil:
		log.Error(err, "Failed to acquire lock")
		return false
	default: // successfully found an existing config map
		info, err := sema.information()
		if err != nil {
			log.Error(err, "Failed to get information about the existing lock")
			return false
		}
		// the machine requesting the lock is the machine that created the lock, therefore the lock is acquired
		if info.MachineName == machine.Name {
			return true
		}

		machine := &clusterv1.Machine{}

		err = c.client.Get(ctx, client.ObjectKey{
			Namespace: eksaSystemNamespace,
			Name:      info.MachineName,
		}, machine)
		if err != nil {
			// Release the lock if for some reason the machine that acquired the lock
			// failed to launch due to a catastrophic event
			//
			// without this check we might end up with a deadlock.
			if apierrors.IsNotFound(err) {
				log.Info("Machine that has acquired the lock not found, releasing the lock", "init-machine", info.MachineName)
				if c.Unlock(ctx, cluster) {
					break
				} else {
					return false
				}
			} else {
				log.Error(err, "Failed to retreive machine", "machine", info.MachineName)
				return false
			}
		}

		log.Info("Waiting on another machine to initialize", "init-machine", info.MachineName)
		return false
	}

	// Adds owner reference, namespace and name
	sema.setMetadata(cluster)
	// Adds the additional information
	if err := sema.setInformation(&information{MachineName: machine.Name}); err != nil {
		log.Error(err, "Failed to acquire lock while setting semaphore information")
		return false
	}

	log.Info("Attempting to acquire the lock")
	err = c.client.Create(ctx, sema.ConfigMap)
	switch {
	case apierrors.IsAlreadyExists(err):
		log.Info("Cannot acquire the lock. The lock has been acquired by someone else")
		return false
	case err != nil:
		log.Error(err, "Error acquiring the lock")
		return false
	default:
		return true
	}
}

// Unlock releases the lock
func (c *EtcdadmInitMutex) Unlock(ctx context.Context, cluster *clusterv1.Cluster) bool {
	sema := newSemaphore()
	cmName := configMapName(cluster.Name)
	log := c.log.WithValues("namespace", cluster.Namespace, "cluster-name", cluster.Name, "configmap-name", cmName)
	log.Info("Checking for lock")
	err := c.client.Get(ctx, client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cmName,
	}, sema.ConfigMap)
	switch {
	case apierrors.IsNotFound(err):
		log.Info("Control plane init lock not found, it may have been released already")
		return true
	case err != nil:
		log.Error(err, "Error unlocking the control plane init lock")
		return false
	default:
		// Delete the config map semaphore if there is no error fetching it
		if err := c.client.Delete(ctx, sema.ConfigMap); err != nil {
			if apierrors.IsNotFound(err) {
				return true
			}
			log.Error(err, "Error deleting the config map underlying the control plane init lock")
			return false
		}
		return true
	}
}

type information struct {
	MachineName string `json:"machineName"`
}

type semaphore struct {
	*apicorev1.ConfigMap
}

func newSemaphore() *semaphore {
	return &semaphore{&apicorev1.ConfigMap{}}
}

func configMapName(clusterName string) string {
	return fmt.Sprintf("%s-etcd-lock", clusterName)
}

func (s semaphore) information() (*information, error) {
	li := &information{}
	if err := json.Unmarshal([]byte(s.Data[semaphoreInformationKey]), li); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal semaphore information")
	}
	return li, nil
}

func (s semaphore) setInformation(information *information) error {
	b, err := json.Marshal(information)
	if err != nil {
		return errors.Wrap(err, "failed to marshal semaphore information")
	}
	s.Data = map[string]string{}
	s.Data[semaphoreInformationKey] = string(b)
	return nil
}

func (s *semaphore) setMetadata(cluster *clusterv1.Cluster) {
	s.ObjectMeta = metav1.ObjectMeta{
		Namespace: cluster.Namespace,
		Name:      configMapName(cluster.Name),
		Labels: map[string]string{
			clusterv1.ClusterNameLabel: cluster.Name,
		},
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: cluster.APIVersion,
				Kind:       cluster.Kind,
				Name:       cluster.Name,
				UID:        cluster.UID,
			},
		},
	}
}
