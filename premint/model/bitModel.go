package model

type BrowserFingerPrint struct {
	CoreVersion                string      `json:"coreVersion"`
	Ostype                     string      `json:"ostype"`
	Os                         string      `json:"os"`
	Version                    string      `json:"version"`
	UserAgent                  string      `json:"userAgent"`
	IsIpCreateTimeZone         bool        `json:"isIpCreateTimeZone"`
	TimeZone                   string      `json:"timeZone"`
	TimeZoneOffset             int         `json:"timeZoneOffset"`
	WebRTC                     string      `json:"webRTC"`
	IgnoreHttpsErrors          bool        `json:"ignoreHttpsErrors"`
	Position                   string      `json:"position"`
	IsIpCreatePosition         bool        `json:"isIpCreatePosition"`
	Lat                        string      `json:"lat"`
	Lng                        string      `json:"lng"`
	PrecisionData              string      `json:"precisionData"`
	IsIpCreateLanguage         bool        `json:"isIpCreateLanguage"`
	Languages                  string      `json:"languages"`
	IsIpCreateDisplayLanguage  bool        `json:"isIpCreateDisplayLanguage"`
	DisplayLanguages           string      `json:"displayLanguages"`
	OpenWidth                  int         `json:"openWidth"`
	OpenHeight                 int         `json:"openHeight"`
	ResolutionType             string      `json:"resolutionType"`
	Resolution                 string      `json:"resolution"`
	WindowSizeLimit            bool        `json:"windowSizeLimit"`
	DevicePixelRatio           int         `json:"devicePixelRatio"`
	FontType                   string      `json:"fontType"`
	Font                       string      `json:"font"`
	Canvas                     string      `json:"canvas"`
	CanvasValue                interface{} `json:"canvasValue"`
	WebGL                      string      `json:"webGL"`
	WebGLValue                 interface{} `json:"webGLValue"`
	WebGLMeta                  string      `json:"webGLMeta"`
	WebGLManufacturer          string      `json:"webGLManufacturer"`
	WebGLRender                string      `json:"webGLRender"`
	AudioContext               string      `json:"audioContext"`
	AudioContextValue          interface{} `json:"audioContextValue"`
	SpeechVoices               string      `json:"speechVoices"`
	SpeechVoicesValue          interface{} `json:"speechVoicesValue"`
	HardwareConcurrency        string      `json:"hardwareConcurrency"`
	DeviceMemory               string      `json:"deviceMemory"`
	DoNotTrack                 string      `json:"doNotTrack"`
	ClientRectNoiseEnabled     bool        `json:"clientRectNoiseEnabled"`
	ClientRectNoiseValue       int         `json:"clientRectNoiseValue"`
	PortScanProtect            string      `json:"portScanProtect"`
	PortWhiteList              string      `json:"portWhiteList"`
	DeviceInfoEnabled          bool        `json:"deviceInfoEnabled"`
	ComputerName               string      `json:"computerName"`
	MacAddr                    string      `json:"macAddr"`
	DisableSslCipherSuitesFlag bool        `json:"disableSslCipherSuitesFlag"`
	DisableSslCipherSuites     interface{} `json:"disableSslCipherSuites"`
	EnablePlugins              bool        `json:"enablePlugins"`
	Plugins                    string      `json:"plugins"`
}
type BrowserWindow struct {
	GroupID                     string `json:"groupId"`
	Platform                    string `json:"platform"`
	PlatformIcon                string `json:"platformIcon"`
	URL                         string `json:"url"`
	Name                        string `json:"name"`
	Remark                      string `json:"remark"`
	UserName                    string `json:"userName"`
	Password                    string `json:"password"`
	Cookie                      string `json:"cookie"`
	ProxyMethod                 int    `json:"proxyMethod"`
	ProxyType                   string `json:"proxyType"`
	Host                        string `json:"host"`
	Port                        int    `json:"port"`
	IPCheckService              string `json:"ipCheckService"`
	IsIPv6                      bool   `json:"isIpv6"`
	ProxyUserName               string `json:"proxyUserName"`
	ProxyPassword               string `json:"proxyPassword"`
	RefreshProxyURL             string `json:"refreshProxyUrl"`
	IP                          string `json:"ip"`
	Country                     string `json:"country"`
	Province                    string `json:"province"`
	City                        string `json:"city"`
	Workbench                   string `json:"workbench"`
	AbortImage                  bool   `json:"abortImage"`
	AbortMedia                  bool   `json:"abortMedia"`
	StopWhileNetError           bool   `json:"stopWhileNetError"`
	DynamicIPURL                string `json:"dynamicIpUrl"`
	DynamicIPChannel            string `json:"dynamicIpChannel"`
	IsDynamicIPChangeIP         bool   `json:"isDynamicIpChangeIp"`
	DuplicateCheck              int    `json:"duplicateCheck"`
	IsGlobalProxyInfo           bool   `json:"isGlobalProxyInfo"`
	SyncTabs                    bool   `json:"syncTabs"`
	SyncCookies                 bool   `json:"syncCookies"`
	SyncIndexedDB               bool   `json:"syncIndexedDb"`
	SyncLocalStorage            bool   `json:"syncLocalStorage"`
	SyncBookmarks               bool   `json:"syncBookmarks"`
	SyncAuthorization           bool   `json:"syncAuthorization"`
	CredentialsEnableService    bool   `json:"credentialsEnableService"`
	SyncHistory                 bool   `json:"syncHistory"`
	SyncExtensions              bool   `json:"syncExtensions"`
	SyncUserExtensions          bool   `json:"syncUserExtensions"`
	IsValidUsername             bool   `json:"isValidUsername"`
	AllowedSignin               bool   `json:"allowedSignin"`
	ClearCacheFilesBeforeLaunch bool   `json:"clearCacheFilesBeforeLaunch"`
	ClearCookiesBeforeLaunch    bool   `json:"clearCookiesBeforeLaunch"`
	ClearHistoriesBeforeLaunch  bool   `json:"clearHistoriesBeforeLaunch"`
	RandomFingerprint           bool   `json:"randomFingerprint"`
	DisableGPU                  bool   `json:"disableGpu"`
	EnableBackgroundMode        bool   `json:"enableBackgroundMode"`
	MuteAudio                   bool   `json:"muteAudio"`
	DisableTranslatePopup       bool   `json:"disableTranslatePopup"`
}

