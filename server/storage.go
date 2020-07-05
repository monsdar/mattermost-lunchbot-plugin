package main

import (
	"bytes"
	"encoding/json"
)

const (
	//KVKEY is the key used for storing the data in the KVStorage
	KVKEY = "LunchbotData"
)

// ReadFromStorage reads LunchbotData from the KVStore. Makes sure that data is inited for the given team and channel
func (p *Plugin) ReadFromStorage() LunchbotData {
	data := LunchbotData{}
	kvData, err := p.API.KVGet(KVKEY)
	if err != nil {
		//do nothing.. we'll return an empty LunchbotData then...
	}
	if kvData != nil {
		json.Unmarshal(kvData, &data)
	}

	return data
}

// WriteToStorage writes the given data to storage
func (p *Plugin) WriteToStorage(data *LunchbotData) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(data)
	p.API.KVSet(KVKEY, reqBodyBytes.Bytes())
}

// ClearStorage removes all stored data from KVStorage
func (p *Plugin) ClearStorage() {
	p.API.KVDelete(KVKEY)
}
