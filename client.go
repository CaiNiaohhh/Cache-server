package main

import (
   "net"
   "fmt"
)

// 编写一个发送请求的函数
func Send_Message(message string){
   // 主动发起连接请求
   conn, err := net.Dial("tcp", "127.0.0.1:8080")
   if err != nil {
      fmt.Println("Dial err:", err)
      return
   }
   defer conn.Close()         // 结束时，关闭连接
   // 发送数据
   conn.Write([]byte(message))
   
   // 接收服务器返回的数据
   buf := make([]byte, 1024)
   n, err := conn.Read(buf)
   if err != nil {
      fmt.Println("read err: ", err)
      return
   }
   fmt.Println("client receive the message from server, the content is: ", string(buf[:n]))

}

func main() {
   for {
      var str string
      var id string
      var _id string
      var body string
      var size_of_string string
      var size int
      fmt.Printf("请输入1或2或3 【1表示写数据 2表示读数据 3表示批量写入数据】： ")
      fmt.Scanln(&str)
      if str == "1" {
         fmt.Printf("请输入id: ")
         fmt.Scanln(&id)
         fmt.Printf("请输入具体内容: ")
         fmt.Scanln(&body)
         size = len(body)
         size_of_string = fmt.Sprintf("%d", size)
         Send_Message(id + "|" + body + "|" + size_of_string)
      } else if str == "2" {
         fmt.Printf("请输入id: ")
         fmt.Scanln(&id)
         Send_Message(id)
      } else if str == "3" {
         fmt.Printf("请输入起始id: ")
         fmt.Scanln(&id)
         fmt.Printf("请输入结束id: ")
         fmt.Scanln(&_id)
         fmt.Printf("请输入具体内容: ")
         fmt.Scanln(&body)
         Send_Message("TYPE_LIST_OPERATION|" + id + "-" + _id + "|" + body)
      } else {
         fmt.Println("输入有误 请重新输入")
         continue
      }
   }
}