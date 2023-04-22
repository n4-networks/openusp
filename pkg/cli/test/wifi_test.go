// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clitest

import (
	"testing"
	"time"
)

func TestWiFi_ShowWiFi(t *testing.T) {
	cmd := "show wifi"
	runAndCheckErr(t, cmd)

	cmd = "show wifi ssid"
	runAndCheckErr(t, cmd)

	cmd = "show wifi accesspoint"
	runAndCheckErr(t, cmd)
}

func TestWiFi_AddShowRemoveSSID(t *testing.T) {
	t.Log("Adding a new SSID")
	cmd := "add wifi ssid alias testSsid name myssid radio 1"
	runAndCheckErr(t, cmd)

	time.Sleep(time.Second * 2)

	t.Log("Show wifi ssid")
	cmd = "show wifi ssid"
	runAndCheckErr(t, cmd)

	t.Log("Removing the added ssid")
	cmd = "remove wifi ssid testSsid"
	runAndCheckErr(t, cmd)
}

func TestWiFi_AddShowRemoveAccessPoint(t *testing.T) {
	t.Log("Adding a new AccessPoint")
	cmd := "add wifi ssid alias testSsid name myssid radio 1"
	runAndCheckErr(t, cmd)

	time.Sleep(time.Second * 2)

	t.Log("Show wifi ssid")
	cmd = "show wifi ssid"
	runAndCheckErr(t, cmd)

	ssidPath, err := getInstancePathByAlias("testSsid")
	if err != nil {
		t.Error("Error in locating testSsid in db")
	}
	// Add AccessPoint
	cmd = "add wifi accesspoint alias testAP ssid " + ssidPath + " security open"
	runAndCheckErr(t, cmd)

	t.Log("Removing the added accesspoint")
	cmd = "remove wifi accesspoint testAP"
	runAndCheckErr(t, cmd)

	t.Log("Removing the added ssid")
	cmd = "remove wifi ssid testSsid"
	runAndCheckErr(t, cmd)
}
