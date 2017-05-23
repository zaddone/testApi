package apiInfo
import(
	"fmt"
	"os"
	"math"
	"net/http"
	"io/ioutil"
	"net"
	"golang.org/x/net/proxy"
	"encoding/json"
	"time"
	"net/url"
	"../fitting"
	"../tree"
)
const (
	MaxLevel int = 2
	GranLen int = 120
)
type Grain struct{
	List []uint64
	IsEnd int
	Over  float64
	Tmp float64
	Complete bool
	Max  int64
	T  int
}
type Account struct {
	Id string `json:"id"`
	Tags []string `json:"tags"`
	Instruments map[string]*Instrument `json:"instruments"`
}
type ApiInfo struct {

	Client *http.Client
	Header http.Header
	Host string
	TmpFileName string
//	Proxy string
	Accounts []*Account
	Gran []string
	WaitSample []*Grain
	TreeL  *tree.TreeList
	TCount [3]int

//	tmpCandles  map[string]uint64

}
func (self *ApiInfo) GetWaitCount(v float64) {

	var newWait []*Grain
	for _,g := range self.WaitSample {
		dif := v - g.Over
		if math.Abs(dif) > g.Tmp {
			if dif >0 {
				g.IsEnd = 1
			}else{
				g.IsEnd = 0
			}
			if g.T == 0 {
				self.TCount[2]++
			}else if (g.T==1) == (g.IsEnd==1 ){
				self.TCount[1]++
			}else{
				self.TCount[0]++
			}
			fmt.Printf("%d\r\n",self.TCount)
			self.TreeL.Append(g.List,g.IsEnd)
		}else{
			newWait = append(newWait,g)
		}
	}
	self.WaitSample = newWait

}
func (self *ApiInfo)GranInit () {
	self.Gran = make([]string,10)
	self.Gran[0]="S10"
	self.Gran[1]="S15"
	self.Gran[2]="S30"
	self.Gran[3]="M1"
	self.Gran[4]="M2"
	self.Gran[5]="M4"
	self.Gran[6]="M5"
	self.Gran[7]="M10"
	self.Gran[8]="M15"
	self.Gran[9]="M30"
}
func (self *ApiInfo) GetInstruments(n int) (ins map[string]*Instrument,err error) {

	if len(self.Accounts[n].Instruments) != 0 {
		return self.Accounts[n].Instruments,nil
	}

	path := self.GetAccountPath(n)
	path += "/instruments"
	da := make(map[string]interface{})
	err = self.ClientDO(path,&da)
	if err != nil {
		return nil,err
	}
	in := da["instruments"].([]interface{})
	ins = make(map[string]*Instrument)
	for _,n := range in {
//		fmt.Println(n.(InstrumentTmp).Name)
		in :=new(Instrument)
		in.Init(n.(map[string]interface{}))

		ins[in.Name] =in
	}
	self.Accounts[n].Instruments = ins
	err =  self.SaveAccounts()
	return ins,err

}
func (self *ApiInfo) GetAccountPath(n int) string {

	return fmt.Sprintf("%s/accounts/%s",self.Host,self.Accounts[n].Id)

}
func (self *ApiInfo)Init(tmpFileName string,Proxy string,Host string,Authorization string) error {
	self.TmpFileName = tmpFileName
	self.GranInit()
	self.TreeL = new(tree.TreeList)
	self.TreeL.Init(MaxLevel)
//	self.tmpCandles = make([string]uint64)

	self.Host = Host
	self.Header = make(http.Header)
//	self.Header.Add("Authorization","Bearer d481a687988ea82d0a09c0527221099d-3e9e6c2697e93f64fd58ecf34ad069db")
	self.Header.Add("Authorization",fmt.Sprintf("Bearer %s",Authorization))
	self.Header.Add("Connection","Keep-Alive")
	self.Header.Add("Accept-Datetime-Format", "UNIX")

	if Proxy=="" {
		self.Client = new(http.Client)
	}else{
		dialer, err := proxy.SOCKS5("tcp",Proxy,
		    nil,
		    &net.Dialer {
		        Timeout: 30 * time.Second,
		        KeepAlive: 30 * time.Second,
		    },
		)
		if err != nil {
			return err
		}
		transport := &http.Transport{
		    Proxy: nil,
		    Dial: dialer.Dial,
		    TLSHandshakeTimeout: 10 * time.Second,
		}
		self.Client = &http.Client{Transport:transport}
	}
	err := self.ReadAccounts()
	if err == nil {
		return nil
	}
	da := make(map[string][]*Account)
	err = self.ClientDO(Host+"/accounts",&da)
	if err != nil {
		return err
	}
	self.Accounts = da["accounts"]
//	L := len(self.Accounts)
	if len(self.Accounts) == 0 {
		fmt.Errorf("accounts == nil")
	}
	return self.SaveAccounts()
}
func (self *ApiInfo) SaveAccounts() error {
	f,err :=os.OpenFile(self.TmpFileName,os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_SYNC,0777)
	if err != nil {
		return err
	}
	defer f.Close()
	d,err := json.Marshal(self.Accounts)
	if err != nil {
		return err
	}
	_,err=f.Write(d)
	if err != nil {
		return err
	}
	return nil
}
func (self *ApiInfo) ReadAccounts() error {
	fi,err := os.Stat(self.TmpFileName)
	if err != nil {
		return err
	}
	data := make([]byte,fi.Size())
	f,err := os.Open(self.TmpFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	n,err := f.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("%d %d",n,len(data))
	}
	return json.Unmarshal(data,&(self.Accounts))
}
func (self *ApiInfo) ClientDO (path string,da interface{} ) error {

	Req,err := http.NewRequest("GET",path,nil)
	if err != nil {
		return err
	}
	Req.Header = self.Header
	res,err := self.Client.Do(Req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
//	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		b,err:=ioutil.ReadAll(res.Body)
		fmt.Println(string(b),err)
		return fmt.Errorf("%d",res.StatusCode)
	}
	b,err:=ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b,da)
}

