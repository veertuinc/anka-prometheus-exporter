package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/veertuinc/anka-prometheus-exporter/src/state"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

var lock = &sync.Mutex{}

type Communicator struct {
	controllerAddress string
	username          string
	password          string
}

func NewCommunicator(addr, username, password string, certs TLSCerts) (*Communicator, error) {

	if err := setUpTLS(certs); err != nil {
		return nil, err
	}

	return &Communicator{
		controllerAddress: addr,
		username:          username,
		password:          password,
	}, nil
}

func (comm *Communicator) TestConnection() error {
	endpoint := "/api/v1/status"
	r, err := comm.getResponse(endpoint, comm.username, comm.password)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	resp := &types.DefaultResponse{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if resp.Status == "OK" {
		return nil
	} else {
		return errors.New(resp.Message)
	}
}

func (comm *Communicator) GetNodesData() (interface{}, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint := "/api/v1/node"
	resp := &types.NodesResponse{}
	d, err := comm.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting node data error: %s", err)
	}
	return d, nil
}

func (comm *Communicator) GetVmsData() (interface{}, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint := "/api/v1/vm"
	resp := &types.InstancesResponse{}
	d, err := comm.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting vms data error: %s", err)
	}
	templatesMap := state.GetState().GetTemplatesMap()
	instances := d.([]types.Instance)
	for i, v := range instances {
		template, ok := templatesMap[v.Vm.TemplateUUID]
		if !ok {
			continue
		}
		instances[i].Vm.TemplateName = template.Name
	}
	return instances, nil
}

func (comm *Communicator) GetRegistryDiskData() (interface{}, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint := "/api/v1/registry/disk"
	resp := &types.RegistryDiskResponse{}
	d, err := comm.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting registry disk data error: %s", err)
	}
	return d, nil
}

func (comm *Communicator) GetRegistryTemplatesData() (interface{}, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint := "/api/v1/registry/vm"
	resp := &types.RegistryTemplateResponse{}
	templates, err := comm.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting registry templates error: %s", err)
	}
	templatesArray := templates.([]types.Template)
	templatesMap := state.GetState().GetTemplatesMap()
	for i, template := range templatesArray {
		if templatesMap[template.UUID].Size != template.Size {
			endpoint := "/api/v1/registry/vm?id=" + template.UUID
			resp := &types.RegistryTemplateTagsResponse{}
			tagsData, err := comm.getData(endpoint, resp)
			if err != nil {
				return nil, fmt.Errorf("getting registry template %s/%s tags error: %s", template.UUID, template.Name, err)
			}
			tags := tagsData.(types.RegistryTemplateTags)
			templatesArray[i].Tags = tags.Versions
		} else {
			templatesArray[i].Tags = templatesMap[template.UUID].Tags
		}
	}
	state.GetState().SetTemplatesMap(templatesArray)
	return templatesArray, nil
}

func (comm *Communicator) getData(endpoint string, repsObject types.Response) (interface{}, error) {
	r, err := comm.getResponse(endpoint, comm.username, comm.password)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &repsObject); err != nil {
		return nil, err
	}
	if repsObject.GetStatus() != "OK" {
		return nil, errors.New(repsObject.GetMessage())
	}
	return repsObject.GetBody(), nil
}

func (comm *Communicator) getResponse(endpoint, username, password string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", comm.controllerAddress, endpoint)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}
	// set auth
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return r, nil
}
