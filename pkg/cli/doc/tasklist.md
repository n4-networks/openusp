- [X] Delete Wireless network (delete SSID and AccessPoint instances)
- [X] Change get to show 
- [X] Dynamic auto generated incremental msg id implementation
- [X] show wifi ssid, radio, accesspoint implementation
- [X] Implementation of reconnect to MTP command
- [X] Adding set commands for WiFi
- [X] Add and removal of WiFi network
- [X] Generic function to parse show command line arguments, fetch from param, dm and instance tables and then print
- [X] msgId mismatch problem
- [X] show eth interface is not working
- [X] show ip intf is not showing newly created interface
- [ ] extended AP environment with ingress and egress as WiFi
- [ ] IP QoS management
- [ ] CLI to talk to Http instead of database directly
- [ ] While creating SSID check if ssid object with the same name and radio is there, if so use the same
- [ ] IP Implementation (intf, port, addr)
      [X] show ip intf, port, addr
      [X] add ip intf, addr
      [X] remove ip intf, addr, port
      [X] set ip intf, addr, port 
	  [X] update ip
      [ ] add ip port
- [X ] Wifi refactoring
      [X] show wifi
      [X] add wifi
	  [X] remove wifi
	  [X] set wifi
	  [X] update wifi
- [ ] Bridging 
      [X] show bridging bridge | port
      [X] add bridging bridge  | port
	  [X] remove bridging bridge | port
	  [ ] show bridging filter
	  [ ] add bridging filter
	  [ ] remove bridging filter
	  [X] set bridging
	  [X] update bridging
- [X] Ethernet
      [X] show eth intf
      [X] show eth link
      [X] show eth intf stats
	  [X] set param
	  [X] update param
- [X] Time 
      [X] show time
	  [X] set param
	  [X] update param
- [X] DeviceInfo 
      [X] show deviceinfo 
	  [X] set param
	  [X] update param
- [ ] After receiving boot notify, push required configuration to agent
