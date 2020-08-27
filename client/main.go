package main

import (
  "bufio"
  "math"
  "flag"
  "fmt"
  "log"
  "math/rand"
  "os"
  "time"
  "strings"

  "golang.org/x/net/websocket"
)
var key string;

type Message struct {
  Text []rune `json:"text"`
}

var (
  port = flag.String("port", "9000", "port used for ws connection")
)

func connect() (*websocket.Conn, error) {
  return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

func mockedIP() string {
  var arr [4]int
  for i := 0; i < 4; i++ {
    rand.Seed(time.Now().UnixNano())
    arr[i] = rand.Intn(256)
  }
  return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}
func getN(key string) int {
  var strSplitEncoded []int
  var a = []rune(key)
  for _, letter := range a {
    strSplitEncoded = append(strSplitEncoded, int(letter))
  }
  var sum = 0
  for _, num := range strSplitEncoded {
    sum = sum + num
  }
  return sum % 94;
}
func chunkSubstr(data string, size int) []string {
  var chunks []string
  var numChunks int 
  var o = 0
  numChunks = int(math.Ceil(float64(len(data))/float64(size)))
  for i := 1 ; i < numChunks; i++ {
    chunks = append(chunks, data[o: o + size])
    o = o + size;
  }
  chunks = append(chunks, data[o:])
  return chunks
}
func reverse(s string) string { 
  rns := []rune(s) // convert to rune 
  for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 { 
    rns[i], rns[j] = rns[j], rns[i] 
  } 
  return string(rns) 
}
func shiftUpByN(val rune, n int) string {
  rns := val // convert to rune 
  if(32 <= int(rns) && int(rns) <= 125) {
    if((int(rns) + n) > 125) {
      return string(31 + ((int(rns) + n) - 125))
    } else {
      return string(int(rns) + n)
    }
  } else {
    return string(val)
  }
}
func shiftDownByN(val rune, n int) string {
  rns := val;
  if(32 <= int(rns) && int(rns) <= 125) {
    if((int(rns) - n) < 32) {
      return string(125 + ((int(rns) - n) - 31))
    } else {
      return string(int(rns) - n)
    }
  } else {
    return string(val)
  }
}

func encrypt(data string, key string) []rune {
  var chunks []string
  var encryptedChunks []string
  var encryptedChunk []string
  var encryptedstring string
  var n int
  var thisChunk string
  n = getN(key)
  chunks = chunkSubstr(data, len(key)) 
  for _, chunk := range chunks {
    thisChunk = reverse(chunk)
    //s := strings.Split(thisChunk, "")
    rns := []rune(thisChunk) // convert to rune 
    for _, letter := range rns {
      encryptedChunk = append(encryptedChunk, shiftUpByN(letter, n))
    }
    thisChunk = reverse(strings.Join(encryptedChunk, ""))
    encryptedChunks = append(encryptedChunks, thisChunk)
  }
  encryptedstring = strings.Join(encryptedChunks, "")
  //fmt.Println(encryptedstring)
  return []rune(encryptedstring)
}

func decrypt(dat []rune, key string) string {
  data := string(dat)
  var chunks []string
  var thisChunk string
  var decryptedChunks []string
  var decryptedChunk []string
  var decryptedstring string
  var n int
  n = getN(key)
  chunks = chunkSubstr(data, len(key)) 
  for _, chunk := range chunks {
    thisChunk = reverse(chunk)
    //s := strings.Split(thisChunk, "")
    rns := []rune(thisChunk) // convert to rune 
    for _, letter := range rns {
      decryptedChunk = append(decryptedChunk, shiftDownByN(letter, n))
    }
    thisChunk = reverse(strings.Join(decryptedChunk, ""))
    decryptedChunks = append(decryptedChunks, thisChunk)
  }
  decryptedstring = strings.Join(decryptedChunks, "")
  return decryptedstring
}


func main() {
  fmt.Println("Enter secret key: ")
  fmt.Scan(&key)
  fmt.Println("You can not start communicating ...")
  fmt.Println("____________________________________________")
  fmt.Println("")

  flag.Parse()

  ws, err := connect()
  if err != nil {
    log.Fatal(err)
  }
  defer ws.Close()
  
  // code to receive
  var m Message
  go func() {
    for {
      err := websocket.JSON.Receive(ws, &m)
      if err != nil {
        fmt.Println("Error receiving message: ", err.Error())
        break
      }
      fmt.Println("Message: ", decrypt(m.Text, key))
    }
  }()


  //code to send messages
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    text := scanner.Text()
    if text == "" {
      continue
    }
    m := Message {
      Text: encrypt(text,key),
    }
    err = websocket.JSON.Send(ws, m)
    if err != nil {
      fmt.Println("Erro sending message: ", err.Error())
    }
  }
}


