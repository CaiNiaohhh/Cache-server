import socket
from concurrent.futures import ThreadPoolExecutor
import time

class Node:
    '''
    pre : 保存前向指针
    pnext : 保存后向指针
    '''
    def __init__(self, id="", body="", size=0, pre=None, pnext=None):
        self.id = id
        self.body = body
        self.size = size
        self.pre = pre
        self.pnext = pnext
    # 将新来的节点插入到该节点的前面
    def Insert(self, new_node):
        global head, tail
        self.pre = new_node
        new_node.pnext = self
        head = new_node

# 记录头指针
head = Node()
# 记录尾指针
tail = Node()

def id_in_list(node):
    global head, tail
    if node == head:
        return
    if node ==tail:
        tmp_tail = tail
        tmp_tail.pnext = head
        head.pre = tmp_tail
        head = tmp_tail
        tail.pre.pnext = None
        tail = tail.pre
        return
    node.pre.pnext = node.pnext
    node.pnext.pre = node.pre
    node.pnext = head
    head.pre = node
    head = node

def Search_Id(id):
    global head, tail
    tmp = head
    while True:
        if tmp == tail:
            if id == tail.id:
                return tmp.body
            return "FFFFFFFF"
        if tmp.id == id:
            id_in_list(tmp)
            return tmp.body
        tmp = tmp.pnext

# 根据接收到的字符串 判断是读还是写 用 | 作为区别标志
# 读的话只有id，没有 | 
# 写的话有两个 | 即 id|body|size
def is_read(content):
    if content.find("|") == -1:
        return True
    return False

def is_list_operation(content):
    if content.find("TYPE_LIST_OPERATION") != -1:
        return True
    return False

def Cut_String(content):
    k = content.find("|")
    id = content[:k]
    res = content[k + 1:]
    k = res.find("|")
    body = res[:k]
    size = res[k + 1:]
    return id, body, size

def cut_list_str(content):
    k = content.find("|")
    tmp = content[k + 1:]
    k = tmp.find("-")
    id_start = tmp[:k]
    tmp = tmp[k + 1:]
    k = tmp.find("|")
    id_end = tmp[:k]
    body = tmp[k + 1:]
    return id_start, id_end, body

def Handle_conn(conn):
    global head, tail
    client_data = conn.recv(1024).decode()
    # print("receive client send message of ", client_data)
    if is_list_operation(client_data):
        id_start, id_end, body = cut_list_str(client_data)
        id_start = int(id_start)
        id_end = int(id_end)
        l = len(body)
        for i in range(id_start, id_end):
            new_node = Node()
            new_node.id = str(i)
            new_node.body = body
            new_node.size = str(l)
            head.Insert(new_node)
        res = str(id_start) + "-" + str(id_end) + " ok"
        conn.sendall(res.encode())    # 回馈信息给客户端
    elif is_read(client_data):
        res = Search_Id(client_data)
        conn.sendall(res.encode())
    else:
        id, body, size = Cut_String(client_data)
        if head.id == "":
            head.id = id
            head.body = body
            head.size = size
            tail = head
        else:
            new_node = Node()
            new_node.id = id
            new_node.body = body
            new_node.size = size
            head.Insert(new_node)
        res = id + " ok"
        conn.sendall(res.encode())
    conn.close()    # 关闭连接

if __name__ == '__main__':
    ip_port = ('127.0.0.1', 8080)
    sk = socket.socket()            # 创建套接字
    sk.bind(ip_port)                # 绑定服务地址
    sk.listen(50)                    # 监听连接请求
    i = 0
    while True:
        i = i + 1
        # print("wait for the ", str(i), " time of connect")
        conn, address = sk.accept()     # 等待连接，此处自动阻塞
        with ThreadPoolExecutor(max_workers=50) as task:
            task.submit(Handle_conn, conn)
        # Handle_conn(conn)

    # wait(all_task, return_when=ALL_COMPLETED)
    # print("done")