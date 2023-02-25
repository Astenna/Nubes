package parser

// Every type to be used in the system must
// implement the interface with the NobjectImplementationMethod method
// METHODS
const NobjectImplementationMethod = "GetTypeName"
const CustomIdImplementationMethod = "GetId"
const ConstructorPrefix = "New"
const DestructorPrefix = "Delete"
const LibraryGetFieldOfType = "GetFieldOfType"
const LibraryGetObjectStateMethod = "GetObjectState"
const InitFunctionName = "Init"
const ReferenceNavigationListCtor = "NewReferenceNavigationList"
const SetField = "SetField"
const Upsert = "Upsert"

// FIELDS & PARAMETER TYPES
const GetStateParamType = "GetStateParam"
const HandlerInputParameterType = "lib.HandlerParameters"
const ReferenceType = "lib.Reference"
const ReferenceListType = "lib.ReferenceList"
const LibraryReferenceNavigationList = "lib.ReferenceNavigationList"
const IsInitializedFieldName = "isInitialized"
const SetFieldParam = "SetFieldParam"

// TAGS
const NubesTagKey = "nubes"
const ReadonlyTag = "readonly"
const IndexTag = "index"
const HasOneTag = "hasOne"
const HasManyTag = "hasMany"
const DynamoDBIgnoreTag = "dynamodbav:\"-\""
const DynamoDBIgnoreValueTag = "-"
const DynamoDBTagKey = "dynamodbav"
const DynamoDBIdTagValue = "Id"
const DynamoDBIdTag = "dynamodbav:\"Id\""
const CustomIdTag = "Id"

// PARAMETER NAMES
const Id = "Id"
const TypeName = "TypeName"
const FieldName = "FieldName"
const HandlerInputParameterName = "input"
const HandlerParameters = "(" + HandlerInputParameterName + " " + HandlerInputParameterType + ")"
const HandlerInputParameterFieldName = "Parameter"
const LibErrorVariableName = "_libError"
const TemporaryReceiverName = "tempReceiverName"

// OTHERS
const LibImportPath = "\"github.com/Astenna/Nubes/lib\""
const OrginalPackageAlias = "org"

// Prefixes of repository operations
const (
	GetPrefix    = "Get"
	SetPrefix    = "Set"
	CreatePrefix = "Create"
	DeletePrefix = "Delete"
	UpdatePrefix = "Update"
)
