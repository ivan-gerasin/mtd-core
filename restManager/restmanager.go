package mtdrest

import (
	"encoding/json"
	"fmt"
	"github.com/ivan-gerasin/mtdcore/mtdmodels"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RemoteRestServerManager struct {
	address string
}

func (manager *RemoteRestServerManager) UseAddress(addr string) {
	manager.address = addr
}

func (manager *RemoteRestServerManager) List() (error, *mtdmodels.ToDoGlobal) {
	resp, err := http.Get(manager.address + "/list")
	if err != nil {
		return &mtdmodels.MtdError{
			Why:           "Failed on sending request to remote REST mtd-server",
			Where:         "RemoteRestServerManager->List->http.Get",
			OriginalError: &err,
		}, nil
	}
	defer resp.Body.Close()
	buffer := make([]byte, resp.ContentLength)

	_, err = resp.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return &mtdmodels.MtdError{
			Why:           "Failed while reading response from List endpoint from remote mtd REST server",
			Where:         "RemoteRestServerManager->List->http.Get->Body.Read",
			OriginalError: &err,
		}, nil
	}

	results := make(mtdmodels.ToDoGlobal, 10) // TODO: figure out what is best way identify size
	err = json.Unmarshal(buffer, &results)
	if err != nil {
		return &mtdmodels.MtdError{
			Why:           "Failed while unmarshalling json response from server. Is it json at all?",
			Where:         "RemoteRestServerManager->List->http.Get->Body.Read->json.Unmarshal",
			OriginalError: &err,
		}, nil
	}

	return nil, &results
}

func (manager RemoteRestServerManager) AddItem(item string, priority mtdmodels.Priority) error {
	var stringValue = "{\"item\": \"" + item + "\"}"
	resp, err := http.Post(manager.address+"/add", "application/json", strings.NewReader(stringValue))
	if err != nil {
		return &mtdmodels.MtdError{
			Why:           "Failed while trying to make a request to server",
			Where:         "RemoteRestServerManager->AddItem->http.Post",
			OriginalError: &err,
		}
	}

	defer resp.Body.Close()
	buffer := make([]byte, resp.ContentLength)

	_, err = resp.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return &mtdmodels.MtdError{
			Why:           "Failed while reading response body from server",
			Where:         "RemoteRestServerManager->List->http.POST->Body.Read",
			OriginalError: &err,
		}
	}
	return nil
}

func (manager RemoteRestServerManager) Done(id int) error {
	resp, err := http.Post(manager.address+"/done/"+strconv.Itoa(id), "application/json", nil)
	if err != nil {
		return &mtdmodels.MtdError{
			Why:           "Failed while trying to make a request to server",
			Where:         "RemoteRestServerManager->AddItem->http.Post",
			OriginalError: &err,
		}
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	buffer := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return &mtdmodels.MtdError{
			Why:           "Failed while trying to make a request to server",
			Where:         "RemoteRestServerManager->AddItem->http.Post",
			OriginalError: &err,
		}
	}
	var dat map[string]string
	err = json.Unmarshal(buffer, &dat)
	fmt.Println("Status: " + dat["status"])
	return nil
}
