package typespec

import "github.com/Astenna/Nubes/generator/parser"

type ExportTemplateInput struct {
	OrginalPackage            string
	OrginalPackageAlias       string
	IsNobjectInOrginalPackage map[string]bool
	TypesWithCustomExport     map[string]parser.CustomExportDefinition
}

type DeleteTemplateInput struct {
	OrginalPackage        string
	OrginalPackageAlias   string
	TypesWithCustomDelete map[string]parser.CustomDeleteDefinition
}
