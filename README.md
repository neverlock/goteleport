# go-teleport

***Send and receieve data with channel***

transmit data over http protocol with channel interface. 

### Get
```
go get github.com/tspn/goteleport
```

### Example
##### server.go
```go
func main (){
	in, out := goteleport.New(9099, 100)
  
	go func (){
		for {
			v, ok := <- in
			if !ok{
				fmt.Println(ok)
			}

			fmt.Println(string(v.([]byte)))
		}
	}()

	i := 0
	for{
		time.Sleep(5 * time.Second)
		out <- fmt.Sprintf("%s - %d", "this is from server", i)
		i++
	}
}
```

##### client.go
```go
func main(){
	in, out := goteleport.Connect("127.0.0.1:9099", 9100, 100)
	go func (){
		i := 0
		for{
			time.Sleep(3 * time.Second)
			out <- fmt.Sprintf("%s - %d", "this is from client", i)
			i++
		}
	}()

	for {
		if v, ok := <- in; ok{
			fmt.Println(string(v.([]byte)))
		}
	}
	close(in)
}

```
