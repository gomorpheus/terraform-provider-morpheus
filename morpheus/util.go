package morpheus

import (
	"encoding/json"
	"reflect"
)

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func parseEnvironmentVariables(variables []interface{}) []map[string]interface{} {
	var evars []map[string]interface{}
	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]interface{})
		evarconfig := variables[i].(map[string]interface{})
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "value":
				row["value"] = v.(string)
			case "export":
				row["export"] = v.(bool)
			case "masked":
				row["masked"] = v
			}
		}
		evars = append(evars, row)
		//log.Printf("evars payload: %s", evars)
	}
	return evars
}
