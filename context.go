package magic

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// Context structure
// Writer - standert ResponseWriter
// Request - standart Request
// Body - all in body
// Params - params in url (/:id)
// QueryParams - query params in url (?param1=1&param2=2)
// MultipartParams - standart multipartParams without files
// PostParams - all params in post request(form-data)
// FileParams - all files
// Headers - standart headers
// Storage - storage for all; transfer data middleware -> middleware -> ... -> handler
type Context struct {
	Writer          http.ResponseWriter
	Request         *http.Request
	Body            string
	Params          Values
	QueryParams     ValuesArr
	MultipartParams ValuesArr
	PostParams      ValuesArr
	FileParams      FilesArr
	Headers         ValuesArr
	Storage         map[string]interface{}
}

// SendError function
// send your error like {"message": "error message"}
func (context *Context) SendError(err error) error {
	message := make(map[string]interface{})
	message["message"] = err.Error()
	str, _ := json.MarshalIndent(message, "", "    ")
	fmt.Fprint(context.Writer, string(str))
	return err
}

// SendErrorString function
// Send your error like {"message": "your string"}
func (context *Context) SendErrorString(errorStr string) error {
	message := make(map[string]interface{})
	message["message"] = errorStr
	str, _ := json.MarshalIndent(message, "", "    ")
	fmt.Fprint(context.Writer, string(str))
	return errors.New(errorStr)
}

// SendJSON function
// Send pretty json like {"string" : "interface", ...}
func (context *Context) SendJSON(j map[string]interface{}) error {
	str, _ := json.MarshalIndent(j, "", "    ")
	fmt.Fprint(context.Writer, string(str))
	return nil
}

// SendString function
// Send string like "something"
func (context *Context) SendString(str string) error {
	fmt.Fprint(context.Writer, str)
	return nil
}

// SendFile function
// Send os.file by name like "index.html"
func (context *Context) SendFile(fileName string) error {
	http.ServeFile(context.Writer, context.Request, fileName)
	return nil
}

// SendFileBytes function
// Send bytes as file
func (context *Context) SendFileBytes(fileName string, bytes []byte) error {
	mas := strings.Split(fileName, ".")
	if len(mas) != 0 {
		mas[0] = mime.TypeByExtension("." + mas[len(mas)-1])
	} else {
		mas = append(mas, "text/plain")
	}
	map[string][]string(context.Headers)["Content-Type"] = []string{mas[0]}
	return context.SendString(string(bytes))
}

// ParseJSON function
// Parse JSON in body to your interface
func (context *Context) ParseJSON(iface interface{}) error {
	err := json.Unmarshal([]byte(context.Body), iface)
	return err
}

// TEST
// // Bind function
// // Get data from all and put in interface
// func (context *Context) Bind(iface interface{}) error {
// 	values := reflect.ValueOf(iface).Elem()
// 	for i := 0; i < values.NumField(); i++ {
// 		f := values.Field(i)
// 		fmt.Print(values.Type().Field(i).Name)
// 		fmt.Print(" ")
// 		fmt.Print(f.Type())
// 		fmt.Print(" ")
// 		fmt.Println(f.Interface())
// 	}
// 	return nil
// }

// FilesArr structure
// It is map[string][]*multipart.FileHeader
type FilesArr map[string][]*multipart.FileHeader

// Values structure
// It is map[string]string
type Values map[string]string

// ValuesArr structure
// It is map[string][]string
type ValuesArr map[string][]string

// ParseInt function
// Return int and error if can't parse
func (values Values) ParseInt(key string) (int, error) {
	str := values[key]
	i, err := strconv.ParseInt(str, 10, 32)
	return int(i), err
}

// ParseUint function
// Return uint and error if can't parse
func (values Values) ParseUint(key string) (uint, error) {
	str := values[key]
	i, err := strconv.ParseUint(str, 10, 32)
	return uint(i), err
}

// ParseFloat function
// Return float and error if can't parse
func (values Values) ParseFloat(key string) (float32, error) {
	str := values[key]
	i, err := strconv.ParseFloat(str, 32)
	return float32(i), err
}

