package parser

// Every type to be used in the system must
// implement the interface with the NobjectImplementationMethod method
const NobjectImplementationMethod = "GetTypeName"
const CustomIdImplementationMethod = "GetId"

const HandlerSuffix = "Handler"
const OrginalPackageAlias = "org"
const HandlerInputParameterType = "lib.HandlerParameters"
const ReferenceType = "lib.Reference"
const ReadonlyTag = "nubes:\"readonly\""
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
