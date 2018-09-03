## V1接口

---

### WebSocket接口

---

>建立连接接口

Url:/ws

---

>总述:

Description:均通过ws.send（）访问

Params:

``` json
{
    "type":"bind"/"send"/"join"/"exit",(sting) //通过该数据判断接口事件
    "uid":,                            (sting) //用户id
    "touid":,                          (sting) //与用户交互的用户id
    "msg":                             (sting) //仅用于"send"时间
}
```

Return:

``` json
{
    "type":"bind"/"send"/"join"/"exit", (sting)  //保留项
    "error":true/false,                 (bool)   //true为错误
    "code":,                            (int)    //错误代码
    "msg":,                             (string) //错误消息/send消息 视情况而定
    "uid":,                             (string) //send事件时发送者的uid
    "touid":,                           (string) //保留项
}
```

保留项暂时均可忽略

---

>bind事件

Description: 推荐在onopen事件请求，用于绑定用户uid

Params:

``` json
{
    "type":"bind", (string)
    "uid":         (string)//绑定的用户id
}
```

Return:

``` json
{
    "error":true/false, (bool)   //true为错误
    "code":,            (int)    //错误代码
    "msg":,             (string) //错误消息
}
```

---

> send事件

Description:发送消息事件

Params:

``` json
{
    "type":"send" (sting)
    "uid":,       (sting) //发送者id(可选)
    "touid":,     (sting) //接收者id
    "msg":        (sting) //具体消息内容
}
```

Return:

``` json
{
    "error":true/false, (bool)   //true为错误
    "code":,            (int)    //错误代码
    "msg":,             (string) //错误消息/send消息 视error而定
    "uid":,             (string) //发送者的uid
}
```

---

> join事件

Description: 均通过ws.send（）访问

Params:

``` json
{
    "type":"join", (sting)
    "uid":,        (sting) //用户id
    "touid":,      (sting) //与用户交互的用户id
}
```

Return:

``` json
{
    "error":true/false, (bool)   //true为错误
    "code":,            (int)    //错误代码
    "msg":,             (string) //错误消息
}
```

---

> exit事件

Description: 均通过ws.send（）访问

Params:

``` json
{
    "type":"exit",(sting)
    "uid":,       (sting) //用户id
    "touid":,     (sting) //与用户交互的用户id
    "msg":        (sting) //仅用于"send"时间
}
```

Return:

``` json
{
    "error":true/false, (bool)   //true为错误
    "code":,            (int)    //错误代码
    "msg":,             (string) //错误消息
}
```