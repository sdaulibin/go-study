|类型 |格式 |作用|
|---|---|---|
|测试函数|函数名前缀为Test|测试程序的一些逻辑是否正确|
|基准测试|函数前缀名为Benchmark|测试程序的性能|
|示例函数|||

```Shell
go test
go test -v
go test -run=Sub/chinese
go test -run=Group
go test -cover -coverprofile=split.out
go tool cover -html=split.out 
go test -bench=Split -benchmem
```