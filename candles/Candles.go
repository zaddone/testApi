package candles
import(
	"fmt"
	"time"
	"strconv"
)
type minTmp struct {
	C string `json:"c"`
	H string `json:"h"`
	L string `json:"l"`
	O string `json:"o"`
}
type CandlesTmp struct {
	Complete bool `json:"complete"`
	Min *minTmp `json:"min"`
	Time string `json:"time"`
	Volume int64 `json:"volume"`
}
type Candles struct {
	Complete bool
	Min [4]float64
	Time time.Time
	Volume int64
}
func (self * Candles) Init (tmp *CandlesTmp) error {
	t, err := time.Parse("2006-01-02T15:04:05", tmp.Time)
	if err != nil {
		return err
	}
	self.Min[0],err = strconv.ParseFloat(tmp.Min.C,64)
	if err != nil {
		return err
	}
	self.Min[0],err = strconv.ParseFloat(tmp.Min.H,64)
	if err != nil {
		return err
	}
	self.Min[0],err = strconv.ParseFloat(tmp.Min.L,64)
	if err != nil {
		return err
	}
	self.Min[0],err = strconv.ParseFloat(tmp.Min.O,64)
	if err != nil {
		return err
	}
	self.Complete = tmp.Complete
	self.Volume = tmp.Volume
}
func GetCandlesArray(){

}
