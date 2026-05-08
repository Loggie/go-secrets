package secrets

const (
	SecretService     = "org.freedesktop.secrets"
	SecretServiceRoot = "org.freedesktop.Secret"
	SecretServicePath = "/org/freedesktop/secrets"
)

const (
	SecretServiceInterface              = SecretServiceRoot + ".Service"
	SecretServiceMethodOpenSession      = SecretServiceInterface + ".OpenSession"
	SecretServiceMethodCreateCollection = SecretServiceInterface + ".CreateCollection"
	SecretServiceMethodSearchItems      = SecretServiceInterface + ".SearchItems"
	SecretServiceMethodUnlock           = SecretServiceInterface + ".Unlock"
	SecretServiceMethodLock             = SecretServiceInterface + ".Lock"
	SecretServiceMethodGetSecrets       = SecretServiceInterface + ".GetSecrets"
	SecretServiceMethodReadAlias        = SecretServiceInterface + ".ReadAlias"
	SecretServiceMethodSetAlias         = SecretServiceInterface + ".SetAlias"
	SecretServicePropertyCollections    = SecretServiceInterface + ".Collections"

	SecretServiceSessionAlgorithmPlain  = "plain"
	SecretServiceSessionAlgorithmDhSha1 = "dh-ietf1024-sha1-aes128-cbc-pkcs7" //gosec:disable G101 -- false positive, not a credential
)

const (
	SecretSessionInterface   = SecretServiceRoot + ".Session"
	SecretSessionMethodClose = SecretSessionInterface + ".Close"
)

const (
	SecretCollectionBasePath    = SecretServicePath + "/collection"
	SecretCollectionDefaultPath = SecretCollectionBasePath + "/default"

	SecretCollectionInterface         = SecretServiceRoot + ".Collection"
	SecretCollectionMethodDelete      = SecretCollectionInterface + ".Delete"
	SecretCollectionMethodSearchItems = SecretCollectionInterface + ".SearchItems"
	SecretCollectionMethodCreateItem  = SecretCollectionInterface + ".CreateItem"
	SecretCollectionPropertyItems     = SecretCollectionInterface + ".Items"
	SecretCollectionPropertyLabel     = SecretCollectionInterface + ".Label"
	SecretCollectionPropertyLocked    = SecretCollectionInterface + ".Locked"
	SecretCollectionPropertyCreated   = SecretCollectionInterface + ".Created"
	SecretCollectionPropertyModified  = SecretCollectionInterface + ".Modified"
)

const (
	SecretItemInterface          = SecretServiceRoot + ".Item"
	SecretItemMethodSetSecret    = SecretItemInterface + ".SetSecret"
	SecretItemMethodGetSecret    = SecretItemInterface + ".GetSecret"
	SecretItemMethodDelete       = SecretItemInterface + ".Delete"
	SecretItemPropertyLocked     = SecretItemInterface + ".Locked"
	SecretItemPropertyAttributes = SecretItemInterface + ".Attributes"
	SecretItemPropertyLabel      = SecretItemInterface + ".Label"
	SecretItemPropertyCreated    = SecretItemInterface + ".Created"
	SecretItemPropertyModified   = SecretItemInterface + ".Modified"
)

const (
	SecretPromptInterface       = SecretServiceRoot + ".Prompt"
	SecretPromptMethodPrompt    = SecretPromptInterface + ".Prompt"
	SecretPromptMethodDismiss   = SecretPromptInterface + ".Dismiss"
	SecretPromptSignalCompleted = SecretPromptInterface + ".Completed"
)
