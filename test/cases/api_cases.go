package cases

import "trading/internal/api"

var Kline = []api.KlineParams{
	{Symbol: "BTCUSDT", Interval: "1h", TimeStart: 1728432000000, TimeEnd: 1730231999999},
	{Symbol: "BTCUSDT", Interval: "1h", TimeStart: 1715832000000, TimeEnd: 1717631999999},
	{Symbol: "BTCUSDT", Interval: "1h", TimeStart: 1658232000000, TimeEnd: 1660031999999},
}
