package main
import(
	"fmt"
	"flag"
	"time"
	"./apiInfo"
)
var (
	Host = flag.String("h","https://api-fxpractice.oanda.com/v3","host")
	Proxy = flag.String("p","127.0.0.1:1080","proxy")
//	Proxy = flag.String("p","","proxy")
	TmpFile = flag.String("t","tmpInfo.log","TmpFile")
	Auth = flag.String("a","d481a687988ea82d0a09c0527221099d-3e9e6c2697e93f64fd58ecf34ad069db","Auth")
	UserID = flag.Int("u",1,"UserID")
	InsName = flag.String("n","EUR_JPY","INS NAME")
	BeginTime = flag.String("b","2015-01-04T08:00:00","2015-01-04T08:00:00")
	Api *apiInfo.ApiInfo
)

func main(){

	flag.Parse()
	Api = new(apiInfo.ApiInfo)
	err := Api.Init(*TmpFile,*Proxy,*Host,*Auth)
	if err != nil {
		panic(err)
	}
	ins,err := Api.GetInstruments(*UserID)
	if err != nil {
		panic(err)
	}
	from,err:=time.Parse("2006-01-02T15:04:05",*BeginTime)
	if err != nil {
		panic(err)
	}
//	to,err:=time.Parse("2006-01-02T15:04:05","2006-01-01T00:00:00")
//	if err != nil {
//		panic(err)
//	}
//	for k,v := range ins {
//		fmt.Println(k,v.MinimumTrailingStopDistance)
//	}
//	count := int64(12)
	f := from.Unix()
//	t := f + Api.Gran["S10"]*count
//	fmt.Println(time.Unix(t,0))
//	key,err := Api.GetCandlesList(ins["EUR_USD"],f,count)
	fmt.Println(f)
	Ins := ins[*InsName]
	var cans []*apiInfo.Candles
	for {
		if f > time.Now().Unix() {
			break
		}
		cans,err = Api.GetCandlesSampleExp(Ins,"S5",&f,cans,-1,-1,-1)
		if err != nil {
			fmt.Println(f,err)
			time.Sleep(time.Second*10)
			continue
//		}else{
//			f = cans[len(cans)-1].Time.Unix()+5
		}
	//	fmt.Println(can)
	}

}