// 更新
type BitDetailStruct struct {
	Ids                []string `json:"ids"`
	GroupId            string   `json:"group_id,omitempty"`
	Platform           string   `json:"platform"`
	Host               string   `json:"host"`
	Port               string   `json:"port"`
	BrowserFingerPrint `json:"browserFingerPrint"`
}
type UpdateResponseData struct {
	Success bool `json:"success"`
}

// list  分页查询
type ListRequestData struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"pageSize"`
	GroupId  string `json:"groupId"` //分组id
	//Name     string `json:"name"`    // 窗口名称，模糊匹配
	//Seq      int    `json:"seq"`     // 序号，精确查询
	//MinSeq   int    `json:"minSeq"`  //最小序号，范围查询，不可与seq同时使用
	//MaxSeq   int    `json:"maxSeq"`  //最大序号，范围查询，不可与seq同时使用

}
type ListIDs struct {
	ID   string `json:"id"`
	Seq  int    `json:"seq"`
	Code string `json:"code"`
	BrowserWindow
}

type ListResponseData struct {
	Success bool        `json:"success"`
	Data    ListDataRes `json:"data"`
}
type List struct {
	ID                          string      `json:"id"`
	Seq                         int         `json:"seq"`
	Code                        string      `json:"code"`
	GroupID                     string      `json:"groupId"`
	Platform                    string      `json:"platform"`
	URL                         string      `json:"url"`
	Name                        string      `json:"name"`
	UserName                    string      `json:"userName"`
	Password                    string      `json:"password"`
	Cookie                      string      `json:"cookie"`
	OtherCookie                 string      `json:"otherCookie"`
	IsGlobalProxyInfo           bool        `json:"isGlobalProxyInfo"`
	IsIPv6                      bool        `json:"isIpv6"`
	ProxyMethod                 int         `json:"proxyMethod"`
	ProxyType                   string      `json:"proxyType"`
	AgentID                     string      `json:"agentId"`
	IPCheckService              string      `json:"ipCheckService"`
	Host                        string      `json:"host"`
	Port                        int         `json:"port"`
	ProxyUserName               string      `json:"proxyUserName"`
	ProxyPassword               string      `json:"proxyPassword"`
	LastIP                      string      `json:"lastIp"`
	LastCountry                 string      `json:"lastCountry"`
	IsIPNoChange                bool        `json:"isIpNoChange"`
	IP                          string      `json:"ip"`
	Country                     string      `json:"country"`
	Province                    string      `json:"province"`
	City                        string      `json:"city"`
	DynamicIPURL                string      `json:"dynamicIpUrl"`
	IsDynamicIPChangeIP         bool        `json:"isDynamicIpChangeIp"`
	Remark                      string      `json:"remark"`
	Status                      int         `json:"status"`
	OperUserName                string      `json:"operUserName"`
	OperTime                    string      `json:"operTime"`
	CloseTime                   string      `json:"closeTime"`
	IsDelete                    int         `json:"isDelete"`
	DelReason                   string      `json:"delReason"`
	IsMostCommon                int         `json:"isMostCommon"`
	IsRemove                    int         `json:"isRemove"`
	TempStr                     interface{} `json:"tempStr"`
	CreatedBy                   string      `json:"createdBy"`
	UserID                      string      `json:"userId"`
	CreatedTime                 string      `json:"createdTime"`
	UpdateBy                    string      `json:"updateBy"`
	UpdateTime                  string      `json:"updateTime"`
	RecycleBinRemark            string      `json:"recycleBinRemark"`
	ImportID                    string      `json:"importId"`
	MainUserID                  string      `json:"mainUserId"`
	AbortImage                  bool        `json:"abortImage"`
	AbortMedia                  bool        `json:"abortMedia"`
	StopWhileNetError           bool        `json:"stopWhileNetError"`
	SyncTabs                    bool        `json:"syncTabs"`
	SyncCookies                 bool        `json:"syncCookies"`
	SyncIndexedDB               bool        `json:"syncIndexedDb"`
	SyncBookmarks               bool        `json:"syncBookmarks"`
	SyncAuthorization           bool        `json:"syncAuthorization"`
	SyncHistory                 bool        `json:"syncHistory"`
	SyncGoogleAccount           bool        `json:"syncGoogleAccount"`
	AllowedSignin               bool        `json:"allowedSignin"`
	SyncSessions                bool        `json:"syncSessions"`
	Workbench                   string      `json:"workbench"`
	ClearCacheFilesBeforeLaunch bool        `json:"clearCacheFilesBeforeLaunch"`
	ClearCookiesBeforeLaunch    bool        `json:"clearCookiesBeforeLaunch"`
	ClearHistoriesBeforeLaunch  bool        `json:"clearHistoriesBeforeLaunch"`
	RandomFingerprint           bool        `json:"randomFingerprint"`
	MuteAudio                   bool        `json:"muteAudio"`
	DisableGPU                  bool        `json:"disableGpu"`
	EnableBackgroundMode        bool        `json:"enableBackgroundMode"`
	SyncExtensions              bool        `json:"syncExtensions"`
	SyncUserExtensions          bool        `json:"syncUserExtensions"`
	SyncLocalStorage            bool        `json:"syncLocalStorage"`
	CredentialsEnableService    bool        `json:"credentialsEnableService"`
	DisableTranslatePopup       bool        `json:"disableTranslatePopup"`
	GroupName                   string      `json:"groupName"`
	CreatedName                 string      `json:"createdName"`
	BelongUserName              string      `json:"belongUserName"`
	UpdateName                  string      `json:"updateName"`
	AgentIPCount                interface{} `json:"agentIpCount"`
	BelongToMe                  bool        `json:"belongToMe"`
	SeqExport                   interface{} `json:"seqExport"`
	GroupIDs                    interface{} `json:"groupIDs"`
	BrowserShareID              interface{} `json:"browserShareID"`
	Share                       interface{} `json:"share"`
	IsValidUsername             bool        `json:"isValidUsername"`
	CreateNum                   int         `json:"createNum"`
	IsRandomFinger              bool        `json:"isRandomFinger"`
	RemarkType                  int         `json:"remarkType"`
	RefreshProxyURL             interface{} `json:"refreshProxyUrl"`
	DuplicateCheck              int         `json:"duplicateCheck"`
}
type ListDataRes struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	TotalNum int    `json:"totalNum"`
	List     []List `json:"list"`
}
