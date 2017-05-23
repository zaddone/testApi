package tree
import(
	"fmt"
	"math"
	"sync"
)
const (
//	TreeLen int  = 1024
//	TreeLen int  = 16
//	NextLen int = 13
	MaxLen int = 3
	MAX float64 =5
//)
//var (
	Kv float64 = 0.05
)
type Node struct {
	T []*Tree //32
}
func (self *Node)Shows() {
	str := ""
	count :=0
	for i,t := range self.T {
		if t!= nil {
			str = fmt.Sprintf("%s (%d %d)",str,i,t.Val)
			count += int(t.Val[0])+int(t.Val[1])
		}
	}
	fmt.Println(str,count);
}
func (self *Node)Show(count * int) {

	if len(self.T) == 0 {
		return
	}
	for _,t := range self.T {
		if t != nil {
//			fmt.Println(t.Val)
			*count++
			t.Show(count)
		}
	}

}
func GetJeiTi(num uint64,level uint) int {
	return int(num >> level)
//	if level == 0 {
//		return int(num)
//	}
//	bc:=uint(64-level)
//	return int(num<<bc>>bc)
}
func powInt(v int,l int) int {
	if l == 0 {
		return 1;
	}
	for i:=1;i<l;i++{
		v *=v
	}
	return v
}
func (self *Node)Append(data []uint64,index int,level int,num int,val int,count * int){

	if len(self.T) == 0 {
	//	self.T = make([]*Tree,powInt(2,num - level))
		self.T = make([]*Tree,powInt(2,num))
	}
//	nu:=GetJeiTi(data[index],uint(level))
	nu := data[index]
//	fmt.Println(num,level)
	if self.T[nu] ==nil {
		if count != nil {
			*count ++
		}
		self.T[nu] = new(Tree)
	}
	self.T[nu].Append(data,index,level,num,val,count)

}
func (self *Node) Find(data []uint64,index int,Vals []int,level int) {
	if len(self.T) == 0 {
	//	if level == 1 {
	//		Vals[2]+=len(data)-index
	//		return
	//	}else {
	//		Vals[2]++
	//	}
		return
	}
//	num:=GetJeiTi(data[index],uint(level))
	num:=data[index]
	if self.T[num] ==nil {
		if level == 1 {
			Vals[2]+=len(data)-index-1
		}else if level == 2 {
			Vals[2]++
		}
		return
	}
	self.T[num].Find(data,index,Vals,level)
}
type Tree struct {
//	Val  []uint8
	Val  []int
	Next  []*Node
	level int
}
func (self *Tree) Show(count *int){
	if len(self.Next) == 0 {
		return
	}
	for _,n := range self.Next {
		if n != nil {
			n.Show(count)
		}
	}
}
func (self *Tree) AppendVal(val int){
	if self.Val == nil {
		self.Val = make([]int,2)
	}
	self.Val[val]++

	v := 1 - val
	if self.Val[val] == 255 {
		self.Val[val] = 15
		self.Val[v]  /= 17
	}else if self.Val[val] == self.Val[v] {
		self.Val[val] = 0
		self.Val[v] = 0
	}else if self.Val[val] > int(MAX) {
		self.Val[val] = self.Val[val] - self.Val[v]
		self.Val[v] = 0
	}
}

func (self *Tree)Append(data []uint64,index int,level int,num int,val int,count * int){
	L:=len(data)
	index++
	level++
//	self.level = level
//	self.AppendVal(val)
	if index == L || level == MaxLen {
		self.AppendVal(val)
		return
	}

	if len(self.Next) == 0 {
		self.Next = make([]*Node,L-index)
	}
	for i:=index; i < L; i++{
		I := i - index
		if self.Next[I] == nil {
			self.Next[I] = new(Node)
		}
		self.Next[I].Append(data,i,level,num,val,count)
	}

}
func (self *Tree) AddVal(Vals []int) {
	if self.Val == nil {
		Vals[2]++
		return
	}
	v1:= float64(self.Val[1])
	v0:= float64(self.Val[0])
	sum :=v1 + v0
//	if sum < MAX {
//		return
//	}

	dif :=v1 - v0
	if sum<MAX || math.Abs(dif)/sum < Kv {
		Vals[2]++
		return
	}
	if dif >0 {
		Vals[1]++
	}else{
		Vals[0]++
	}
}

func (self *Tree) Find(data []uint64,index int,Vals []int,level int) {
	L:=len(data)
	level++
	if level == MaxLen {
		self.AddVal(Vals)
		return
	}
	if len(self.Next) == 0 {
		return
	}
	index++
	for i:=index; i < L; i++{
		I := i - index
		if self.Next[I] == nil {
		//	if level == 1 {
		//		Vals[2]+=L-i
		//	}else if level == 2 {
		//		Vals[2]++
		//	}

			continue
		}
		self.Next[I].Find(data,i,Vals,level)
	}
	return
}
type TreeList struct {
	Tr *Node
//	Name string
	Count int
	Num int
	sync.RWMutex
	Test  [3]float64
	Estimat  [3]float64
}
func (self *TreeList) Init(n int){
	self.Tr = new(Node)
	self.Num = n
}

func (self *TreeList) Show(count *int){
	self.Tr.Show(count)
}
func (self *TreeList)Tests(data []uint64,val int){
	b,err:= self.Find(data)
	if err != nil {
		self.Test[2]++
		return
	}
	if b == (val==1) {
		self.Test[0]++
		return
	}else{
		self.Test[1]++
		return
	}
}

func (self *TreeList) AppendWait(data []uint64,val int) error {
	self.Lock()
	defer self.Unlock()
	return self.Append(data,val)
}
func (self *TreeList) Append(data []uint64,val int) error {

	self.Count ++
	self.Tr.Append(data[0:],0,0,self.Num,val,nil)
	return nil

}
func (self *TreeList) FindWait(data []uint64) (bool,error) {
	self.Lock()
	defer self.Unlock()
	return self.Find(data)
}
func (self *TreeList) Find(data []uint64) (bool,error) {
	val := make([]int,3)
	self.Tr.Find(data[0:],0,val,0)

	sum := val[1] + val[0]
//	fmt.Println(sum+val[2])

//	if sum == 0 {
//		return false,fmt.Errorf("val = 0")
//	}
//	fmt.Println(sum +val[2],len(data))
//	if (sum +val[2]) < len(data)-1 {
//	if val[2] == 0 {
//		return false,fmt.Errorf("val2 =0")
//	}


	if sum < 3{
		return false,fmt.Errorf("val <= %d %d", val[2],sum+val[2])
	}
	dif := val[1] - val[0]
//	if math.Abs(float64(dif))/float64(sum) < Kv{
	if math.Abs(float64(dif)) != float64(sum)  {
		return false,fmt.Errorf("dif != sum %d",sum+val[2])
	}

////	fmt.Println(val[0]+val[1])
//	if val[2] > sum {
//		return false,fmt.Errorf("val > kv")
//	}
//	fmt.Println(dif)
//	fmt.Println(val)
	return dif > 0,nil


}
