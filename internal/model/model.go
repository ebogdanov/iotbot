package model

type BaseResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	T       int64  `json:"t"`
}

type DeviceModel struct {
	UUID   string `json:"uuid"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Sub    bool   `json:"sub"`
	Model  string `json:"model"`
	Status []struct {
		Code  string      `json:"code"`
		Value interface{} `json:"value"`
	} `json:"status,omitempty"`
	Category    string `json:"category"`
	Online      bool   `json:"online"`
	ID          string `json:"id"`
	TimeZone    string `json:"time_zone"`
	LocalKey    string `json:"local_key"`
	UpdateTime  int64  `json:"update_time"`
	ActiveTime  int64  `json:"active_time"`
	OwnerID     string `json:"owner_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
}

type DeviceResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Success bool        `json:"success"`
	Result  DeviceModel `json:"result,omitempty"`
	T       int64       `json:"t"`
}

type DeviceList struct {
	Result struct {
		LastRowKey string   `json:"last_row_key,omitempty"`
		List       []Device `json:"list"`
		Total      int      `json:"total,omitempty"`
		HasMore    bool     `json:"has_more,omitempty"`
	} `json:"result"`
	T       int64 `json:"t"`
	Success bool  `json:"success"`
}

type Device struct {
	Sub          bool   `json:"sub"`
	CategoryName string `json:"category_name"`
	CreateTime   int64  `json:"create_time"`
	LocalKey     string `json:"local_key"`
	OwnerId      string `json:"owner_id"`
	Ip           string `json:"ip"`
	Icon         string `json:"icon"`
	Lon          string `json:"lon"`
	TimeZone     string `json:"time_zone"`
	ProductName  string `json:"product_name"`
	Uuid         string `json:"uuid"`
	GatewayId    string `json:"gateway_id"`
	ActiveTime   int64  `json:"active_time"`
	UpdateTime   int64  `json:"update_time"`
	ProductId    string `json:"product_id"`
	Name         string `json:"name"`
	Online       bool   `json:"online"`
	Model        string `json:"model"`
	Id           string `json:"id"`
	Category     string `json:"category"`
	Lat          string `json:"lat"`
}

type UserHomesResponse struct {
	Code    int    `json:"code,omitempty"`
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
	T       int64  `json:"t"`
	Result  []struct {
		GeoName string `json:"geo_name"`
		HomeId  int    `json:"home_id"`
		Lat     int    `json:"lat"`
		Lon     int    `json:"lon"`
		Name    string `json:"name"`
		Role    string `json:"role"`
	} `json:"result,omitempty"`
}

type ScenariosResponse struct {
	Result []struct {
		Actions []struct {
			ActionExecutor   string `json:"action_executor"`
			EntityId         string `json:"entity_id"`
			ExecutorProperty struct {
				Switch1 bool   `json:"switch_1,omitempty"`
				Hours   string `json:"hours,omitempty"`
				Minutes string `json:"minutes,omitempty"`
				Seconds string `json:"seconds,omitempty"`
			} `json:"executor_property"`
		} `json:"actions"`
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
		SceneId string `json:"scene_id"`
		Status  string `json:"status"`
	} `json:"result,omitempty"`
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
	T       int64  `json:"t"`
	Tid     string `json:"tid"`
}

type DeviceItem struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	ProductId string `json:"product_id"`
	HomeId    string `json:"home_id"`
}

type SceneDeviceResponse struct {
	Result  []DeviceItem `json:"result,omitempty"`
	Success bool         `json:"success"`
	Msg     string       `json:"msg,omitempty"`
	T       int64        `json:"t"`
	Tid     string       `json:"tid"`
}

type UserRequestName struct {
	Success bool
	Msg     string
}
