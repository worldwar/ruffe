package main 

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
    "types"
    "bytes"
)

func Print(in interface{}) string {
    switch in.(type) {
        case int:
          return fmt.Sprintf("(integer) %d", in.(int))
        case *string:
            if in == (*string)(nil) {
                return "nil"
            }
            return fmt.Sprintf("\"%s\"", *(in.(*string)))
        case []interface{}:
            b := new(bytes.Buffer)
            if len(in.([]interface{})) == 0 {
                return "(empty list or set)"
            }
            sep := ""
            for i, x := range in.([]interface{}) {
                b.WriteString(sep)
                sep = "\n"
                b.WriteString(fmt.Sprintf("%d) %s", i + 1, Print(x)))
            }
            return b.String()
        case string:
            return in.(string)
    }
    return ""
}

func send(cmd string, conn net.Conn) {
    go func() {
        raw := make([]byte, 1024)
        length, _ := conn.Read(raw)
        s := raw[:length]
        result := types.Decode(bytes.NewBuffer(s))
        fmt.Fprintln(os.Stdout, Print(result))
    }()
    go func() {
        fmt.Fprint(conn, cmd)
    }()
}

func main() {
    conn, err := net.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("network error")
    }
    fmt.Println("connection established")
    reader := bufio.NewReader(os.Stdin)
    for {
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        if strings.ToLower(text) == "exit" {
            break
        }
        cmds := strings.Split(text, " ")
    	send(*(types.Encode(cmds)), conn)
    }
}