func (self *ApiInfo)GetCandlesArray(Ins_name,granularity string,from ,count int64) (can []*Candles,err error) {

	path := fmt.Sprintf("%s/instruments/%s/candles?",self.Host,Ins_name)
	uv := url.Values{}
	uv.Add("granularity",granularity)
	uv.Add("price","M")
	uv.Add("count",fmt.Sprintf("%d",count))
	uv.Add("from",fmt.Sprintf("%d",from))
//	uv.Add("to",fmt.Sprintf("%d",to))
	path += uv.Encode()
//	fmt.Println(path)
	da := make(map[string]interface{})
	err = self.ClientDO(path,&da)
	if err != nil {
		return nil,err
	}
	ca := da["candles"].([]interface{})
	lc := len(ca)
	if lc == 0 {
		return nil,fmt.Errorf("candles len = 0")
	}
	can = make([]*Candles,lc)
	for i,c := range ca {
		can[i] = new(Candles)
		can[i].Init(c.(map[string]interface{}))
	}
	return can,nil

}
func (self *ApiInfo)GetCandlesList(Ins *Instrument,to  int64) (key []uint64,err error) {
	key = make([]uint64,len(self.Gran))
	for i,g := range self.Gran {
		key[i],err = self.GetCandlesUrl(Ins,g,to)
		if err!= nil {
			fmt.Println(err)
		}
	}
	return key,nil
}