// ParseBool function
// Return bool and error if can't parse
func (values Values) ParseBool(key string) (bool, error) {
	str := values[key]
	i, err := strconv.ParseBool(str)
	return i, err
}

// ParseString function
// Return string, error always nil
func (values Values) ParseString(key string) (string, error) {
	str := values[key]
	return str, nil
}

// ParseInt function
// Return first int and error if can't parse or empty string
func (valuesArr ValuesArr) ParseInt(key string) (int, error) {
	str := valuesArr[key]
	if len(str) == 0 {
		return 0, errors.New("no this key in map")
	}
	strInt := str[0]
	i, err := strconv.ParseInt(strInt, 10, 32)
	return int(i), err
}

// ParseUint function
// Return first uint and error if can't parse or empty string
func (valuesArr ValuesArr) ParseUint(key string) (uint, error) {
	str := valuesArr[key]
	if len(str) == 0 {
		return 0, errors.New("no this key in map")
	}
	strUint := str[0]
	i, err := strconv.ParseUint(strUint, 10, 32)
	return uint(i), err
}

// ParseFloat function
// Return first float and error if can't parse or empty string
func (valuesArr ValuesArr) ParseFloat(key string) (float32, error) {
	str := valuesArr[key]
	if len(str) == 0 {
		return 0, errors.New("no this key in map")
	}
	strFloat := str[0]
	i, err := strconv.ParseFloat(strFloat, 32)
	return float32(i), err
}

// ParseBool function
// Return first bool and error if can't parse or empty string
func (valuesArr ValuesArr) ParseBool(key string) (bool, error) {
	str := valuesArr[key]
	if len(str) == 0 {
		return false, errors.New("no this key in map")
	}
	strBool := str[0]
	i, err := strconv.ParseBool(strBool)
	return i, err
}

// ParseString function
// Return first string and error if empty string
func (valuesArr ValuesArr) ParseString(key string) (string, error) {
	str := valuesArr[key]
	if len(str) == 0 {
		return "", errors.New("no this key in map")
	}
	strStr := str[0]
	return strStr, nil
}

// ParseInts function
// Return []int and error if can't parse or empty string
func (valuesArr ValuesArr) ParseInts(key string) ([]int, error) {
	str := valuesArr[key]
	var mas []int
	for _, s := range str {
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		mas = append(mas, int(i))
	}
	return mas, nil
}

// ParseUints function
// Return []uint and error if can't parse or empty string
func (valuesArr ValuesArr) ParseUints(key string) ([]uint, error) {
	str := valuesArr[key]
	var mas []uint
	for _, s := range str {
		i, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}
		mas = append(mas, uint(i))
	}
	return mas, nil
}

// ParseFloats function
// Return []float and error if can't parse or empty string
func (valuesArr ValuesArr) ParseFloats(key string) ([]float32, error) {
	str := valuesArr[key]
	var mas []float32
	for _, s := range str {
		i, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		mas = append(mas, float32(i))
	}
	return mas, nil
}

// ParseBools function
// Return []bool and error if can't parse or empty string
func (valuesArr ValuesArr) ParseBools(key string) ([]bool, error) {
	str := valuesArr[key]
	var mas []bool
	for _, s := range str {
		i, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		mas = append(mas, i)
	}
	return mas, nil
}

// ParseStrings function
// Return []string and error always nil
func (valuesArr ValuesArr) ParseStrings(key string) ([]string, error) {
	str := valuesArr[key]
	return str, nil
}

// ParseFile function
// Return file in []byte, filename and error
func (filesArr FilesArr) ParseFile(key string) ([]byte, string, error) {
	fileHeaderStr := filesArr[key]
	if len(fileHeaderStr) == 0 {
		return []byte{}, "", errors.New("no this key in map")
	}
	fileHeader := fileHeaderStr[0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	bytes := make([]byte, MaxBytes)
	_, err = file.Read(bytes)
	return bytes, fileHeader.Filename, err
}

func getContext(writer http.ResponseWriter, request *http.Request) *Context {
	context := Context{
		Writer:  writer,
		Request: request,
	}
	context.Storage = make(map[string]interface{})
	return &context
}
