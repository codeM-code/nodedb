package nodedb

import "encoding/json"

type GenericDataSlice []map[string]interface{}
type GenericData map[string]interface{}

func (g GenericData) Bytes() ([]byte, error) {
	return json.Marshal(g)
}

func (gs GenericDataSlice) Bytes() ([]byte, error) {
	return json.Marshal(gs)
}

func GenericDataFromBytes(bytes []byte) (GenericData, error) {
	result := GenericData{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func GenericDataSliceFromBytes(bytes []byte) (GenericDataSlice, error) {
	result := GenericDataSlice{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil

}
