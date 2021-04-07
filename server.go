/*
1. 构建一个双向链表
2. 当缓存空间达到1G即2 ^ 22时，将链表尾与链表头连接
   形成一个双向循环链表
3. 时刻维护一个head的指针，动态调整
4. lru:当新数据进入时，链表尾需要被替换，head->pre为链表尾
   1）先将旧数据写入磁盘[文件]
   2）将新数据覆盖旧数据 包括id, body, size
   3）将链表尾赋值给head
*/
package main

import (
    "net"
    "log"
    "fmt"
    "os"
    "strings"
    "strconv"   
    "path/filepath"
    "io/ioutil"
)

type Node struct{
	pre *Node  // 前向指针
	next *Node // 后向指针
	id string
	body string
	size string
}

// 最大长度
var Max_Length int = 1024 * 1024 * 4
// 头指针
var head *Node
// 尾指针
var tail *Node
// 记录数据大小
var total int = 0
// 判断时候成环
var is_circle int = 0
// 新建文件个数
var file_count int = 0
// 定义一个文件写入的缓存变量
var file_content string

// 判断目录是否存在
func IsExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}
// 判断文件是否存在
func FileExist(path string) bool {
  _, err := os.Lstat(path)
  return !os.IsNotExist(err)
}
//获取文件的长度，以KB为单位
func getFileSize(filename string) int64 {
    var result int64
    filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
        result = f.Size()
        return nil
    })
    return result
}
// 传入字符串 将字符串写到文件中
func Write_File(content string, filename string) {
	str := []byte(content + "\n")
    // 以追加模式打开文件 文件不存在则创建
    txt, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
    defer txt.Close()
    if err != nil {
        panic(err)
    }
    // 写入文件
    n, err := txt.Write(str)
    // 当 n != len(b) 时，返回非零错误
    if err == nil && n != len(str) {
        println(`错误代码：`, n)
        panic(err)
    }
}

// 将旧数据写入文件
func Save_Into_File(id string, body string, size string){
	if file_count > 200 {
		return
	}
	content := "ID: " + id + " BODY: " + body + " SIZE: " + size
	// 缓存变量达到1M才写入文件中
	if len(file_content) < 1024 * 1024 {
		file_content += (content + "\n")
		return
	}
	// 在当前目录下判断是否存在./oldfile的路径
	if IsExist("./old_message") == false {
		os.Mkdir("./old_message", os.ModePerm)
	}
	// int转string
	name := fmt.Sprintf("%d", file_count)  
	if IsExist("./old_message/" + name + ".txt") == false{
		os.Create("./old_message/" + name + ".txt")
	}
 	// 打开文件 判断文件的内容是否大于10M 
 	filename := "./old_message/" + name + ".txt"
 	if getFileSize(filename) < 1024 * 1024 * 10 {
 		Write_File(file_content, filename)
 	} else {
 		file_count += 1
 		name := fmt.Sprintf("%d", file_count)  
 		filename := "./old_message/" + name + ".txt"
 		Write_File(file_content, filename)
 	}
 	file_content = ""
}


// 链表插入元素
func Insert_Link_List(id string, body string, size string){
	node := &Node{nil, nil, id, body, size}
	if total == 0 {
		head = node
		tail = node
		total += 1
	} else{
		// 如果缓冲区没满
		if total < Max_Length {
			tmp := head
			head = node
			head.next = tmp
			tmp.pre = head
			total += 1
		} else {
			//当缓冲区满了的时候
			// 还没有成环
			if is_circle == 0 {
				tail.next = head
				head.pre = tail 
				is_circle = 1
			}
			// 将旧数据写入文件
			Save_Into_File(tail.id, tail.body, tail.size)
			// 将新数据替换旧数据
			tail.id = node.id
			tail.body = node.body
			tail.size = node.size
			head = tail
			tail = tail.pre
		}
	}
}

// 如果id在缓存中 需要将id提到链表头 然后将id所在的节点删除
func id_in_list(node *Node) {
	// 如果刚好是head
	if node == head {
		return
	}
	// 如果刚好是尾节点
	if node == tail {
		if is_circle == 0 {
			// 提到头节点
			tmp_tail := tail
			tmp_tail.next = head
			head.pre = tmp_tail
			head = tmp_tail
			// 将tail.pre设置为tail
			tail.pre.next = nil
			head.pre = nil
			tail = tail.pre
			return
		} else {
			tmp_tail := tail
			tail = tail.pre
			head = tmp_tail
			return 
		}
	}
	// 删除节点
	node.pre.next = node.next
	node.next.pre = node.pre
	// 提到表头
	node.next = head
	head.pre = node
	head = node
}

// 读取一个文件中的全部内容 并以字符串的形式返回
func read_file(filepath string) string {
	// fmt.Println(filepath)
	file, err := os.Open(filepath)
    if err != nil {
        fmt.Println(err)
        return "FFFFFFFF"
    }
    defer file.Close()
    fileinfo, err := file.Stat()
    if err != nil {
        fmt.Println(err)
        return "FFFFFFFF"
    }

    filesize := fileinfo.Size()
    buffer := make([]byte, filesize)
    bytesread, err := file.Read(buffer)
    if err != nil {
        fmt.Println(err)
        fmt.Println(bytesread)
        return "FFFFFFFF"
    }
    // fmt.Println("bytes read: ", bytesread)
    return string(buffer)
}

