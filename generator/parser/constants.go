package parser

// Every type to be used in the system must
// implement the interface with the NobjectImplementationMethod method
const NobjectImplementationMethod = "GetTypeName"
const CustomIdImplementationMethod = "GetId"
const ConstructorPrefix = "New"
const ReConstructorPrefix = "ReNew"

const OrginalPackageAlias = "org"
const HandlerInputParameterType = "lib.HandlerParameters"
const ReferenceType = "lib.Reference"
const ReferenceListType = "lib.ReferenceList"
const ReadonlyTag = "\"readonly\""
const IndexTag = "\"index\""
const TagKey = "nubes"
const HandlerInputParameterName = "input"
const HandlerParameters = "(" + HandlerInputParameterName + " " + HandlerInputParameterType + ")"
const HandlerInputEmbededOrginalFunctionParameterName = "Parameter"
const LibErrorVariableName = "_libError"
const TemporaryReceiverName = "tempReceiverName"
const LibImportPath = "\"github.com/Astenna/Nubes/lib\""

// Prefixes of repository operations
const (
	GetPrefix    = "Get"
	SetPrefix    = "Set"
	CreatePrefix = "Create"
	DeletePrefix = "Delete"
	UpdatePrefix = "Update"
)
