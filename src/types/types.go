package types

import(
    "bytes"
    "fmt"
    "reflect"
)
func Pointer(s string) *string {
    t := s
    return &t
}

func Encode(i interface{}) *string {
    return EncodeValue(reflect.ValueOf(i))
}

func EncodePointer(v reflect.Value) *string {
    switch v.Type().Elem().Kind() {
        case reflect.String:
            if v.IsNil() {
                return Pointer("$-1\r\n")
            } else {
                return EncodeValue(v.Elem())
            }
     }
     return nil
}

func EncodeValue(v reflect.Value) *string {
    switch v.Kind() {
        case reflect.Int:
            return Pointer(fmt.Sprintf(":%d\r\n", v.Int())) 
        case reflect.String:
            length := len(v.String())
            return Pointer(fmt.Sprintf("$%d\r\n%v\r\n", length, v.String()))
        case reflect.Ptr:
            return EncodePointer(v)
        case reflect.Interface:
            return Encode(v.Interface())
        case reflect.Array, reflect.Slice:
            if v.IsNil() {
                return Pointer("*-1\r\n")
            }
            b := new(bytes.Buffer)
            length := v.Len()
            b.WriteString(fmt.Sprintf("*%d\r\n", length)) 
            for index := 0; index < length; index++ {
                b.WriteString(*EncodeValue(v.Index(index)))
            }
            return Pointer(b.String())
    }
    return nil 
}

func Decode(b *bytes.Buffer) interface{} {
    c, err := b.ReadByte()
    if err == nil {
        switch c {
            case '+':
                return decodeSimpleString(b)
            case ':':
                return decodeInteger(b)
            case '$':
                return decodeBulkString(b)
            case '*':
                return decodeArray(b)
        }
    }
    return nil
}

func decodeSimpleString(b *bytes.Buffer) *string {
    result := new(bytes.Buffer)
    for {
        if c, err := b.ReadByte(); err == nil && c != '\r' {
            result.WriteByte(c)
        } else {
            break
        }
    }
    return Pointer(result.String())
}

func decodeInteger(b *bytes.Buffer) int {
    result := 0
    positive := 0 
    for {
        if c, err := b.ReadByte(); err == nil && c != '\r' {
            if positive == 0 {
                if c == '-' {
                    positive = -1
                } else {
                    positive = 1
                    result = result * 10 + (int(c) - '0')
                }
            } else {
                result = result * 10 + (int(c) - '0')
            }
        } else {
            if c == '\r' {
                b.ReadByte()
            }
            break
        }
    }
    if positive == -1 {
        return -result
    }
    return result
}

func decodeBulkString(b *bytes.Buffer) *string {
    length := decodeInteger(b)
    if length == -1 {
        return nil
    }
    result := make([]byte, length)
    b.Read(result)
    b.Read(make([]byte, 2))
    return Pointer(string(result))
}

func decodeArray(b *bytes.Buffer) []interface{} {
    length := decodeInteger(b)
    if length == -1 {
        return nil
    }
    result := make([]interface{}, length)
    for i := 0; i < length; i++ {
        result[i] = Decode(b)
    }
    return result
}
