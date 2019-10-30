package checks

import (
	"strings"

	"github.com/liamg/tfsec/internal/app/tfsec/models"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// TODO move these to util package

func getAttribute(block *hcl.Block, ctx *hcl.EvalContext, name string) (cty.Value, *models.Range, bool) {
	attributes, diagnostics := block.Body.JustAttributes()
	if diagnostics != nil && diagnostics.HasErrors() {
		return cty.NilVal, nil, false
	}

	for _, attribute := range attributes {
		if attribute.Name == name {
			val, diagnostics := attribute.Expr.Value(ctx)
			if diagnostics != nil && diagnostics.HasErrors() {
				return cty.NilVal, nil, false
			}
			return val, convertRange(attribute.Range), true
		}
	}

	return cty.NilVal, nil, false
}

func getBlockName(block *hcl.Block) string {
	var prefix string
	if block.Type != "resource" {
		prefix = block.Type
	}
	return prefix + strings.Join(block.Labels, ".")
}

func convertRange(r hcl.Range) *models.Range {
	return &models.Range{
		Filename:  r.Filename,
		StartLine: r.Start.Line,
		EndLine:   r.End.Line,
	}
}