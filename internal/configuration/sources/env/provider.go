package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
)

func (r *Reader) readProvider(vpnType string) (provider settings.Provider, err error) {
	provider.Name = r.readVPNServiceProvider(vpnType)
	var providerName string
	if provider.Name != nil {
		providerName = *provider.Name
	}

	provider.ServerSelection, err = r.readServerSelection(providerName, vpnType)
	if err != nil {
		return provider, fmt.Errorf("server selection: %w", err)
	}

	provider.PortForwarding, err = r.readPortForward()
	if err != nil {
		return provider, fmt.Errorf("port forwarding: %w", err)
	}

	return provider, nil
}

func (r *Reader) readVPNServiceProvider(vpnType string) (vpnProviderPtr *string) {
	_, s := r.getEnvWithRetro("VPN_SERVICE_PROVIDER", "VPNSP")
	s = strings.ToLower(s)
	switch {
	case vpnType != vpn.Wireguard &&
		os.Getenv("OPENVPN_CUSTOM_CONFIG") != "": // retro compatibility
		return stringPtr(providers.Custom)
	case s == "":
		return nil
	case s == "pia": // retro compatibility
		return stringPtr(providers.PrivateInternetAccess)
	}
	return stringPtr(s)
}
