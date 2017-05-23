package fitting
import(
	"testing"
	"os"
	"fmt"
)
func Test_CurveFitting(t *testing.T){

	file := "../dba/EURGBP.m/20170426/1493192260199/11"
	fi,err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	Size := int(fi.Size())
	data:=make([]byte,Size)
	f,err := os.Open(file)
	if err != nil {
		panic(err)
	}
	begin:=0
	for{
		n,err :=f.Read(data[begin:]);
		begin+=n
		if begin ==  Size{
			break
		}
		if err.Error() == "EOF" {
			break
		}
	}
	key,err := CurveFitting(data,4)
	fmt.Println(key,Size,err)

}
