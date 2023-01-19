package parser

// Every type to be used in the system must
// implement the interface with the NobjectImplementationMethod method
// METHODS
const NobjectImplementationMethod = "GetTypeName"
const CustomIdImplementationMethod = "GetId"
const ConstructorPrefix = "New"
const ReConstructorPrefix = "ReNew"
const DestructorPrefix = "Delete"
const LibraryGetObjectStateMethod = "GetObjectState"

// FIELDS
const OrginalPackageAlias = "org"
const HandlerInputParameterType = "lib.HandlerParameters"
const ReferenceType = "lib.Reference"
const ReferenceListType = "lib.ReferenceList"
const LibraryReferenceNavigationList = "lib.ReferenceNavigationList"

// TAGS
const NubesTagKey = "nubes"
const ReadonlyTag = "readonly"
const IndexTag = "index"
const HasOneTag = "hasOne"
const HasManyTag = "hasMany"

// PARAMETER NAMES
const HandlerInputParameterName = "input"
const HandlerParameters = "(" + HandlerInputParameterName + " " + HandlerInputParameterType + ")"
const HandlerInputEmbededOrginalFunctionParameterName = "Parameter"
const LibErrorVariableName = "_libError"
const TemporaryReceiverName = "tempReceiverName"

// OTHERS
const LibImportPath = "\"github.com/Astenna/Nubes/lib\""

// Prefixes of repository operations
const (
	GetPrefix    = "Get"
	SetPrefix    = "Set"
	CreatePrefix = "Create"
	DeletePrefix = "Delete"
	UpdatePrefix = "Update"
)
