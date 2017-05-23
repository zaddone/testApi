package apiInfo
import(
//	"fmt"
	"time"
	"strconv"
)
type Candles struct {
	Complete bool
	Mid [4]float64
	Time time.Time
	Volume float64
}
func (self * Candles) Init (tmp map[string]interface{}) (err error) {
	ti,err := strconv.ParseFloat(tmp["time"].(string),64)
	self.Time = time.Unix(int64(ti),0)
//	self.Time, err = time.Parse("2006-01-02T15:04:05.000000000", tmp["time"].(string))
	if err != nil {
		return err
	}
//	fmt.Println(tmp)
	Mid := tmp["mid"].(map[string]interface{})
	if Mid != nil {
//		fmt.Println(Min)
		self.Mid[0],err = strconv.ParseFloat(Mid["o"].(string),64)
		if err != nil {
			return err
		}
		self.Mid[1],err = strconv.ParseFloat(Mid["c"].(string),64)
		if err != nil {
			return err
		}
		self.Mid[2],err = strconv.ParseFloat(Mid["h"].(string),64)
		if err != nil {
			return err
		}
		self.Mid[3],err = strconv.ParseFloat(Mid["l"].(string),64)
		if err != nil {
			return err
		}
	}
	self.Complete = tmp["complete"].(bool)
	self.Volume = tmp["volume"].(float64)
	return nil
}
