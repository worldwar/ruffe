package types

import(
    "testing"
    "reflect"
    "bytes"
)

func Equal(a interface{}, b interface{}) bool {
    if a == nil && b == nil {
        return true
    }
    if reflect.TypeOf(a) != reflect.TypeOf(b) {
        return false
    }
    if a == nil {
        return b == nil   
    }
    switch a.(type) {
        case *string:
            if a.(*string) == nil || b.(*string) == nil {
                return a == b
            }
            return *(a.(*string)) == *(b.(*string))
        case int:
            return a == b
        case []interface{}:
            return EqualArray(a.([]interface{}), b.([]interface{})) 
    }
    return false 
}

func EqualArray(a []interface{}, b []interface{}) bool {
    if len(a) == len(b) {
        for i, x := range a {
            if !Equal(x, b[i]) {
                return false
            }
        }
        return true
    }
    return false
}


func TestDecode(t *testing.T) {
    //nilString *string = nil

    var tests = []struct{
        input string
        want interface{}
    }{
        {"+OK\r\n", Pointer("OK")},
        {":12345\r\n", 12345},
        {":-12345\r\n", -12345},
        {"$6\r\nzhuran\r\n", Pointer("zhuran")},
        {"$0\r\n\r\n", Pointer("")},
        {"$-1\r\n", (*string)(nil)},
        {"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", []interface{}{Pointer("foo"), Pointer("bar")}},
        {"*0\r\n\r\n", []interface{}{}},
        {"*3\r\n:1\r\n:2\r\n:3\r\n", []interface{}{1, 2, 3}},
        {"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n", []interface{}{1, 2, 3, 4, Pointer("foobar")}},
        {"*-1\r\n", ([]interface{})(nil)},
        {"*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n", []interface{}{Pointer("foo"), (*string)(nil), Pointer("bar")}},
    }
    for _, test := range tests {
        if !Equal(Decode(bytes.NewBufferString(test.input)), test.want) {
            t.Errorf("Decode(%q) = %v, want %v", test.input, Decode(bytes.NewBufferString(test.input)), test.want)
        }
    }
}

func TestEncode(t *testing.T) {
    var tests = []struct{
        input interface{}
        want *string
    }{
        {"OK", Pointer("$2\r\nOK\r\n")},
        {12345, Pointer(":12345\r\n")},
        {-12345, Pointer(":-12345\r\n")},
        {(*string)(nil), Pointer("$-1\r\n")},
        {[]string{"foo", "bar"}, Pointer("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")},
        {[]interface{}{}, Pointer("*0\r\n")},
        {[]interface{}{1, 2, 3}, Pointer("*3\r\n:1\r\n:2\r\n:3\r\n")},
        {[]interface{}{1, 2, 3, 4, Pointer("foobar")}, Pointer("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n")},
        {([]interface{})(nil), Pointer("*-1\r\n")},
        {[]interface{}{Pointer("foo"), (*string)(nil), Pointer("bar")}, Pointer("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")},
    }

    for _, test := range tests {
        if !Equal(Encode(test.input), test.want) {
            t.Errorf("Encode(%v) = %v, want %v", test.input, *Encode(test.input), *test.want)
        }
    }
}
