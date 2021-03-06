package models

import (
	"encoding/json"
	"testing"

	"github.com/projectcypress/cdatools/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestImportHQMFTemplateJSON(t *testing.T) {
	var origID = "2.16.840.1.113883.3.560.1.24"
	var def = hds.GetTemplateDefinition(origID)
	assert.Equal(t, def.Definition, "diagnosis")
	assert.Equal(t, def.Status, "resolved")
	var id = hds.GetID(def)
	assert.Equal(t, id, origID)
}

func TestHqmfToQrdaOid(t *testing.T) {
	// ["Device, Applied", "Encounter, Performed", "Diagnostic Study, Intolerance"]
	hqmfOids := []string{"2.16.840.1.113883.3.560.1.10", "2.16.840.1.113883.3.560.1.79", "2.16.840.1.113883.3.560.1.39"}
	vsOids := []string{"", "", "", "2.16.840.1.113883.3.117.1.7.1.403", "2.16.840.1.113883.3.526.3.1279"}
	qrdaOids := []string{"2.16.840.1.113883.10.20.24.3.7", "2.16.840.1.113883.10.20.24.3.23", "2.16.840.1.113883.10.20.24.3.16"}
	for i, hqmfOid := range hqmfOids {
		assert.Equal(t, qrdaOids[i], hds.HqmfToQrdaOid(hqmfOid, vsOids[i]))
	}
}

func TestCodeDisplayForQrdaOid(t *testing.T) {
	// invalid qrda oid
	codeDisplays := hds.CodeDisplayForQrdaOid("not a valid qrda oid", "r3")
	assert.Equal(t, 0, len(codeDisplays))

	// qrda oid with multiple code displays
	codeDisplays = hds.CodeDisplayForQrdaOid("2.16.840.1.113883.10.20.24.3.23", "r3_1")
	assert.Equal(t, 5, len(codeDisplays))
	if len(codeDisplays) > 0 {
		assert.Equal(t, "entryCode", codeDisplays[0].CodeType)
		assert.Equal(t, "code", codeDisplays[0].TagName)
		assert.Equal(t, false, codeDisplays[0].ExcludeNullFlavor)
		assert.Equal(t, []string{"2.16.840.1.113883.3.88.12.80.32"}, codeDisplays[0].PreferredCodeSets)
	}
}

func TestGetAllDataCriteriaForOneMeasure(t *testing.T) {
	mes := make([]Measure, 1)
	measureData := append([]byte("["), append(fixtures.Cms9v4a, []byte("]")...)...)
	json.Unmarshal(measureData, &mes)

	assert.Equal(t, len(allDataCriteria(mes)), 28)
}

func TestGetAllDatacriteriaForMultipleMeasures(t *testing.T) {
	mes := make([]Measure, 2)
	measureData := append([]byte("["), append(append(fixtures.Cms9v4a, append([]byte(","), fixtures.Cms26v3...)...), []byte("]")...)...)
	json.Unmarshal(measureData, &mes)

	assert.Equal(t, len(allDataCriteria(mes)), 48)
}

func TestGetUniqueDataCriteriaForOneMeasure(t *testing.T) {
	mes := make([]Measure, 1)
	measureData := append([]byte("["), append(fixtures.Cms9v4a, []byte("]")...)...)
	json.Unmarshal(measureData, &mes)
	assert.Equal(t, len(UniqueDataCriteria(allDataCriteria(mes))), 15)
}

func TestSetCodeDisplaysForEntry(t *testing.T) {
	entry := &Entry{Oid: "2.16.840.1.113883.3.560.1.79"} // encounter performed hqmf oid
	mapDataCriteria := Mdc{}
	assert.Equal(t, 0, len(entry.CodeDisplays))
	hds.SetCodeDisplaysForEntry(entry, mapDataCriteria, "r3")
	assert.Equal(t, 5, len(entry.CodeDisplays)) // three code displays the encounter performed entry
}

func TestAllPreferredCodeSetsIfNeeded(t *testing.T) {
	// code sets without the "*" string
	preferredCodeSets := []string{"one", "two", "three"}
	codeDisplays := make([]CodeDisplay, 1)
	codeDisplays[0] = CodeDisplay{PreferredCodeSets: preferredCodeSets}
	allPerferredCodeSetsIfNeeded(codeDisplays)
	assert.Equal(t, preferredCodeSets, codeDisplays[0].PreferredCodeSets)

	// all code sets indicated by the "*" string
	preferredCodeSets = append(preferredCodeSets, "*")
	codeDisplays[0] = CodeDisplay{PreferredCodeSets: preferredCodeSets}
	allPerferredCodeSetsIfNeeded(codeDisplays)
	assert.Equal(t, true, unorderedStringSlicesEqual(CodeSystemNames(), codeDisplays[0].PreferredCodeSets), "preferred code sets should include all code system names")
}

func TestStringInSlice(t *testing.T) {
	stringSlice := []string{"one", "two", "four"}
	assert.Equal(t, true, stringInSlice("one", stringSlice))
	assert.Equal(t, true, stringInSlice("two", stringSlice))
	assert.Equal(t, false, stringInSlice("three", stringSlice))
}

// Test helper
func unorderedStringSlicesEqual(a, b []string) bool {
	if a == nil || b == nil || len(a) != len(b) {
		return false
	}
	for _, elemA := range a {
		if !stringInSlice(elemA, b) {
			return false
		}
	}
	return true
}
