1.启动服务端
```
go run main.go server.go user.go

# 或者构建二进制
go build -o server main.go server.go user.go
./server
```

2.启动客户端(可以启动多个进行测试)
```
go run client.go 
```

![image](https://github.com/user-attachments/assets/ff8720ca-245a-43b7-ad36-d8951f22adc8)
