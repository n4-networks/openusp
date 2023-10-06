package cli

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type RestObjParam struct {
	Path   string       `json:"path"`
	Params []*RestParam `json:"params"`
}
type RestParam struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Access string `json:"access"`
}

type RestObjInstance struct {
	EndpointId string            `json:"endpoint_id"`
	Path       string            `json:"path"`
	UniqueKeys map[string]string `json:"unique_keys"`
}

type RestDmCommand struct {
	Name    string   `json:"name"`
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
}
type RestDmEvent struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}
type RestDmParam struct {
	Name   string `json:"name"`
	Access string `json:"access"`
}

type RestObjDm struct {
	EndpointId    string          `json:"endpoint_id"`
	Path          string          `json:"path"`
	MultiInstance bool            `json:"multi_instance"`
	Access        string          `json:"access"`
	Params        []RestDmParam   `json:"params"`
	Events        []RestDmEvent   `json:"events"`
	Cmds          []RestDmCommand `json:"cmds"`
}

const (
	RECONNECT_MTP    = "/reconnect/mtp/"
	RECONNECT_DB     = "/reconnect/db/"
	GET_AGENTS       = "/get/agents/"
	GET_PARAMS       = "/get/params/"
	GET_INSTANCES    = "/get/instances/"
	GET_DM           = "/get/dm/"
	UPDATE_DM        = "/update/dm/"
	DELETE_DBCOLL    = "/delete/dbcoll/"
	UPDATE_PARAMS    = "/update/params/"
	UPDATE_INSTANCES = "/update/instances/"
	DELETE_INSTANCES = "/delete/instances/"
	ADD_INSTANCES    = "/add/instances/"
	SET_PARAMS       = "/set/params/"
	OPERATE_CMD      = "/operate/cmd/"
	GET_MTPINFO      = "/get/mtpinfo/"
)

func (cli *Cli) restInit() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli.rest.client = &http.Client{Transport: tr}
	return nil
}
func (cli *Cli) restReconnectMtp() error {
	url := cli.cfg.apiServerAddr + RECONNECT_MTP

	_, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt reconnect MTP, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restReconnectDb() error {
	url := cli.cfg.apiServerAddr + RECONNECT_DB

	_, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt reconnect DB, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restReadAgents() ([]string, error) {
	url := cli.cfg.apiServerAddr + GET_AGENTS

	data, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return nil, err
	}

	var agents []string
	if err := json.Unmarshal(data, &agents); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}
	return agents, nil
}
func (cli *Cli) restReadParams(path string) ([]*RestObjParam, error) {
	if path == "" {
		log.Println("Err: Blank USP path")
		return nil, errors.New("Blank USP path")
	}
	if !cli.agent.isSet.epId {
		return nil, errors.New("agent endpoint id is not set")
	}
	url := cli.cfg.apiServerAddr + GET_PARAMS + cli.agent.epId + "/" + path

	data, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return nil, err
	}

	var objInstances []*RestObjParam
	if err := json.Unmarshal(data, &objInstances); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}
	return objInstances, nil

}

func (cli *Cli) restReadInstances(path string) ([]*RestObjInstance, error) {
	url := cli.cfg.apiServerAddr + GET_INSTANCES + cli.agent.epId + "/" + path

	data, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return nil, err
	}

	var objInstances []*RestObjInstance
	if err := json.Unmarshal(data, &objInstances); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}
	return objInstances, nil
}
func (cli *Cli) restReadDm(path string) ([]*RestObjDm, error) {
	url := cli.cfg.apiServerAddr + GET_DM + cli.agent.epId + "/" + path

	data, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return nil, err
	}

	var objInstances []*RestObjDm
	if err := json.Unmarshal(data, &objInstances); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}
	return objInstances, nil
}

