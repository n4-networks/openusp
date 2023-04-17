package cli

type node struct {
	pathName      string
	multiInstance bool
}

var nodeSchema = map[string]node{
	//DeviceInfo
	"device": {"Device", false},

	"devinfo":              {"DeviceInfo", false},
	"devinfo.vendorcfg":    {"VendorConfigFile", true},
	"devinfo.memory":       {"MemoryStatus", false},
	"devinfo.process":      {"ProcessStatus", false},
	"devinfo.process.proc": {"Process", true},
	"devinfo.temp":         {"TemperatureStatus", false},
	"devinfo.temp.sensor":  {"TemperatureSensor", true},
	"devinfo.net":          {"NetworkProperties", false},
	"devinfo.cpu":          {"Processor", true},
	"devinfo.logfile":      {"VendorLogFile", true},
	"devinfo.loc":          {"Location", true},
	"devinfo.image":        {"DeviceImageFile", true},
	"devinfo.firmware":     {"FirmwareImage", true},

	// Time
	"time": {"Time", false},

	// IP
	"ip":               {"IP", false},
	"ip.intf":          {"Interface", true},
	"ip.intf.stats":    {"Stats", false},
	"ip.intf.ipv4addr": {"IPv4Address", true},
	"ip.intf.ipv6addr": {"IPv6Address", true},
	"ip.acport":        {"ActivePort", true},
	// WiFi
	"wifi":                        {"WiFi", false},
	"wifi.ssid":                   {"SSID", true},
	"wifi.ssid.stats":             {"Stats", false},
	"wifi.radio":                  {"Radio", true},
	"wifi.radio.stats":            {"Stats", false},
	"wifi.accesspoint":            {"AccessPoint", true},
	"wifi.accesspoint.security":   {"Security", false},
	"wifi.accesspoint.wps":        {"WPS", false},
	"wifi.accesspoint.stations":   {"AssociatedDevice", true},
	"wifi.accesspoint.qos":        {"AC", true},
	"wifi.accesspoint.accounting": {"Accounting", false},
	"wifi.endpoint":               {"EndPoint", true},
	"wifi.endpoint.stats":         {"Stats", true},
	"wifi.endpoint.profile":       {"Profile", true},
	// Ethernet
	"eth":                {"Ethernet", false},
	"eth.intf":           {"Interface", true},
	"eth.intf.stats":     {"Stats", false},
	"eth.link":           {"Link", true},
	"eth.vlanterm":       {"VLANTermination", true},
	"eth.vlanterm.stats": {"Stats", false},
	"eth.rmonstats":      {"RMONStats", false},
	"eth.wol":            {"WoL", false},
	"eth.lag":            {"LAG", true},
	"eth.lag.stats":      {"Stats", false},
	// Bridging
	"bridging":                           {"Bridging", false},
	"bridging.bridge":                    {"Bridge", true},
	"bridging.bridge.port":               {"Port", true},
	"bridging.bridge.port.stats":         {"Port", true},
	"bridging.bridge.port.priocodepoint": {"PriorityCodePoint", false},
	"bridging.vlan":                      {"VLAN", true},
	"bridging.vlanport":                  {"VLANPort", true},
	"bridging.filter":                    {"Filter", true},
	"bridging.providerbridge":            {"ProviderBridge", true},

	// DHCPv4 Server
	"dhcpv4":                             {"DHCPv4", false},
	"dhcpv4.server":                      {"Server", false},
	"dhcpv4.server.client":               {"Client", true},
	"dhcpv4.server.pool":                 {"Pool", true},
	"dhcpv4.server.pool.staticaddr":      {"StaticAddress", true},
	"dhcpv4.server.pool.option":          {"Option", true},
	"dhcpv4.server.pool.client":          {"Client", true},
	"dhcpv4.server.pool.client.ipv4addr": {"IPv4Address", true},
	"dhcpv4.server.pool.client.option":   {"Option", true},

	// DHCPv4 Relay
	"dhcpv4.relay":            {"Relay", false},
	"dhcpv4.relay.forwarding": {"Forwarding", true},

	// NAT
	"nat":              {"NAT", false},
	"nat.intf-setting": {"InterfaceSetting", true},
	"nat.port-mapping": {"PortMapping", true},

	// LocalAgent
	"agent":                                  {"LocalAgent", false},
	"agent.mtp":                              {"MTP", true},
	"agent.mtp.coap":                         {"CoAP", false},
	"agent.mtp.stomp":                        {"STOMP", false},
	"agent.mtp.websocket":                    {"WebSocket", false},
	"agent.mtp.mqtt":                         {"MQTT", false},
	"agent.threshold":                        {"Threshold", true},
	"agent.controller":                       {"Controller", true},
	"agent.controller.mtp":                   {"MTP", true},
	"agent.controller.mtp.coap":              {"CoAP", false},
	"agent.controller.mtp.stomp":             {"STOMP", false},
	"agent.controller.mtp.websocket":         {"WebSocket", false},
	"agent.controller.mtp.mqtt":              {"MQTT", false},
	"agent.controller.boot-params":           {"BootParameter", true},
	"agent.controller.e2e":                   {"E2ESession", false},
	"agent.subscription":                     {"Subscription", true},
	"agent.request":                          {"Request", true},
	"agent.cert":                             {"Certificate", true},
	"agent.controller-trust":                 {"ControllerTrust", false},
	"agent.controller-trust.role":            {"Role", true},
	"agent.controller-trust.role.permission": {"Permission", true},
	"agent.controller-trust.creds":           {"Credential", true},
	"agent.controller-trust.challenge":       {"Challenge", true},
}

var cmdSchema = map[string]string{
	"reboot":         "Reboot()",
	"factory-reset":  "FactoryReset()",
	"self-test-diag": "SelfTestDiagnostic()",
	"download":       "Download()",
}
