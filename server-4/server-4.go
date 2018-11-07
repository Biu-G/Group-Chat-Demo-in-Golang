package main
import(
	"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"net"
	"log"
	"io"
)
type client chan string
var(
	messages=make(chan string)
	leaving=make(chan client)
	entering=make(chan client)
)
func checker(err error){
	if err!=nil{
		log.Fatal(err)
	}
}
func WriteToClients(conn net.Conn,ch chan string){
	for msg:= range ch{
		conn.Write([]byte(msg))
	}
}
func HandleConn(conn net.Conn){
	defer conn.Close()
	db,err:=sql.Open("mysql","uroot:xiao201379@tcp(127.0.0.1:3306)/tester")
	checker(err)
	err=db.Ping()
	checker(err)
	defer db.Close()
	var name string
	var pass string
	byt:=make([]byte,1024)
	ch:=make(chan string)
	go WriteToClients(conn,ch)
	for{
		n,err:=conn.Read(byt)
		if err!=nil{
			if n==0&&err==io.EOF{
				str:="QUIT EVEN NOT LOGGED!"
				fmt.Println(str)
				ch<-str
				close(ch)
				conn.Close()
				return 
			}else{
				checker(err)
			}
		}else{
			str:=string(byt[:n])
			if(n==9&&str[:1]=="l"){
				strs:="LOGGING IN"
				fmt.Println(strs)
				ch<-strs
				tx,err:=db.Begin()
				checker(err)
				defer tx.Rollback()
				myname:=str[1:5]
				mypass:=str[5:9]
				err=tx.QueryRow("select pass from players where name=?",myname).Scan(&pass)
				if err!=nil{
					if err==sql.ErrNoRows{
						strs:="NON-EXIST USER"
						fmt.Println(strs)
						ch<-strs
					}else{
						checker(err)
					}
				}else{
					if pass==mypass{
						strs:="WELCOME BACK->"+myname+"!"
						name=myname
						fmt.Println(strs)
						ch<-strs
						tx.Commit()
						break
					}else{
						strs:="WRONG PASSWD"
						fmt.Println(strs)
						ch<-strs
					}
				}
				
			}else if(n==9&&str[:1]=="s"){
				strs:="SIGNING UP"
				fmt.Println(strs)
				ch<-strs
				tx,err:=db.Begin()
				checker(err)
				defer tx.Rollback()
				myname:=str[1:5]
				mypass:=str[5:9]
				err=tx.QueryRow("select pass from players where name=?",myname).Scan(&pass)
				if err==nil{
					strs:="ALREADY REGISTERED!!"
					fmt.Println(strs)
					ch<-strs
				}else{
					if err!=sql.ErrNoRows{
						checker(err)
					}else{
						name=myname
						pass=mypass
						strs:="WELCOME NEW!->"+name
						fmt.Println(strs)
						ch<-strs
						_,err:=tx.Exec("insert into players values(?,?)",name,pass)
						checker(err)
						tx.Commit()
						break
					}
				}
			}else{
				strs:="GET AN ID FIRST"
				fmt.Println(strs)
				ch<-strs
			}
		}
	}
	//ID CONFIRMED!!~
	entering<-ch
	messages<-"NEWCOMER->"+name
	for{
		n,err:=conn.Read(byt)
		if err!=nil{
			if n==0&&err==io.EOF{
				leaving<-ch
				messages<-name+"->LEFT"
				fmt.Println(name+"->LEFT")
				break
			}else{
				checker(err)
			}
		}else{
			strs:=string(byt[:n])
			strs=name+"->"+strs
			messages<-strs
		}
	}
}
func broadcast(){
	clients:=make(map[client]bool)
	for{
		select{
		case msg:=<-messages:
			for cli:= range clients{
				cli<-msg
			}
		case cli:=<-leaving:
			delete(clients,cli)
			close(cli) 
		case cli:=<-entering:
			clients[cli]=true
		}
	}
}
func main(){
	listen,err:=net.Listen("tcp",":2077")
	checker(err)
	fmt.Println("LISTENING..")
	go broadcast()
	for{
		conn,err:=listen.Accept()
		if err!=nil{
			fmt.Println("one client error")
			continue
		}
		go HandleConn(conn)
	}
}