func (cli *Cli) restUpdateDm(path string) error {
	url := cli.cfg.apiServerAddr + UPDATE_DM + cli.agent.epId + "/" + path

	if _, err := cli.restGet(url); err != nil {
		log.Println("Error in RESt DM Update, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restDeleteCollection(collName string) error {
	url := cli.cfg.apiServerAddr + DELETE_DBCOLL + collName

	if _, err := cli.restGet(url); err != nil {
		log.Println("Error in RESt delete Db collection, err:", err)
		return err
	}
	return nil
}
func (cli *Cli) restUpdateParams(path string) error {
	url := cli.cfg.apiServerAddr + UPDATE_PARAMS + cli.agent.epId + "/" + path

	if _, err := cli.restGet(url); err != nil {
		log.Println("Error in RESt Params Update, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restUpdateInstances(path string) error {
	url := cli.cfg.apiServerAddr + UPDATE_INSTANCES + cli.agent.epId + "/" + path

	if _, err := cli.restGet(url); err != nil {
		log.Println("Error in RESt Instance Update, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restDeleteInstances(path string) error {
	url := cli.cfg.apiServerAddr + DELETE_INSTANCES + cli.agent.epId + "/" + path

	if _, err := cli.restGet(url); err != nil {
		log.Println("Error in RESt Delete Instances, err", err)
		return err
	}
	return nil
}
func (cli *Cli) restAddInstance(path string, params map[string]string) (*Instance, error) {
	url := cli.cfg.apiServerAddr + ADD_INSTANCES + cli.agent.epId + "/" + path

	dataBytes, err := json.Marshal(params)
	if err != nil {
		log.Println("Marshal error:", err)
		return nil, err
	}

	resBytes, err := cli.restPost(url, dataBytes)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return nil, err
	}
	log.Println("resBytes len:", len(resBytes))

	var inst Instance
	if err := json.Unmarshal(resBytes, &inst); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}

	return &inst, nil
}

func (cli *Cli) restSetParams(path string, params map[string]string) error {
	url := cli.cfg.apiServerAddr + SET_PARAMS + cli.agent.epId + "/" + path

	dataBytes, err := json.Marshal(params)
	if err != nil {
		log.Println("Marshal error:", err)
		return err
	}

	_, err = cli.restPost(url, dataBytes)
	if err != nil {
		log.Println("Error in RESt Read, err", err)
		return err
	}
	return nil
}
func (cli *Cli) restOperateCmd(path string, inputArgs map[string]string) error {
	url := cli.cfg.apiServerAddr + OPERATE_CMD + cli.agent.epId + "/" + path

	dataBytes, err := json.Marshal(inputArgs)
	if err != nil {
		log.Println("Marshal error:", err)
		return err
	}

	_, err = cli.restPost(url, dataBytes)
	if err != nil {
		log.Println("Error in RESt Operate Cmd, err", err)
		return err
	}
	return nil
}

func (cli *Cli) restMtpGetInfo() (*MtpInfo, error) {
	url := cli.cfg.apiServerAddr + GET_MTPINFO

	data, err := cli.restGet(url)
	if err != nil {
		log.Println("Error in RESt Operate Cmd, err", err)
		return nil, err
	}

	mtpInfo := &MtpInfo{}
	if err := json.Unmarshal(data, mtpInfo); err != nil {
		log.Println("Unmarshal error:", err)
		return nil, err
	}
	return mtpInfo, nil
}

func (cli *Cli) restGet(url string) ([]byte, error) {
	log.Println("Sending GET to:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("n4admin", "n4defaultpass")

	resp, err := cli.rest.client.Do(req)
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}

	log.Println("HTTP Status:", resp.Status)
	if resp.StatusCode != 200 {
		errStr := string(bodyBytes)
		log.Println("HTTP Error Msg:", errStr)
		return nil, errors.New(errStr)
	}
	return bodyBytes, nil
}

func (cli *Cli) restPost(url string, data []byte) ([]byte, error) {
	log.Println("Sending POST to:", url)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("n4admin", "n4defaultpass")

	log.Println("HTTP post to url:", url)
	resp, err := cli.rest.client.Do(req)
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("HTTP Error code:", resp.Status)
		return nil, errors.New(resp.Status)
	}
	log.Println("HTTP Status:", resp.Status)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("restErr:", err)
		return nil, err
	}
	return bodyBytes, nil
}
