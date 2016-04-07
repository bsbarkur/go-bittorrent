package main

/*
This implements the bencoding required for BitTorrent
http://www.bittorrent.org/beps/bep_0003.html
*/
import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

/*
Strings are length-prefixed base ten followed by a colon and the string.
For example 6:foobar corresponds to 'foobar'.
*/
func bencode_string(x string) string {
	var stringLength int = utf8.RuneCountInString(x)
	var buffer bytes.Buffer

	/* http://herman.asia/efficient-string-concatenation-in-go */
	buffer.WriteString(strconv.Itoa(stringLength))
	buffer.WriteString(":")

	for i := 0; i < stringLength; i++ {
		buffer.WriteString(string(x[i]))
	}

	return buffer.String()
}

/*
Integers are represented by an 'i' followed by the number in base 10 followed
by an 'e'. For example i3e corresponds to 3 and i-3e corresponds to -3.
Integers have no size limitation. i-0e is invalid. All encodings with a leading
zero, such as i03e, are invalid, other than i0e, which of course
corresponds to 0.
*/
func bencode_int(x int) string {
	var buffer bytes.Buffer

	buffer.WriteString("i")
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(x))
	buffer.WriteString("e")

	return buffer.String()
}

/*
Lists are encoded as an 'l' followed by their elements (also bencoded)
followed by an 'e'. For example l4:spam4:eggse corresponds to ['spam', 'eggs'].
*/
func bencode_list(args []string) string {
	var buffer bytes.Buffer
	buffer.WriteString("l")

	for _, arg := range args {
		buffer.WriteString(strconv.Itoa(utf8.RuneCountInString(arg)))
		buffer.WriteString(":")
		buffer.WriteString(arg)
	}
	buffer.WriteString("e")

	return buffer.String()
}

/*
Dictionaries are encoded as a 'd' followed by a list of alternating keys and
their corresponding values followed by an 'e'. For example,
d3:cow3:moo4:spam4:eggse corresponds to {'cow': 'moo', 'spam': 'eggs'} and
d4:spaml1:a1:bee corresponds to {'spam': ['a', 'b']}. Keys must be strings and
appear in sorted order (sorted as raw strings, not alphanumerics).
*/
type BencodedDict struct {
	strValue string
}

/* Return string of type BencodedDict */
func (bdict *BencodedDict) bencoded_dict_string() string {
	return "d" + bdict.strValue + "e"
}

/* Add key with value of type int */
func (bdict *BencodedDict) add_int(key string, value int) {
	bdict.strValue += bencode_string(key) + bencode_int(value)
}

/* Add key with value of type string */
func (bdict *BencodedDict) add_string(key string, value string) {
	bdict.strValue += bencode_string(key) + bencode_string(value)
}

/* Add key with value of type list */
func (bdict *BencodedDict) add_list(key string, value []string) {
	bdict.strValue += bencode_string(key) + bencode_list(value)
}

/* Add key with value of type dict */
func (bdict *BencodedDict) add_dict(key string, value BencodedDict) {
	bdict.strValue += bencode_string(key) + value.strValue
}

func main() {
	var encodedString string
	var encodedInteger string
	var encodedList string

	encodedString = bencode_string("foobar")
	fmt.Println(encodedString)

	encodedInteger = bencode_int(123)
	fmt.Println(encodedInteger)

	listArray := []string{"spam", "eggs", "baz"}
	encodedList = bencode_list(listArray[0:])
	fmt.Println(encodedList)

	dict := BencodedDict{}
	dict.add_int("foobar", 42)
	fmt.Println(dict.bencoded_dict_string())
}
