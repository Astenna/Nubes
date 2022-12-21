package parser

// Every type to be used in the system must
// implement the interface with the GetTypeName method
const GetTypeName = "GetTypeName"

const HandlerSuffix = "Handler"
const OrginalPackageAlias = "org"
const HandlerInputParameterType = "lib.HandlerParameters"
const HandlerInputParameterName = "input"
const HandlerParameters = "(" + HandlerInputParameterName + " " + HandlerInputParameterType + ")"
const HandlerInputEmbededOrginalFunctionParameterName = "Parameter"

// Prefixes of repository operations
const (
	GetPrefix    = "Get"
	CreatePrefix = "Create"
	DeletePrefix = "Delete"
	UpdatePrefix = "Update"
)
