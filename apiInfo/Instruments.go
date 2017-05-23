package apiInfo
import(
	"strconv"
)
type Instrument struct {

	Name string

	DisplayPrecision float64
	MarginRate float64

	MaximumOrderUnits float64
	MaximumPositionSize float64
	MaximumTrailingStopDistance float64

	MinimumTradeSize float64
	MinimumTrailingStopDistance float64

	PipLocation float64
	TradeUnitsPrecision float64
	Type string

}
func (self *Instrument)Init(tmp map[string]interface{}) (err error) {
	self.Name = tmp["name"].(string)
	self.PipLocation = tmp["pipLocation"].(float64)
	self.TradeUnitsPrecision = tmp["tradeUnitsPrecision"].(float64)
	self.Type = tmp["type"].(string)
	self.DisplayPrecision = tmp["displayPrecision"].(float64)

	self.MarginRate,err = strconv.ParseFloat(tmp["marginRate"].(string),64)
	if err != nil {
		return err
	}
	self.MaximumOrderUnits,err = strconv.ParseFloat(tmp["maximumOrderUnits"].(string),64)
	if err != nil {
		return err
	}
	self.MaximumPositionSize,err = strconv.ParseFloat(tmp["maximumPositionSize"].(string),64)
	if err != nil {
		return err
	}
	self.MaximumTrailingStopDistance,err = strconv.ParseFloat(tmp["maximumTrailingStopDistance"].(string),64)
	if err != nil {
		return err
	}
	self.MinimumTradeSize,err = strconv.ParseFloat(tmp["minimumTradeSize"].(string),64)
	if err != nil {
		return err
	}
	self.MinimumTrailingStopDistance,err =  strconv.ParseFloat(tmp["minimumTrailingStopDistance"].(string),64)
	if err != nil {
		return err
	}

	return nil
}
