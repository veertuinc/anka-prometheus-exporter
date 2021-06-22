package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/veertuinc/anka-prometheus-exporter/src/state"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type Communicator struct {
	controllerAddress string
}

func NewCommunicator(addr string, certs TLSCerts) (*Communicator, error) {

	if err := setUpTLS(certs); err != nil {
		return nil, err
	}

	return &Communicator{
		controllerAddress: addr,
	}, nil
}

func (this *Communicator) TestConnection() error {
	endpoint := "/api/v1/status"
	r, err := this.getResponse(endpoint)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	resp := &types.DefaultResponse{}
	body, err := ioutil.ReadAll(r.Body)
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

func (this *Communicator) GetNodesData() (interface{}, error) {
	endpoint := "/api/v1/node"
	resp := &types.NodesResponse{}
	d, err := this.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting node data error: %s", err)
	}
	return d, nil
}

func (this *Communicator) GetVmsData() (interface{}, error) {
	endpoint := "/api/v1/vm"
	resp := &types.InstancesResponse{}
	d, err := this.getData(endpoint, resp)
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
		instances[i].Vm.TemplateNAME = template.Name
	}
	return instances, nil
}

func (this *Communicator) GetRegistryDiskData() (interface{}, error) {
	endpoint := "/api/v1/registry/disk"
	resp := &types.RegistryDiskResponse{}
	d, err := this.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting registry disk data error: %s", err)
	}
	return d, nil
}

func (this *Communicator) GetRegistryTemplatesData() (interface{}, error) {
	endpoint := "/api/v1/registry/vm"
	resp := &types.RegistryTemplateResponse{}
	templates, err := this.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting registry templates error: %s", err)
	}
	templatesArray := templates.([]types.Template)
	templatesMap := state.GetState().GetTemplatesMap()
	for i, template := range templatesArray {
		if templatesMap[template.UUID].Size != template.Size {
			endpoint := "/api/v1/registry/vm?id=" + template.UUID
			resp := &types.RegistryTemplateTagsResponse{}
			tagsData, err := this.getData(endpoint, resp)
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

func (this *Communicator) getData(endpoint string, repsObject types.Response) (interface{}, error) {
	r, err := this.getResponse(endpoint)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
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

func (this *Communicator) getResponse(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", this.controllerAddress, endpoint)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return r, nil
}
