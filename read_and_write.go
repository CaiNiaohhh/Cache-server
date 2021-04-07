/*
1. 先写id 1-10000 的数据到服务器
2. 然后读 1-10000 和写 10001-20000 同时进行
3. 记录运行结束的时间
*/
package main

import(
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup //定义一个同步等待的组

func Send_Message_First(message string) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Dial err: ", err)
		return 
	}
	defer conn.Close()
	conn.Write([]byte(message))
	// 接收数据
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read err: ", err)
		return 
	}
	fmt.Println(string(buf[:n]))
}

func Send_Message(message string) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Dial err: ", err)
		return 
	}
	defer conn.Close()
	conn.Write([]byte(message))
	// 接收数据
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read err: ", err)
		return 
	}
	fmt.Println(string(buf[:n]))
	wg.Done() //减去一个计数
}


func main() {
	var id, _id, body, size_of_string string
	var l int
	// 先写入10000条数据,数据大小为255Byte
	buf_string := ""
	for i := 0; i < 255; i += 1 {
		buf_string += "#"
	}
	Send_Message_First("TYPE_LIST_OPERATION|1-10001|"+buf_string)
	// 开始计时
	t_start := time.Now() // get current time
	// 读写交替
	for i := 1; i <= 10000; i += 1 {
		j := 10000 + i
		wg.Add(1) //添加1个计数
		wg.Add(1) //添加1个计数
		id = fmt.Sprintf("%d", j)
		_id = fmt.Sprintf("%d", i)
		body = buf_string
		l = len(body)
		size_of_string = fmt.Sprintf("%d", l)
		if i % 2 == 0 {
			go Send_Message(id + "|" + body + "|" + size_of_string)	
			Send_Message(_id)
		} else {
			go Send_Message(_id)	
			Send_Message(id + "|" + body + "|" + size_of_string)	
		}	
	}
	wg.Wait() //阻塞直到所有任务完成
	tol := time.Since(t_start)
	fmt.Println("Spend ", tol, " second")

}
