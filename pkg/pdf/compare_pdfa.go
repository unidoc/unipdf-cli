package pdf

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/pdfa"

	"github.com/unidoc/unipdf-cli/internal/verapdf"
)

// ComparePdfaResult is the result of the PDF/A UniPDF Verification.
type ComparePdfaResult struct {
	MismatchedRules map[string]string `json:"mismatchedRules"`
	IsExpectedValid bool              `json:"isExpectedValid"`
	IsUniPdfValid   bool              `json:"isUniPdfValid"`
	Passed          bool              `json:"passed"`
}

// ComparePdfARules is the command that parses VeraPDF XML report, and compares its result with the UniPDF PDF/A validator.
func ComparePdfARules(pdfDoc io.ReadSeeker, veraPdfXMLReport io.Reader, reportType verapdf.ReportType) (*ComparePdfaResult, error) {
	if pdfDoc == nil {
		return nil, errors.New("undefined pdf document")
	}
	if veraPdfXMLReport == nil {
		return nil, errors.New("undefined veraPdf Xml report")
	}

	if !reportType.Valid() {
		return nil, errors.New("undefined report type")
	}

	var (
		report *verapdf.Report
		err    error
	)
	switch reportType {
	case verapdf.GuiXMLReport:
		report, err = verapdf.ParseGuiXML(veraPdfXMLReport)
	case verapdf.OnlineXMLReport:
		report, err = verapdf.ParseOnlineXML(veraPdfXMLReport)
	case verapdf.CliXMLReport:
		report, err = verapdf.ParseCliXML(veraPdfXMLReport)
	default:
		return nil, fmt.Errorf("unknown report type: %v", reportType)
	}
	if err != nil {
		return nil, err
	}

	var profile pdfa.Profile
	switch report.ConformanceLevel {
	case 1:
		switch strings.ToLower(report.ConformanceVariant) {
		case "a":
			profile = pdfa.NewProfile1A(nil)
		case "b":
			profile = pdfa.NewProfile1B(nil)
		default:
			return nil, fmt.Errorf("unknown vera pdf conformance variant: %v", report.ConformanceVariant)
		}
	default:
		return nil, errors.New("currently command support only first PDF/A profile (1A, 1B) ")
	}

	r, err := model.NewCompliancePdfReader(pdfDoc)
	if err != nil {
		return nil, fmt.Errorf("reading pdf document failed: %w", err)
	}

	isEncrypted, err := r.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		return nil, errors.New("pdf document is encrypted - not supported by PDF/A standard")
	}

	err = profile.ValidateStandard(r)
	if err == nil && len(report.ViolatedRules) == 0 {
		return &ComparePdfaResult{Passed: true}, nil
	}
	if err == nil {
		mp := map[string]string{}
		for _, r := range report.ViolatedRules {
			mp[r] = "not found"
		}
		return &ComparePdfaResult{MismatchedRules: mp, IsExpectedValid: false, IsUniPdfValid: true, Passed: false}, nil
	}

	var pdfaErr pdfa.VerificationError
	if !errors.As(err, &pdfaErr) {
		return nil, fmt.Errorf("unexpected validation error type: %T", err)
	}

	expected := map[string]struct{}{}
	for _, rule := range report.ViolatedRules {
		expected[rule] = struct{}{}
	}
	found := map[string]struct{}{}

	for _, rule := range pdfaErr.ViolatedRules {
		found[rule.RuleNo] = struct{}{}
	}

	unmatched := map[string]string{}
	for exp := range expected {
		if _, ok := found[exp]; !ok {
			unmatched[exp] = "not found"
		}
	}
	for fr := range found {
		if _, ok := expected[fr]; !ok {
			unmatched[fr] = "unexpected"
		}
	}

	result := ComparePdfaResult{
		MismatchedRules: unmatched,
		IsExpectedValid: len(report.ViolatedRules) == 0,
		IsUniPdfValid:   false,
	}
	result.Passed = result.IsUniPdfValid == result.IsExpectedValid && len(unmatched) == 0
	return &result, nil
}
