package metrika_sdk

const visitsFields = "ym:s:visitID,ym:s:counterID,ym:s:watchIDs,ym:s:date,ym:s:dateTime,ym:s:dateTimeUTC,ym:s:isNewUser,ym:s:startURL,ym:s:endURL,ym:s:pageViews,ym:s:visitDuration,ym:s:bounce,ym:s:ipAddress,ym:s:regionCountry,ym:s:regionCity,ym:s:regionCountryID,ym:s:regionCityID,ym:s:clientID,ym:s:counterUserIDHash,ym:s:networkType,ym:s:goalsID,ym:s:goalsSerialNumber,ym:s:goalsDateTime,ym:s:goalsPrice,ym:s:goalsOrder,ym:s:goalsCurrency,ym:s:lastTrafficSource,ym:s:lastAdvEngine,ym:s:lastReferalSource,ym:s:lastSearchEngineRoot,ym:s:lastSearchEngine,ym:s:lastSocialNetwork,ym:s:lastSocialNetworkProfile,ym:s:referer,ym:s:lastDirectClickOrder,ym:s:lastDirectBannerGroup,ym:s:lastDirectClickBanner,ym:s:lastDirectClickOrderName,ym:s:lastClickBannerGroupName,ym:s:lastDirectClickBannerName,ym:s:lastDirectPhraseOrCond,ym:s:lastDirectPlatformType,ym:s:lastDirectPlatform,ym:s:lastDirectConditionType,ym:s:lastCurrencyID,ym:s:from,ym:s:UTMCampaign,ym:s:UTMContent,ym:s:UTMMedium,ym:s:browserLanguage,ym:s:browserCountry,ym:s:clientTimeZone,ym:s:deviceCategory,ym:s:mobilePhone,ym:s:mobilePhoneModel,ym:s:operatingSystemRoot,ym:s:operatingSystem,ym:s:browser,ym:s:browserMajorVersion,ym:s:browserMinorVersion,ym:s:browserEngine,ym:s:browserEngineVersion1,ym:s:browserEngineVersion2,ym:s:browserEngineVersion3,ym:s:browserEngineVersion4,ym:s:cookieEnabled,ym:s:javascriptEnabled,ym:s:screenFormat,ym:s:screenColors,ym:s:screenOrientation,ym:s:screenWidth,ym:s:screenHeight,ym:s:physicalScreenWidth,ym:s:physicalScreenHeight,ym:s:windowClientWidth,ym:s:windowClientHeight,ym:s:purchaseID,ym:s:purchaseDateTime,ym:s:purchaseAffiliation,ym:s:purchaseRevenue,ym:s:purchaseTax,ym:s:purchaseShipping,ym:s:purchaseCoupon,ym:s:purchaseCurrency,ym:s:purchaseProductQuantity,ym:s:productsPurchaseID,ym:s:productsID,ym:s:productsName,ym:s:productsBrand,ym:s:productsCategory,ym:s:productsCategory1,ym:s:productsCategory2,ym:s:productsCategory3,ym:s:productsCategory4,ym:s:productsCategory5,ym:s:productsVariant,ym:s:productsPosition,ym:s:productsPrice,ym:s:productsCurrency,ym:s:productsCoupon,ym:s:productsQuantity,ym:s:impressionsURL,ym:s:impressionsDateTime,ym:s:impressionsProductID,ym:s:impressionsProductName,ym:s:impressionsProductBrand,ym:s:impressionsProductCategory,ym:s:impressionsProductCategory1,ym:s:impressionsProductCategory2,ym:s:impressionsProductCategory3,ym:s:impressionsProductCategory4,ym:s:impressionsProductCategory5,ym:s:impressionsProductVariant,ym:s:impressionsProductPrice,ym:s:impressionsProductCurrency,ym:s:impressionsProductCoupon,ym:s:offlineCallTalkDuration,ym:s:offlineCallHoldDuration,ym:s:offlineCallMissed,ym:s:offlineCallTag,ym:s:offlineCallFirstTimeCaller,ym:s:offlineCallURL,ym:s:parsedParamsKey1,ym:s:parsedParamsKey2,ym:s:parsedParamsKey3,ym:s:parsedParamsKey4,ym:s:parsedParamsKey5,ym:s:parsedParamsKey6,ym:s:parsedParamsKey7,ym:s:parsedParamsKey8,ym:s:parsedParamsKey9,ym:s:parsedParamsKey10"
const hitsFiels = ""

const (
	LogsListUrl    = "https://api-metrika.yandex.net/management/v1/counter/%v/logrequests"
	LogsStatusUrl  = "https://api-metrika.yandex.net/management/v1/counter/%v/logrequest/%v"
	DownloadLogUrl = "https://api-metrika.yandex.net/management/v1/counter/%v/logrequest/%v/part/%v/download"
	CountersUrl    = "https://api-metrika.yandex.net/management/v1/counters"
	CreateLogUrl   = "https://api-metrika.yandex.net/management/v1/counter/%v/logrequests"
	DeleteLogUrl   = "https://api-metrika.yandex.net/management/v1/counter/%v/logrequest/%v/clean"
)

type CounterResponse struct {
	Counters []Counter `json:"counters"`
}

type Counter struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MetrikaResponse struct {
	LogReq LogRequest `json:"log_request"`
}

type LogRequest struct {
	RequestID   int      `json:"request_id"`
	CounterID   int      `json:"counter_id"`
	Source      string   `json:"source"`
	Date1       string   `json:"date1"`
	Date2       string   `json:"date2"`
	Fields      []string `json:"fields"`
	Status      string   `json:"status"`
	Size        int      `json:"size"`
	Parts       []Part   `json:"parts"`
	Attribution string   `json:"attribution"`
}

type Part struct {
	PartNumber int `json:"part_number"`
	Size       int `json:"size"`
}
