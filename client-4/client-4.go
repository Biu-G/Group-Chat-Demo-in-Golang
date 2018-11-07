package main
import(
	"fmt"
	"net"
	"log"
	"bufio"
	"os"
	"time"
)
var reader *bufio.Reader
func checker(err error){
	if err!=nil{
		log.Fatal(err)
	}
}
func ListenToMsg(conn net.Conn){
	byt:=make([]byte,1024)
	for{
		n,err:=conn.Read(byt)
		if err!=nil{
			fmt.Println("SERVER REJECTING")
			break
		}else{
		strs:=string(byt[:n])
		fmt.Println(strs)
		}
	}
}
func main(){
	reader=bufio.NewReader(os.Stdin)
	for{
		fmt.Print("Input 'go' to continue->")
		strs,err:=reader.ReadString('\n')
		checker(err)
		strs=strs[:len(strs)-1]
		if strs=="go"{
			fmt.Println("INITIATING..")
			break
		}
	}
	//GO!
	conn,err:=net.DialTimeout("tcp","localhost:2077",2*time.Second)
	checker(err)
	fmt.Println("CONNECTED!")
	defer conn.Close()
	go ListenToMsg(conn)
	for{
		strs,err:=reader.ReadString('\n')
		checker(err)
		strs=strs[:len(strs)-1]
		if strs=="q"||strs=="Q"{
			fmt.Println("BYE")
			break
		}else{
			conn.Write([]byte(strs))
		}
	}
}