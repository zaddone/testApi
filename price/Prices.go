package price
/**

prices [map[
	type:PRICE 
	time:2017-05-05T20:59:58.502857541Z 
	status:non-tradeable 
	closeoutBid:1.09924
	closeoutAsk:1.10054 
	instrument:EUR_USD
	tradeable:false 

	quoteHomeConversionFactors:map[
		positiveUnits:1.00000000 
		negativeUnits:1.00000000] 
	unitsAvailable:map[
		reduceOnly:map[short:0 long:0] 
		default:map[long:4543846 short:4547855] 
		openOnly:map[long:4543846 short:4547855] 
		reduceFirst:map[long:4543846 short:4547855]] 
	bids:[map[price:1.09940 liquidity:1e+07]]
	asks:[map[price:1.10037 liquidity:1e+07]] 
	]]

**/
type QHCF struct {
	PositiveUnits string `json:"positiveUnits"`
	NegativeUnits string `json:"negativeUnits"`
}
type OO struct {
	Long string `json:"long"`
	Short string `json:"short"`
}
type UA struct {
	OpenOnly *OO `json:"openOnly"`
	ReduceFirst *OO `json:"reduceFirst"`
	ReduceOnly  *OO `json:"reduceOnly"`
	Default  *OO `json:"default"`
}
type Pr struct {
	Price string `json:"price"`
	Liquidity float64 `json:"liquidity"`
}
type Prices struct {
	Type string `json:"type"`
	Time string `json:"time"`
	Instrument string `json:"instrument"`
	CloseoutBid string `json:"closeoutBid"`
	CloseoutAsk string `json:"closeoutAsk"`
	Tradeable bool `json:"tradeable"`
	Status string `json:"status"`

	Bids []*Pr  `json:"bids"`
	Asks []*Pr  `json:"asks"`
	QuoteHomeConversionFactors *QHCF `json:"quoteHomeConversionFactors"`
	UnitsAvailable *UA `json:"unitsAvailable"`
}
