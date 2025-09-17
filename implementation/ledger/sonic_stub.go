// +build !sonic

package sonic

import "encoding/json"

var (
    Marshal = json.Marshal
    Unmarshal = json.Unmarshal
    ConfigDefault = struct{}{}
)

type API struct{}

func (API) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }
func (API) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }
