package verapdf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/unidoc/unipdf/v3/common"
)

var _ pflag.Value = (*ReportType)(nil)

// ReportType is the type of the VeraPDF XML report.
type ReportType int

// String implements pflag.Value interface.
func (r ReportType) String() string {
	switch r {
	case GuiXMLReport:
		return "gui"
	case OnlineXMLReport:
		return "online"
	case CliXMLReport:
		return "cli"
	default:
		return "unknown"
	}
}

// Set implements pflag.Value interface.
func (r *ReportType) Set(s string) error {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "gui":
		*r = GuiXMLReport
	case "online":
		*r = OnlineXMLReport
	case "cli":
		*r = CliXMLReport
	default:
		return fmt.Errorf("unknown report type: %s. Allowed values: 'gui','online','cli'", s)
	}
	return nil
}

// Type implements pflag.Value interface.
func (r *ReportType) Type() string {
	return "report-type"
}

const (
	_ ReportType = iota
	GuiXMLReport
	OnlineXMLReport
	CliXMLReport
)

// Valid checks if the report type is valid.
func (r ReportType) Valid() bool {
	return r >= GuiXMLReport && r <= CliXMLReport
}

// Report is the parsed veraPDF XML reportGuiXML.
type Report struct {
	// ViolatedRules are the rules that were violated during error verification.
	ViolatedRules []string
	// ConformanceLevel defines the standard on verification failed.
	ConformanceLevel int
	// ConformanceVariant is the standard variant used on verification.
	ConformanceVariant string
}

// ParseGuiXML parses xml document and gets it's violated rules.
func ParseGuiXML(r io.Reader) (*Report, error) {
	var rp reportGuiXML

	if err := xml.NewDecoder(r).Decode(&rp); err != nil {
		return nil, err
	}

	if rp.Jobs.Job.ValidationReport == nil {
		return nil, errors.New("invalid input. nil validation report")
	}
	return rp.Jobs.Job.ValidationReport.report()
}

// ParseOnlineXML parses XML document created by the online verapdf tool.
func ParseOnlineXML(r io.Reader) (*Report, error) {
	var vr validationResultsXML
	if err := xml.NewDecoder(r).Decode(&vr); err != nil {
		return nil, err
	}

	return vr.report()
}

// ParseReducedXML parses a reduced xml document.
func ParseReducedXML(r io.Reader) (*Report, error) {
	var vr validationResultsXML
	if err := xml.NewDecoder(r).Decode(&vr); err != nil {
		return nil, err
	}

	return vr.report()
}

// ParseCliXML parses XML document created by the verapdf Cli.
func ParseCliXML(r io.Reader) (*Report, error) {
	var rr rawResultsXML
	if err := xml.NewDecoder(r).Decode(&rr); err != nil {
		return nil, err
	}
	if rr.ValidationResults == nil {
		return nil, errors.New("invalid input document")
	}
	return rr.ValidationResults.report()
}

// ReduceCliXML is a function that reduces the size of the input CLI XML document.
// The result is in Online - reduced format.
func ReduceCliXML(r io.Reader, w io.Writer) error {
	var rr rawResultsXML
	e := xml.NewEncoder(w)
	if err := xml.NewDecoder(r).Decode(&rr); err != nil {
		return err
	}
	if err := e.Encode(rr.ValidationResults); err != nil {
		return err
	}
	return nil
}

type rawResultsXML struct {
	ValidationResults *validationResultsXML `xml:"validationResult"`
}

type validationResultsXML struct {
	XMLName xml.Name             `xml:"validationResult"`
	Flavour string               `xml:"flavour,attr"`
	Rules   []assertionRuleIdXML `xml:"assertions>assertion>ruleId"`
}

func (v *validationResultsXML) report() (r *Report, err error) {
	r = &Report{}

	r.ConformanceLevel, r.ConformanceVariant, err = v.parseFlavour()
	if err != nil {
		return nil, err
	}

	m := map[string]struct{}{}

	var sb strings.Builder
	for _, rule := range v.Rules {
		sb.WriteString(rule.Clause)
		sb.WriteRune('-')
		sb.WriteString(rule.TestNumber)

		ruleNo := sb.String()
		if _, ok := m[ruleNo]; ok {
			sb.Reset()
			continue
		}
		m[ruleNo] = struct{}{}
		r.ViolatedRules = append(r.ViolatedRules, ruleNo)
		sb.Reset()
	}
	return r, nil
}

func (v *validationResultsXML) parseFlavour() (level int, conformance string, err error) {
	fv := strings.Split(v.Flavour, "_")
	if len(fv) != 3 {
		common.Log.Debug("invalid validation results flavour: %s", v.Flavour)
		return 0, "", errors.New("invalid validation results flavour")
	}
	level, err = strconv.Atoi(fv[1])
	if err != nil {
		return 0, "", fmt.Errorf("invalid validation results flavor level: %w", err)
	}
	return level, fv[2], nil
}

type assertionRuleIdXML struct {
	Clause     string `xml:"clause,attr"`
	TestNumber string `xml:"testNumber,attr"`
}

type reportGuiXML struct {
	Jobs jobsGuiXML `xml:"jobs"`
}

type jobsGuiXML struct {
	Job jobGuiXML `xml:"job"`
}

type jobGuiXML struct {
	ValidationReport *validationReportGuiXML `xml:"validationReport"`
}

type validationReportGuiXML struct {
	ProfileName string        `xml:"profileName,attr"`
	IsCompliant bool          `xml:"isCompliant,attr"`
	Details     detailsGuiXML `xml:"details"`
}

func (vr *validationReportGuiXML) report() (*Report, error) {
	i := strings.IndexRune(vr.ProfileName, ' ')
	if i == -1 {
		common.Log.Debug("invalid reportGuiXML profile name: %s", vr.ProfileName)
		return nil, errors.New("invalid input veraPDF XML reportGuiXML")
	}
	pn := strings.TrimPrefix(vr.ProfileName[:i], "PDF/A-")
	if len(pn) != 2 {
		common.Log.Debug("invalid reportGuiXML profile name: %s", vr.ProfileName)
		return nil, errors.New("invalid input veraPDF XML reportGuiXML")
	}
	rpt := &Report{}
	rpt.ConformanceLevel = int(pn[0] - '0')
	rpt.ConformanceVariant = pn[1:]

	var sb strings.Builder
	for _, r := range vr.Details.Rule {
		sb.WriteString(r.Clause)
		sb.WriteRune('-')
		sb.WriteString(r.TestNumber)
		rpt.ViolatedRules = append(rpt.ViolatedRules, sb.String())
		sb.Reset()
	}
	return rpt, nil
}

type detailsGuiXML struct {
	Rule []detailsRuleGuiXML `xml:"rule"`
}

type detailsRuleGuiXML struct {
	Clause     string `xml:"clause,attr"`
	TestNumber string `xml:"testNumber,attr"`
}
