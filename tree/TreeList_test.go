package tree
import(
	"testing"
//	"math/rand"
	"fmt"
)
//func TestTree(t *testing.T) {
//
//	tree:=new(TreeList)
//	tree.Init(NextLen)
////	trs := make([]*Node,NextLen)
//	for k:=0;k<5;k++{
//		var sample [NextLen]uint64
//		for i,_:=range sample {
//			sample[i] = uint64(rand.Intn(TreeLen))
//		}
//		Val:=rand.Intn(2)
//		fmt.Println(sample)
//		tree.Append(sample[0:],Val)
//	}
//	var sample [NextLen]uint64
//	for i,_:=range sample {
//		sample[i] = uint64(rand.Intn(TreeLen))
//	}
//	Val := rand.Intn(2)
//	b,e :=tree.Find(sample[0:])
//	fmt.Println(b,e,Val)
//
//}
func TestGetJeiTi(t *testing.T){
	var n uint64 =1023
	fmt.Printf("%b\r\n",n)
	for i:=9;i>0;i--{
		fmt.Printf("%d %b\r\n",i,GetJeiTi(n,i))
	}
}