// 查找文件id是否在字符串中
func find_string(id string, s string) string {
	for {
        line := strings.Index(s, "\n")
        if line == -1 {
            break
        }
        res := strings.Index(s, "ID: ")
        en := strings.Index(s, " BODY:")
        if res == -1 || en == -1 {
        	return "FFFFFFFF"
        }
        _id := s[res + 4:en]
        si := strings.Index(s, " SIZE:")
        _body := s[en + 6:si]
        if _id == id {
        	l := len(_body)
        	slen := fmt.Sprintf("%d", l)
        	stmp := "id: " + _id + " size: " + slen + " content: " + _body
        	return stmp
        }
        s = s[line + 1:]
    }
    return "FFFFFFFF"
}

// 查找id是否在文件中
func in_file(id string) string {
	// 首先查看一下是否在文件缓存变量中
	res := find_string(id, file_content)
	if res != "FFFFFFFF" {
		return res
	}
    // 查找磁盘中的缓存
    pathname := "./old_message/"
    rd, err := ioutil.ReadDir(pathname)
    if err != nil {
    	fmt.Println("folder doesn't exist")
    	return "FFFFFFFF"
    }
    for _, fi := range rd {
        name := fi.Name()
        Str := read_file(pathname + name)
        res = find_string(id, Str)
        if res != "FFFFFFFF" {
        	return res
        }
    }
    return "FFFFFFFF"
}

// 查找id是否在缓存中
func Search_Id(id string) string {
	tmp := head
	for {
		if tmp == tail{
			if id == tail.id{
				res := "id: " + tail.id + " size: " + tail.size + " content: " + tail.body
				return res
			} else{
				return in_file(id)
			}
		} else {
			if tmp.id == id {
				id_in_list(tmp)
				res := "id: " + tmp.id + " size: " + tmp.size + " content: " + tmp.body
				return res
			} else {
				tmp = tmp.next
			}
		}
	}
}

// 根据接收到的字符串 判断是读还是写 用 | 作为区别标志
// 读的话只有id，没有 | 
// 写的话有两个 | 即 id|body|size
func is_read(content string) bool {
	return !strings.Contains(content, "|")
}

func is_list_operation(content string) bool {
	return strings.Contains(content, "TYPE_LIST_OPERATION")
}

// 将id body size切分出来
func Cut_String(content string) (string, string, string) {
	k := strings.Index(content, "|")
    id := content[:k]
    res := content[k + 1:]
    k = strings.Index(res, "|")
    body := res[:k]
    size := res[k + 1:]
    return id, body, size
}

func cut_list_str(content string) (string, string, string){
	k := strings.Index(content, "|")
    tmp := content[k + 1:]
    k = strings.Index(tmp, "-")
    id_start := tmp[:k]
    tmp = tmp[k + 1:]
    k = strings.Index(tmp, "|")
    id_end := tmp[:k]
    body := tmp[k + 1:]
    return id_start, id_end, body
}

func Handle_conn(i int, conn net.Conn) { 
	// 读取客户端的信息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)	
	if err != nil{
		fmt.Println("Read err:", err)
		return
	}
	content := string(buf[:n])
	// fmt.Println("receive client send message of ", string(buf[:n]))
	// 先判断是读操作还是写操作还是批量写操作
	var res string
	if is_list_operation(content) {
		// "TYPE_LIST_OPERATION|id-id|content"
		id_start, id_end, body := cut_list_str(content)
		id_1, err := strconv.Atoi(id_start)
		id_2, _err := strconv.Atoi(id_end)
		if err != nil || _err != nil {
			fmt.Println("param illegal")
		}
		l := len(body)
		size_od_body := fmt.Sprintf("%d", l)
		for i := id_1; i < id_2; i += 1 {
			tmp := fmt.Sprintf("%d", i)
			Insert_Link_List(tmp, body, size_od_body)
		}
		res = id_start + "-" + id_end + " ok"
		conn.Write([]byte(res))
	} else if is_read(content) {
		res = Search_Id(content)
		conn.Write([]byte(res))//通过conn的wirte方法将数据返回给客户端。
	} else {
		// 写操作
		// 要先判断一下size的大小会不会超过2 ^ 8 Byte
		id, body, size := Cut_String(content)
		size_to_int, err := strconv.Atoi(size)
		if err != nil {
			fmt.Println("param illegal")
		}
		if size_to_int > 256 {
			res = "Maximum length limit exceeded"
			conn.Write([]byte(res))
		} else {
			Insert_Link_List(id, body, size)
			res = "id: " + id + " ok"
			conn.Write([]byte(res))
		}
	}
  	conn.Close() //与客户端断开连接。
}
func main() {
    addr := "0.0.0.0:8080" //表示监听本地所有ip的8080端口
    listener,err := net.Listen("tcp",addr)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    i := 0
    fmt.Println("server is running...")
    for  {
    	// fmt.Println("wait for the ", i, " time of connect")
    	i = i + 1
        conn,err := listener.Accept() //用conn接收链接
        if err != nil {
            log.Fatal(err)
        }   
        go Handle_conn(i, conn)  //开启多个协程。
    }
}