func (self *ApiInfo)GetCandlesUrl(Ins *Instrument,granularity string,to  int64) (key uint64,err error) {
//	_tmp :=fmt.Sprintf("%d%s",to,granularity)
//	key =  self.tmpCandles[_tmp]
//	if key != nil {
//		return key,nil
//	}
	path := fmt.Sprintf("%s/instruments/%s/candles?",self.Host,Ins.Name)
	uv := url.Values{}
	uv.Add("granularity",granularity)
	uv.Add("price","M")
	uv.Add("smooth","True")
	uv.Add("count",fmt.Sprintf("%d",GranLen))
//	uv.Add("from",fmt.Sprintf("%d",from))
	uv.Add("to",fmt.Sprintf("%d",to))
	path += uv.Encode()
//	fmt.Println(path)
	da := make(map[string]interface{})
	err = self.ClientDO(path,&da)
	if err != nil {
		return 0,err
	}
	ca := da["candles"].([]interface{})
	lc := len(ca)
//	fmt.Println(lc)
	if lc == 0 {
		return 0,fmt.Errorf("candles len = 0")
	}
	X :=make([]float64,lc)
	Y :=make([]float64,lc)
	var firstDiff,firstTime float64
	var can *Candles
	O := math.Pow(10,Ins.DisplayPrecision)
	for i,c := range ca {
		can = new(Candles)
		can.Init(c.(map[string]interface{}))
//		fmt.Println(can.Time)
		if i==0 {
			firstDiff = can.Mid[0]
			firstTime = float64(can.Time.Unix())
			X[i] = 0
		}else{
			X[i] = float64(can.Time.Unix()) - firstTime
		}
		Y[i] = Round((can.Mid[1] - firstDiff) * O)
	}
	key = fitting.CurveFittingArr(X,Y,MaxLevel)
//	self.tmpCandles[_tmp] = key
	return key,nil
	//MaxLevel
}
func (self *ApiInfo)GetCandlesSampleExp(Ins *Instrument,granularity string,from *int64,arrCan []*Candles,start,begin,tmpbegin int) (Cans []*Candles,err error) {
	path := fmt.Sprintf("%s/instruments/%s/candles?",self.Host,Ins.Name)
	uv := url.Values{}
	uv.Add("granularity",granularity)
	uv.Add("price","M")
	uv.Add("count","5000")
	uv.Add("smooth","True")
	uv.Add("from",fmt.Sprintf("%d",*from))
//	uv.Add("to",fmt.Sprintf("%d",to))
	path += uv.Encode()

//	fmt.Println(path)

	da := make(map[string]interface{})
	err = self.ClientDO(path,&da)
	if err != nil {
		return nil,err
	}
	ca := da["candles"].([]interface{})
	lc := len(ca)
	if lc == 0 {
		return nil,fmt.Errorf("candles len = 0")
	}
	arrlen := len(arrCan)
	SumLen := lc+arrlen
	var cans []*Candles = make([]*Candles,SumLen)
	if arrlen >0 {
		for i,c := range arrCan {
			cans[i] = c
		}
	}
	var h,l float64 = 0,0
	var dif,lastDif float64 = 0,0

	var xOne float64 = 0
	for I,c := range ca {
		i := I + arrlen
		cans[i] = new(Candles)
		cans[i].Init(c.(map[string]interface{}))
		self.GetWaitCount(cans[i].Mid[1])
		dif = cans[i].Mid[1] - cans[i].Mid[0]
		if  (dif != 0 ) && (lastDif != 0 ) && ((dif >0) == (lastDif >0)) {
			if tmpbegin < 0 {
				tmpbegin = i-1
				h = cans[tmpbegin].Mid[2]
				l = cans[tmpbegin].Mid[3]
			}
			if h < cans[i].Mid[2] {
				h = cans[i].Mid[2]
			}
			if l > cans[i].Mid[3] {
				l = cans[i].Mid[3]
			}
			lastDif = dif
			if begin == tmpbegin {
				x := math.Abs(cans[i].Mid[1] - cans[begin].Mid[0])/float64(cans[i].Time.Unix() - cans[begin].Time.Unix())
				if x > xOne {
					xOne = x
				}
			}else{
				if (h - l) > Ins.MinimumTrailingStopDistance {
					begin = tmpbegin
					if start <0 {
						start = begin
					}
					xOne = math.Abs(cans[i].Mid[1] - cans[begin].Mid[0])/float64(cans[i].Time.Unix() - cans[begin].Time.Unix())
				}
			}
			continue
		}
		h = 0
		l = 0
		tmpbegin = -1
		if dif != 0 {
			lastDif = dif
		}
		if start > 0  {
			x := math.Abs(cans[i].Mid[1] - cans[begin].Mid[0])/float64(cans[i].Time.Unix() - cans[begin].Time.Unix())
			if x < xOne/2 {
				//over sample
				if start > GranLen {
					err = self.GetSampleKey(cans[start-GranLen:i+1],Ins)
					if err != nil {
						fmt.Println(err)
					}
				}
				start = -1
			}
		}
	}
	fmt.Println(cans[arrlen].Time,cans[SumLen-1].Time,Ins.Name)
	*from = cans[SumLen-1].Time.Unix()+5
	if *from > time.Now().Unix() {
		return nil,fmt.Errorf("> now")
	}
	var cs []*Candles
	if start > GranLen {
		cs= cans[start-GranLen:]
		begin =begin- start+GranLen
		tmpbegin = tmpbegin - start+GranLen
		start = GranLen
		return self.GetCandlesSampleExp(Ins,granularity,from,cs,start,begin,tmpbegin)
	}
	if tmpbegin > GranLen {
		cs= cans[tmpbegin-GranLen:]
		tmpbegin = GranLen
		return self.GetCandlesSampleExp(Ins,granularity,from,cs,-1,-1,tmpbegin)
	}
	cs = cans[SumLen - GranLen:]
	return cans[SumLen - GranLen:],nil


}
func Round(v float64) (t float64) {
	t = math.Floor(v)
	if v - t > 0.5 {
		return t+1
	}
	return t

}
func (self *ApiInfo)GetSampleKey(cans []*Candles,Ins *Instrument) error {
	O := math.Pow(10,Ins.DisplayPrecision)
	le := len(cans)
	X := make([]float64,le)
	Y := make([]float64,le)
	X[0] = 0
	firstDiff := cans[0].Mid[0]
	Y[0] = Round((cans[0].Mid[1] - firstDiff) * O)
	firstTime := cans[0].Time.Unix()
	keys ,err := self.GetCandlesList(Ins,firstTime)
	if err != nil {
		return err
//		panic(err)
	}
	I := 0
	MaxH:= cans[0].Mid[2]
	MinL:= cans[0].Mid[3]
	for i,c := range cans[1:] {

		if MaxH < cans[i].Mid[2] {
			MaxH = cans[i].Mid[2]
		}
		if MinL > cans[i].Mid[3] {
			MinL = cans[i].Mid[3]
		}
		I = i+1
		X[I] = float64(c.Time.Unix() - firstTime)
		Y[I] = Round((c.Mid[1] - firstDiff) * O)
	}
	key := fitting.CurveFittingArr(X,Y,MaxLevel)
	keys = append([]uint64{key},keys...)

	Gr := new(Grain)
	Gr.List = keys
	Gr.Over = cans[len(cans)-1].Mid[1]
	Gr.Tmp = MaxH - MinL
	Gr.Complete = false
	Gr.Max = cans[len(cans)-1].Time.Unix()
	T,err := self.TreeL.Find(Gr.List)
	if err == nil {
		if T {
			Gr.T = 1
		}else{
			Gr.T = 2
		}
	}

	self.WaitSample = append(self.WaitSample,Gr)
//	fmt.Println(Gr)
//	panic(9)
	return nil
}
