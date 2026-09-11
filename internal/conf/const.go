package conf

const (
	TypeString = "string"
	TypeSelect = "select"
	TypeBool   = "bool"
	TypeText   = "text"
	TypeNumber = "number"
)

const (
	// site
	VERSION      = "version"
	SiteTitle    = "site_title"
	Announcement = "announcement"
	AllowIndexed = "allow_indexed"
	AllowMounted = "allow_mounted"
	RobotsTxt    = "robots_txt"

	Logo                           = "logo" // multi-lines text, L1: light, EOL: dark
	Favicon                        = "favicon"
	MainColor                      = "main_color"
	HideStorageDetails             = "hide_storage_details"
	HideStorageDetailsInManagePage = "hide_storage_details_in_manage_page"

	// preview
	TextTypes                = "text_types"
	AudioTypes               = "audio_types"
	VideoTypes               = "video_types"
	ImageTypes               = "image_types"
	ProxyTypes               = "proxy_types"
	ProxyIgnoreHeaders       = "proxy_ignore_headers"
	AudioAutoplay            = "audio_autoplay"
	VideoAutoplay            = "video_autoplay"
	PreviewDownloadByDefault = "preview_download_by_default"
	ReadMeAutoRender         = "readme_autorender"
	FilterReadMeScripts      = "filter_readme_scripts"
	NonEFSZipEncoding        = "non_efs_zip_encoding"

	// global
	CustomizeHead           = "customize_head"
	CustomizeBody           = "customize_body"
	LinkExpiration          = "link_expiration"
	SignAll                 = "sign_all"
	PrivacyRegs             = "privacy_regs"
	OcrApi                  = "ocr_api"
	FilenameCharMapping     = "filename_char_mapping"
	ForwardDirectLinkParams = "forward_direct_link_params"
	IgnoreDirectLinkParams  = "ignore_direct_link_params"
	WebauthnLoginEnabled    = "webauthn_login_enabled"
	HandleHookAfterWriting  = "handle_hook_after_writing"
	HandleHookRateLimit     = "handle_hook_rate_limit"
	IgnoreSystemFiles       = "ignore_system_files"

	// index
	SearchIndex     = "search_index"
	AutoUpdateIndex = "auto_update_index"
	IgnorePaths     = "ignore_paths"
	MaxIndexDepth   = "max_index_depth"

	// single
	Token         = "token"
	IndexProgress = "index_progress"

	// SSO
	SSOClientId          = "sso_client_id"
	SSOClientSecret      = "sso_client_secret"
	SSOLoginEnabled      = "sso_login_enabled"
	SSOLoginPlatform     = "sso_login_platform"
	SSOOIDCUsernameKey   = "sso_oidc_username_key"
	SSOOrganizationName  = "sso_organization_name"
	SSOApplicationName   = "sso_application_name"
	SSOEndpointName      = "sso_endpoint_name"
	SSOJwtPublicKey      = "sso_jwt_public_key"
	SSOExtraScopes       = "sso_extra_scopes"
	SSOAutoRegister      = "sso_auto_register"
	SSODefaultDir        = "sso_default_dir"
	SSODefaultPermission = "sso_default_permission"
	SSOCompatibilityMode = "sso_compatibility_mode"

	// ldap
	LdapLoginEnabled      = "ldap_login_enabled"
	LdapServer            = "ldap_server"
	LdapSkipTlsVerify     = "ldap_skip_tls_verify"
	LdapManagerDN         = "ldap_manager_dn"
	LdapManagerPassword   = "ldap_manager_password"
	LdapUserSearchBase    = "ldap_user_search_base"
	LdapUserSearchFilter  = "ldap_user_search_filter"
	LdapDefaultPermission = "ldap_default_permission"
	LdapDefaultDir        = "ldap_default_dir"
	LdapLoginTips         = "ldap_login_tips"
)

const (
	UNKNOWN = iota
	FOLDER
	// OFFICE
	VIDEO
	AUDIO
	TEXT
	IMAGE
)

// ContextKey is the type of context keys.
type ContextKey int8

const (
	_ ContextKey = iota

	_ // 1: reserved (was NoTaskKey)
	ApiUrlKey
	UserKey
	ClientIPKey
	ProxyHeaderKey
	RequestHeaderKey
	UserAgentKey
	PathKey
	_ // 11: reserved (was SharingIDKey), kept to preserve the context key numbering
	SkipHookKey
)
