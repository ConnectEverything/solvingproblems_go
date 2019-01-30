# Solving Problems using NATS (Go Edition)

This repository holds the go language templates/solutions for the _Solving Problems using NATS_ courseware.

This project was designed for minimal build dependencies, but does have a go module defined with the go-nats dependency.

The solutions are in `solutions`. Your version of the exercise code should go in `exercises`.

For simplicity, just run the go files directly.

```bash
go run ex01_connecting/main.go
```

## Exercise Hints/Helpers

### General

* The Go doc for the NATS client library is available at [https://godoc.org/github.com/nats-io/go-nats](https://godoc.org/github.com/nats-io/go-nats)
* The Go doc for streaming is available at [https://godoc.org/github.com/nats-io/go-nats-streaming](https://godoc.org/github.com/nats-io/go-nats-streaming)
* Most of the exercises can be implemented by 1 or two main functions, use folders to specify executable names
* You can create a UTF-8 string in go from bytes using `string(theBytes)`
* You can get UTF-8 bytes from a string using `[]byte(theString)`
* You can create random integers with a max using `rand.Int31n(10000)`
* The solutions allow you to put the server URL in the environment

```go
url := os.Getenv("SP_NATS_SERVER")

if url == "" {
    url = "nats://localhost:4222"
}
```

### Exercise 1 - Connecting

* Check out the go doc for options with the name "something" Handler

### Exercise 3 - Notifications

* You can create a timer in Go with:

```go
ticker := time.NewTicker(5 * time.Second)
go func() {
    for range ticker.C {
        // Do something
    }
}()
```

* You can loop over a subscription to get all the messages using:

```go
for {
    msg, err := sub.NextMsg(100 * time.Second)
    // Your code here
}
```

### Exercise 5 & 6 - Mandelbrot

* Use a width of 800 and height of 600 to get a reasonable amount of work but not too much
* Create an image with `img := image.NewRGBA(image.Rect(0, 0, width, height))`
* `palette.Plan9` has 256 colors
* Set a pixel with `img.Set(col, row, color)`
* Save the image to PNG with

```go
out, err := os.Create("mandelbrot.png")
err = png.Encode(out, img)
```

* Each pixel in the Mandelbrot set can be calculated using, with a max iterations of 255:

```go
x0 := (col - width/2.0) * 4.0 / width
y0 := (row - height/2.0) * 4.0 / height

x := 0.0
y := 0.0
iterations := 0

for (x*x+y*y) < 4 && (iterations < max) {
    temp := (x*x - y*y) + x0
    y = (2 * x * y) + y0
    x = temp
    iterations++
}
```

### Exercise 10 - File Copier

* You can read a full file using:

```go
encoded, err := ioutil.ReadFile(filePath)
data := string(encoded)
words := strings.Split(data, " ")
```

### Exercise 15 - Word List

* You can read a file line by line using:

```go
wordList := []string{}
file, err := os.Open(filePath)
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    wordList = append(wordList, scanner.Text())
}
file.Close()
```

Exercise 19 & 20 - Streaming

* You may want to use wait groups to block program exit

```go
wg := sync.WaitGroup{}
wg.Add(1;
```
