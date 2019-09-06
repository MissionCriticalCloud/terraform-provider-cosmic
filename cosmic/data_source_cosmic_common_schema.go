package cosmic

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

// matchFilters returns true or false depending if the passed filters are matched.
func matchFilters(v interface{}, filters *schema.Set) (bool, error) {
	var j map[string]interface{}
	t, _ := json.Marshal(v)
	json.Unmarshal(t, &j)

	for _, f := range filters.List() {
		m := f.(map[string]interface{})

		value, ok := j[m["name"].(string)].(string)
		if !ok {
			return false, fmt.Errorf("Invalid field name: %s", m["name"])
		}

		re, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		if !re.MatchString(value) {
			return false, nil
		}
	}

	return true, nil
}
