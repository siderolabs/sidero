// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package siderolink

import (
	"context"
	"fmt"
	"net/netip"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterctl "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/siderolabs/siderolink/pkg/wireguard"
)

// Config is the internal SideroLink configuration.
type Config struct {
	InstallationID    string
	PrivateKey        wgtypes.Key
	PublicKey         wgtypes.Key
	WireguardEndpoint string
	Subnet            netip.Prefix
	ServerAddress     netip.Prefix
}

const (
	secretInstallationID = "installation-id"
	secretPrivateKey     = "private-key"
)

func (cfg *Config) LoadOrCreate(ctx context.Context, metalClient runtimeclient.Client) error {
	var secret corev1.Secret

retry:
	err := metalClient.Get(ctx, types.NamespacedName{Namespace: corev1.NamespaceDefault, Name: SecretName}, &secret)

	if err == nil {
		if err = cfg.loadFrom(&secret); err != nil {
			return fmt.Errorf("error loading from secret")
		}

		return nil
	}

	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("error fetching secret: %w", err)
	}

	if err = cfg.generate(); err != nil {
		return fmt.Errorf("error generating config: %w", err)
	}

	if err = cfg.save(ctx, metalClient); err != nil {
		if apierrors.IsAlreadyExists(err) {
			// config was already saved by another process, retry loading
			goto retry
		}

		return fmt.Errorf("error saving config: %w", err)
	}

	return nil
}

func (cfg *Config) loadFrom(secret *corev1.Secret) error {
	if b, ok := secret.Data[secretInstallationID]; !ok {
		return fmt.Errorf("missing %q key", secretInstallationID)
	} else {
		cfg.InstallationID = string(b)
	}

	if b, ok := secret.Data[secretPrivateKey]; !ok {
		return fmt.Errorf("missing %q key", secretPrivateKey)
	} else {
		var err error

		cfg.PrivateKey, err = wgtypes.ParseKey(string(b))
		if err != nil {
			return fmt.Errorf("error parsing key: %w", err)
		}
	}

	cfg.fill()

	return nil
}

func (cfg *Config) save(ctx context.Context, metalClient runtimeclient.Client) error {
	return metalClient.Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      SecretName,
			Labels: map[string]string{
				clusterctl.ClusterctlMoveLabel: "",
			},
		},
		Data: map[string][]byte{
			secretInstallationID: []byte(cfg.InstallationID),
			secretPrivateKey:     []byte(cfg.PrivateKey.String()),
		},
	})
}

func (cfg *Config) generate() error {
	installID, err := wgtypes.GeneratePrivateKey() // use private wireguard key as the installation ID, as it's random bytes
	if err != nil {
		return err
	}

	cfg.InstallationID = installID.String()

	cfg.PrivateKey, err = wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}

	cfg.fill()

	return nil
}

func (cfg *Config) fill() {
	cfg.PublicKey = cfg.PrivateKey.PublicKey()

	cfg.Subnet = wireguard.NetworkPrefix(cfg.InstallationID)
	cfg.ServerAddress = netip.PrefixFrom(cfg.Subnet.Addr().Next(), cfg.Subnet.Bits())
}
