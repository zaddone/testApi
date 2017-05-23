package apiInfo
import(
	"testing"
	"fmt"
)
func Test_getInstrumentList(t *testing.T){

	Host := "https://api-fxpractice.oanda.com/v3"
	Proxy := "127.0.0.1:1080"
	TmpFile := "apiInfo.log"
	Auth := "d481a687988ea82d0a09c0527221099d-3e9e6c2697e93f64fd58ecf34ad069db"
	Api := new(ApiInfo)
	err := Api.Init(TmpFile,Proxy,Host,Auth)
	if err != nil {
		panic(err)
	}
	ins,err := Api.GetInstruments(1)
	if err != nil {
		panic(err)
	}
	for _,in := range ins {
		fmt.Println(in.Name)
	}

}
