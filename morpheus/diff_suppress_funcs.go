package morpheus

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func supressOptionListScripts(k, old, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == strings.TrimSpace(new) {
		return true
	} else {
		return false
	}
}
