package main

import "errors"
import "github.com/go-kit/kit/endpoint"
import "context"
import httptransport "github.com/go-kit/kit/transport/http"
import "net/http"
import "time"
import "fmt"
import "encoding/json"
import "bytes"
import "io/ioutil"

type KeyValueService interface {
	Set(string, string) (string, error)
	Get(string) (string, error)
}

var ERROR_NOT_KEY = errors.New("can find this key.")

type MapService struct {
	datas map[string]string
}

func NewMapService() *MapService {
	rst := MapService{
		datas: make(map[string]string),
	}
	return &rst
}

func (service *MapService) Set(key string, value string) (string, error) {
	service.datas[key] = value
	return key, nil
}

func (service *MapService) Get(key string) (string, error) {
	value, ok := service.datas[key]
	if ok {
		return value, nil
	} else {
		return value, ERROR_NOT_KEY
	}

}

type KeyValueStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func makeSetEndpoint(service KeyValueService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(KeyValueStruct)
		_, err := service.Set(req.Key, req.Value)
		return req, err
	}
}

func makeGetEndpoint(service KeyValueService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(KeyValueStruct)
		return service.Get(req.Key)
	}
}

func decodeStruct(_ context.Context, req *http.Request) (interface{}, error) {
	var request KeyValueStruct
	err := json.NewDecoder(req.Body).Decode(&request)
	return request, err
}

func encodeStruct(_ context.Context, writer http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(writer).Encode(resp)
}

func RunService() {
	svc := NewMapService()

	setHandler := httptransport.NewServer(
		makeSetEndpoint(svc),
		decodeStruct,
		encodeStruct,
	)

	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeStruct,
		encodeStruct,
	)

	http.Handle("/set", setHandler)
	http.Handle("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

func ReqTest(key string, value string, url string) (string, error) {

	contentType := "application/json;charset=utf-8"
	req := KeyValueStruct{
		Key:   key,
		Value: value,
	}

	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return key, err
	}

	body := bytes.NewBuffer(data)

	resp, err := http.Post(url, contentType, body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return key, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return key, err
	}

	return string(content), nil
}

func ReqSet(key string, value string) error {

	_, err := ReqTest(key, value, "http://127.0.0.1:8080/set")
	return err
}

func ReqGet(key string) (string, error) {
	return ReqTest(key, "", "http://127.0.0.1:8080/get")
}

func RunReq() {
	time.Sleep(time.Duration(1) * time.Second)

	fmt.Println("set key1:123")
	ReqSet("key1", "123")

	fmt.Println("set key2:456")
	ReqSet("key2", "456")

	value, _ := ReqGet("key1")
	fmt.Println("get key1:", value)

	value, _ = ReqGet("key2")
	fmt.Println("get key2:", value)
}

func main() {
	go RunReq()

	RunService()

}
