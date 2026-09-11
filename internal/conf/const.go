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

	Logo      = "logo" // multi-lines text, L1: light, EOL: dark
	Favicon   = "favicon"
	MainColor = "main_color"

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
