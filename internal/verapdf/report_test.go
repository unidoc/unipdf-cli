package verapdf

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportDecoding(t *testing.T) {
	var r reportGuiXML

	data, err := exampleXmlBytes()
	require.NoError(t, err)

	err = xml.NewDecoder(bytes.NewReader(data)).Decode(&r)
	require.NoError(t, err)

	require.NotNil(t, r.Jobs.Job.ValidationReport)

	rules := r.Jobs.Job.ValidationReport.Details.Rule
	if assert.Len(t, rules, 7) {
		assert.Equal(t, rules[0].Clause, "6.2.3")
		assert.Equal(t, rules[0].TestNumber, "4")
		assert.Equal(t, rules[1].Clause, "6.3.6")
		assert.Equal(t, rules[1].TestNumber, "1")
		assert.Equal(t, rules[2].Clause, "6.1.7")
		assert.Equal(t, rules[2].TestNumber, "2")
		assert.Equal(t, rules[3].Clause, "6.2.3")
		assert.Equal(t, rules[3].TestNumber, "2")
		assert.Equal(t, rules[4].Clause, "6.7.3")
		assert.Equal(t, rules[4].TestNumber, "1")
		assert.Equal(t, rules[5].Clause, "6.7.2")
		assert.Equal(t, rules[5].TestNumber, "1")
		assert.Equal(t, rules[6].Clause, "6.4")
		assert.Equal(t, rules[6].TestNumber, "2")
	}
}

func TestParseXML(t *testing.T) {
	data, err := exampleXmlBytes()
	require.NoError(t, err)

	r, err := ParseGuiXML(bytes.NewReader(data))
	require.NoError(t, err)

	assert.Equal(t, 1, r.ConformanceLevel)
	assert.Equal(t, "B", r.ConformanceVariant)
	if assert.Len(t, r.ViolatedRules, 7) {
		assert.Equal(t, "6.2.3-4", r.ViolatedRules[0])
		assert.Equal(t, "6.3.6-1", r.ViolatedRules[1])
		assert.Equal(t, "6.1.7-2", r.ViolatedRules[2])
		assert.Equal(t, "6.2.3-2", r.ViolatedRules[3])
		assert.Equal(t, "6.7.3-1", r.ViolatedRules[4])
		assert.Equal(t, "6.7.2-1", r.ViolatedRules[5])
		assert.Equal(t, "6.4-2", r.ViolatedRules[6])
	}
}

func TestParseCliXML(t *testing.T) {
	data, err := blanco_cliXmlBytes()
	require.NoError(t, err)

	r, err := ParseCliXML(bytes.NewReader(data))
	require.NoError(t, err)

	assert.Equal(t, 1, r.ConformanceLevel)
	assert.Equal(t, "B", r.ConformanceVariant)

	if assert.Len(t, r.ViolatedRules, 7) {
		assert.Equal(t, "6.1.3-1", r.ViolatedRules[0])
		assert.Equal(t, "6.7.3-1", r.ViolatedRules[1])
		assert.Equal(t, "6.7.2-1", r.ViolatedRules[2])
		assert.Equal(t, "6.2.3-2", r.ViolatedRules[3])
		assert.Equal(t, "6.3.4-1", r.ViolatedRules[4])
		assert.Equal(t, "6.4-6", r.ViolatedRules[5])
		assert.Equal(t, "6.4-5", r.ViolatedRules[6])
	}
}

func TestParseOnlineXML(t *testing.T) {
	data, err := blanco_onlineXmlBytes()
	require.NoError(t, err)

	r, err := ParseOnlineXML(bytes.NewReader(data))
	require.NoError(t, err)

	assert.Equal(t, 1, r.ConformanceLevel)
	assert.Equal(t, "B", r.ConformanceVariant)

	if assert.Len(t, r.ViolatedRules, 7) {
		assert.Equal(t, "6.1.3-1", r.ViolatedRules[0])
		assert.Equal(t, "6.7.3-1", r.ViolatedRules[1])
		assert.Equal(t, "6.7.2-1", r.ViolatedRules[2])
		assert.Equal(t, "6.2.3-2", r.ViolatedRules[3])
		assert.Equal(t, "6.3.4-1", r.ViolatedRules[4])
		assert.Equal(t, "6.4-5", r.ViolatedRules[5])
		assert.Equal(t, "6.4-6", r.ViolatedRules[6])
	}
}

func TestReduceCliXML(t *testing.T) {
	data, err := blanco_cliXmlBytes()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)

	err = ReduceCliXML(bytes.NewReader(data), buf)
	require.NoError(t, err)

	r, err := ParseOnlineXML(buf)
	require.NoError(t, err)

	assert.Equal(t, 1, r.ConformanceLevel)
	assert.Equal(t, "B", r.ConformanceVariant)

	if assert.Len(t, r.ViolatedRules, 7) {
		assert.Equal(t, "6.1.3-1", r.ViolatedRules[0])
		assert.Equal(t, "6.7.3-1", r.ViolatedRules[1])
		assert.Equal(t, "6.7.2-1", r.ViolatedRules[2])
		assert.Equal(t, "6.2.3-2", r.ViolatedRules[3])
		assert.Equal(t, "6.3.4-1", r.ViolatedRules[4])
		assert.Equal(t, "6.4-6", r.ViolatedRules[5])
		assert.Equal(t, "6.4-5", r.ViolatedRules[6])
	}
}
