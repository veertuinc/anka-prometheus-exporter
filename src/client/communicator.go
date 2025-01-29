package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/veertuinc/anka-prometheus-exporter/src/log"
	"github.com/veertuinc/anka-prometheus-exporter/src/state"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

var lock = &sync.Mutex{}
var updateLock = &sync.Mutex{}

type Communicator struct {
	controllerAddress string
	username          string
	password          string
	uak               UAK
	encodedTAPData    string
}

func (comm *Communicator) UpdateEncodedTAPData() error {
	var err error
	if updateLock.TryLock() {
		defer updateLock.Unlock()
		if err = comm.TestConnection(); err.Error() == "Authentication Required" {
			data, err := setUpUAK(comm.uak, comm.controllerAddress)
			if err != nil {
				return err
			}
			comm.encodedTAPData = data
			err = nil
		}
		if err = comm.TestConnection(); err != nil {
			return err
		}
		log.Info("[auth::uak] obtained new UAK session")
	}
	return err
}

func NewCommunicator(addr, username, password string, certs ClientTLSCerts, uak UAK) (*Communicator, error) {
	comm := &Communicator{
		controllerAddress: addr,
		username:          username,
		password:          password,
		uak:               uak,
	}

	if err := setUpTLS(certs); err != nil {
		return nil, err
	}

	if uak.ID != "" {
		log.Info(fmt.Sprintf("[auth::uak] Using User API Key | ID: %s", uak.ID))
		if err := comm.UpdateEncodedTAPData(); err != nil {
			return nil, err
		}
	}

	return comm, nil
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
		log.Debug("status endpoint communication success")
		return nil
	} else {
		return errors.New(resp.Message)
	}
}

func (comm *Communicator) GetStatus() (interface{}, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint := "/api/v1/status"
	resp := &types.StatusResponse{}
	d, err := comm.getData(endpoint, resp)
	if err != nil {
		return nil, fmt.Errorf("getting status error: %s", err)
	}
	return d, nil
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
		return nil, fmt.Errorf("getting registry templates error: %s", err.Error())
	}
	templatesArray := templates.([]types.Template)
	templatesMap := state.GetState().GetTemplatesMap()
	for i, template := range templatesArray {
		if templatesMap[template.UUID].Size != template.Size {
			endpoint := "/api/v1/registry/vm?id=" + template.UUID
			resp := &types.RegistryTemplateTagsResponse{}
			tagsData, err := comm.getData(endpoint, resp)
			if err != nil {
				return nil, fmt.Errorf("getting registry template %s/%s tags error: %s", template.UUID, template.Name, err.Error())
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

func (comm *Communicator) fetchResponseData(endpoint string, repsObject types.Response) (types.Response, error) {
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
	return repsObject, nil
}

func (comm *Communicator) getData(endpoint string, repsObject types.Response) (interface{}, error) {

	repsObject, err := comm.fetchResponseData(endpoint, repsObject)
	if err != nil {
		return nil, err
	}

	retryCount := 2
	for repsObject.GetStatus() != "OK" && retryCount < 4 {
		if comm.uak.ID != "" && repsObject.GetMessage() == "Authentication Required" {
			log.Warn("[auth::uak] uak session expired")
			err = comm.UpdateEncodedTAPData()
			if err != nil {
				log.Error(fmt.Sprintf("could not renew TAP for UAK: %+v", err))
			}
			repsObject, err = comm.fetchResponseData(endpoint, repsObject)
			if err != nil {
				log.Error(fmt.Sprintf("could not get data (after TAP renewal): %+v", err))
			}
		} else {
			return nil, errors.New(repsObject.GetMessage())
		}
		retryCount++
	}
	if err != nil {
		return nil, err
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
	} else if comm.encodedTAPData != "" {
		req.Header.Set("Authorization", "Bearer "+comm.encodedTAPData)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return r, nil
}
