package morpheus

import (
	"bytes"
	"encoding/json"
	"log"
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

func supressInputHereDOC(k, old, new string, d *schema.ResourceData) bool {
	log.Printf("OLD RESPONSE: %v", strings.TrimSuffix(old, "\r\n"))
	log.Printf("NEW RESPONSE: %v", strings.TrimSuffix(new, "\r\n"))
	if strings.TrimSuffix(old, "\r\n") == strings.TrimSuffix(new, "\r\n") {
		log.Println("OLD AND NEW RESPONSES MATCH")
		return true
	} else {
		log.Println("***OLD AND NEW RESPONSES DON'T MATCH***")
		return false
	}
}
