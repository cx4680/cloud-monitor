package sys_rocketmq

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"go.uber.org/atomic"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	ato "sync/atomic"
	"time"
)

type LanguageCode byte

const (
	_Java    = LanguageCode(0)
	_Go      = LanguageCode(9)
	_Unknown = LanguageCode(127)
)

func (lc LanguageCode) MarshalJSON() ([]byte, error) {
	return []byte(`"GO"`), nil
}

func (lc *LanguageCode) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "JAVA":
		*lc = _Java
	case "GO", `"GO"`:
		*lc = _Go
	default:
		*lc = _Unknown
	}
	return nil
}

type RemotingCommand struct {
	Code      int16             `json:"code"`
	Language  LanguageCode      `json:"language"`
	Version   int16             `json:"version"`
	Opaque    int32             `json:"opaque"`
	Flag      int32             `json:"flag"`
	Remark    string            `json:"remark"`
	ExtFields map[string]string `json:"extFields"`
	Body      []byte            `json:"-"`
}

const (
	// 0, REQUEST_COMMAND
	RPCType = 0
	// 1, RPC
	RPCOneWay = 1
	//ResponseType for response
	ResponseType = 1
	_Flag        = 0
	_Version     = 317
)

func (command *RemotingCommand) isResponseType() bool {
	return command.Flag&(ResponseType) == ResponseType
}

func (command *RemotingCommand) WriteTo(w io.Writer) error {
	var (
		header []byte
		err    error
	)

	switch codecType {
	case JsonCodecs:
		header, err = jsonSerializer.encodeHeader(command)
	case RocketMQCodecs:
		header, err = rocketMqSerializer.encodeHeader(command)
	}

	if err != nil {
		return err
	}

	frameSize := 4 + len(header) + len(command.Body)
	err = binary.Write(w, binary.BigEndian, int32(frameSize))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, markProtocolType(int32(len(header))))
	if err != nil {
		return err
	}

	_, err = w.Write(header)
	if err != nil {
		return err
	}

	_, err = w.Write(command.Body)
	return err
}

func markProtocolType(source int32) []byte {
	result := make([]byte, 4)
	result[0] = codecType
	result[1] = byte((source >> 16) & 0xFF)
	result[2] = byte((source >> 8) & 0xFF)
	result[3] = byte(source & 0xFF)
	return result
}

type ClientRequestFunc func(*RemotingCommand, net.Addr) *RemotingCommand

type TcpOption struct {
	// TODO
}

type Client struct {
	responseTable    sync.Map
	connectionTable  sync.Map
	option           TcpOption
	processors       map[int16]ClientRequestFunc
	connectionLocker sync.Mutex
	interceptor      primitive.Interceptor
}

type tcpConnWrapper struct {
	net.Conn
	sync.Mutex
	closed atomic.Bool
}

func (wrapper *tcpConnWrapper) destroy() error {
	wrapper.closed.Swap(true)
	return wrapper.Conn.Close()
}

func (wrapper *tcpConnWrapper) isClosed(err error) bool {
	if !wrapper.closed.Load() {
		return false
	}

	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	return opErr.Err.Error() == "use of closed network connection"
}

func initConn(ctx context.Context, addr string) (*tcpConnWrapper, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpConnWrapper{
		Conn: conn,
	}, nil
}

func (c *Client) connect(ctx context.Context, addr string) (*tcpConnWrapper, error) {
	//it needs additional locker.
	c.connectionLocker.Lock()
	defer c.connectionLocker.Unlock()
	tcpConn, err := initConn(ctx, addr)
	if err != nil {
		return nil, err
	}
	go primitive.WithRecover(func() {
		c.receiveResponse(tcpConn)
	})
	return tcpConn, nil
}

type TopicRouteData struct {
	OrderTopicConf string
	BrokerDataList []*BrokerData `json:"brokerDatas"`
}
type BrokerData struct {
	Cluster             string           `json:"cluster"`
	BrokerName          string           `json:"brokerName"`
	BrokerAddresses     map[int64]string `json:"brokerAddrs"`
	brokerAddressesLock sync.RWMutex
}

var opaque int32

func (c *Client) GetBrokerDataList(nameServerAddress string) (*TopicRouteData, error) {
	var defaultTopic = "TBW102"
	var (
		response *RemotingCommand
		err      error
	)

	var reqMap = make(map[string]string)
	reqMap["topic"] = defaultTopic
	cmd := &RemotingCommand{
		Code:      int16(105),
		Version:   _Version,
		Opaque:    ato.AddInt32(&opaque, 1),
		Body:      nil,
		Language:  _Go,
		ExtFields: reqMap,
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	for _, addr := range strings.Split(nameServerAddress, ";") {
		response, err = c.InvokeSync(ctx, addr, cmd)
		if err == nil {
			break
		}
	}
	if err != nil {
		rlog.Error("get broker list error", map[string]interface{}{
			rlog.LogKeyBroker: err,
		})
	}
	switch response.Code {
	case int16(0):
		if response.Body == nil {
			return nil, primitive.NewMQClientErr(response.Code, response.Remark)
		}
		routeData := &TopicRouteData{}

		err = routeData.decode(string(response.Body))
		if err != nil {
			rlog.Warning("decode TopicRouteData error: %s", map[string]interface{}{
				rlog.LogKeyUnderlayError: err,
				"topic":                  defaultTopic,
			})
			return nil, err
		}
		return routeData, nil
	case int16(17):
		return nil, errors.New("topic not exist")
	default:
		return nil, primitive.NewMQClientErr(response.Code, response.Remark)
	}

}

func (routeData *TopicRouteData) decode(data string) error {
	res := gjson.Parse(data)

	bds := res.Get("brokerDatas").Array()
	routeData.BrokerDataList = make([]*BrokerData, len(bds))
	for idx, v := range bds {
		bd := &BrokerData{
			BrokerName:      v.Get("brokerName").String(),
			Cluster:         v.Get("cluster").String(),
			BrokerAddresses: make(map[int64]string, 0),
		}
		addrs := v.Get("brokerAddrs").String()
		strs := strings.Split(addrs[1:len(addrs)-1], ",")
		if strs != nil {
			for _, str := range strs {
				i := strings.Index(str, ":")
				if i < 0 {
					continue
				}
				id, _ := strconv.ParseInt(str[0:i], 10, 64)
				bd.BrokerAddresses[id] = strings.Replace(str[i+1:], "\"", "", -1)
			}
		}
		routeData.BrokerDataList[idx] = bd
	}
	return nil
}

func (c *Client) InvokeSync(ctx context.Context, addr string, request *RemotingCommand) (*RemotingCommand, error) {
	resp := NewResponseFuture(ctx, request.Opaque, nil)
	c.responseTable.Store(resp.Opaque, resp)
	defer c.responseTable.Delete(request.Opaque)

	conn, err := c.connect(ctx, addr)
	if err != nil {
		return nil, err
	}

	err = c.doRequest(conn, request)
	if err != nil {
		return nil, err
	}
	return resp.waitResponse()
}

func (c *Client) doRequest(conn *tcpConnWrapper, request *RemotingCommand) error {
	conn.Lock()
	defer conn.Unlock()
	err := request.WriteTo(conn)
	if err != nil {
		c.closeConnection(conn)
		return err
	}
	return nil
}

func (c *Client) closeConnection(toCloseConn *tcpConnWrapper) {
	c.connectionTable.Range(func(key, value interface{}) bool {
		if value == toCloseConn {
			c.connectionTable.Delete(key)
			return false
		} else {
			return true
		}
	})
}

func (c *Client) ShutDown() {
	c.responseTable.Range(func(key, value interface{}) bool {
		c.responseTable.Delete(key)
		return true
	})
	c.connectionTable.Range(func(key, value interface{}) bool {
		conn := value.(*tcpConnWrapper)
		err := conn.destroy()
		if err != nil {
			rlog.Warning("close remoting conn error", map[string]interface{}{
				"remote":                 conn.RemoteAddr(),
				rlog.LogKeyUnderlayError: err,
			})
		}
		return true
	})
}

func (c *Client) receiveResponse(r *tcpConnWrapper) {
	var err error
	header := primitive.GetHeader()
	defer primitive.BackHeader(header)
	for {
		if err != nil {
			// conn has been closed actively
			if r.isClosed(err) {
				return
			}
			if err != io.EOF {
				rlog.Error("conn error, close connection", map[string]interface{}{
					rlog.LogKeyUnderlayError: err,
				})
			}
			c.closeConnection(r)
			r.destroy()
			break
		}

		_, err = io.ReadFull(r, header)
		if err != nil {
			continue
		}

		var length int32
		err = binary.Read(bytes.NewReader(header), binary.BigEndian, &length)
		if err != nil {
			continue
		}

		buf := make([]byte, length)

		_, err = io.ReadFull(r, buf)
		if err != nil {
			continue
		}

		cmd, err := decode(buf)
		if err != nil {
			rlog.Error("decode RemotingCommand error", map[string]interface{}{
				rlog.LogKeyUnderlayError: err,
			})
			continue
		}
		c.processCMD(cmd, r)
	}
}

func (c *Client) processCMD(cmd *RemotingCommand, r *tcpConnWrapper) {
	if cmd.isResponseType() {
		resp, exist := c.responseTable.Load(cmd.Opaque)
		if exist {
			c.responseTable.Delete(cmd.Opaque)
			responseFuture := resp.(*ResponseFuture)
			go primitive.WithRecover(func() {
				responseFuture.ResponseCommand = cmd
				if responseFuture.Done != nil {
					close(responseFuture.Done)
				}
			})
		}
	}
}

type ResponseFuture struct {
	ResponseCommand *RemotingCommand
	Err             error
	Opaque          int32
	callback        func(*ResponseFuture)
	Done            chan bool
	callbackOnce    sync.Once
	ctx             context.Context
}

func NewResponseFuture(ctx context.Context, opaque int32, callback func(*ResponseFuture)) *ResponseFuture {
	return &ResponseFuture{
		Opaque:   opaque,
		Done:     make(chan bool),
		callback: callback,
		ctx:      ctx,
	}
}

func (r *ResponseFuture) waitResponse() (*RemotingCommand, error) {
	var (
		cmd *RemotingCommand
		err error
	)
	select {
	case <-r.Done:
		cmd, err = r.ResponseCommand, r.Err
	case <-r.ctx.Done():
		err = errors.New("request timeout")
		r.Err = err
	}
	return cmd, err
}

const (
	JsonCodecs     = byte(0)
	RocketMQCodecs = byte(1)
)

type jsonCodec struct{}

func (c *jsonCodec) encodeHeader(command *RemotingCommand) ([]byte, error) {
	buf, err := jsoniter.Marshal(command)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *jsonCodec) decodeHeader(header []byte) (*RemotingCommand, error) {
	command := &RemotingCommand{}
	command.ExtFields = make(map[string]string)
	command.Body = make([]byte, 0)
	err := jsoniter.Unmarshal(header, command)
	if err != nil {
		return nil, err
	}
	return command, nil
}

type rmqCodec struct{}

const (
	// header + body length
	headerFixedLength = 21
)

// encodeHeader
func (c *rmqCodec) encodeHeader(command *RemotingCommand) ([]byte, error) {
	extBytes, err := c.encodeMaps(command.ExtFields)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, headerFixedLength+len(command.Remark)+len(extBytes)))
	buf.Reset()

	// request code, length is 2 bytes
	err = binary.Write(buf, binary.BigEndian, int16(command.Code))
	if err != nil {
		return nil, err
	}

	// language flag, length is 1 byte
	err = binary.Write(buf, binary.BigEndian, _Go)
	if err != nil {
		return nil, err
	}

	// version flag, length is 2 bytes
	err = binary.Write(buf, binary.BigEndian, int16(command.Version))
	if err != nil {
		return nil, err
	}

	// opaque flag, opaque is request identifier, length is 4 bytes
	err = binary.Write(buf, binary.BigEndian, command.Opaque)
	if err != nil {
		return nil, err
	}

	// request flag, length is 4 bytes
	err = binary.Write(buf, binary.BigEndian, command.Flag)
	if err != nil {
		return nil, err
	}

	// remark length flag, length is 4 bytes
	err = binary.Write(buf, binary.BigEndian, int32(len(command.Remark)))
	if err != nil {
		return nil, err
	}

	// write remark, len(command.Remark) bytes
	if len(command.Remark) > 0 {
		err = binary.Write(buf, binary.BigEndian, []byte(command.Remark))
		if err != nil {
			return nil, err
		}
	}

	err = binary.Write(buf, binary.BigEndian, int32(len(extBytes)))
	if err != nil {
		return nil, err
	}

	if len(extBytes) > 0 {
		err = binary.Write(buf, binary.BigEndian, extBytes)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (c *rmqCodec) encodeMaps(maps map[string]string) ([]byte, error) {
	if maps == nil || len(maps) == 0 {
		return []byte{}, nil
	}
	extFieldsBuf := bytes.NewBuffer([]byte{})
	var err error
	for key, value := range maps {
		err = binary.Write(extFieldsBuf, binary.BigEndian, int16(len(key)))
		if err != nil {
			return nil, err
		}
		err = binary.Write(extFieldsBuf, binary.BigEndian, []byte(key))
		if err != nil {
			return nil, err
		}

		err = binary.Write(extFieldsBuf, binary.BigEndian, int32(len(value)))
		if err != nil {
			return nil, err
		}
		err = binary.Write(extFieldsBuf, binary.BigEndian, []byte(value))
		if err != nil {
			return nil, err
		}
	}
	return extFieldsBuf.Bytes(), nil
}

func (c *rmqCodec) decodeHeader(data []byte) (*RemotingCommand, error) {
	var err error
	command := &RemotingCommand{}
	buf := bytes.NewBuffer(data)
	var code int16
	err = binary.Read(buf, binary.BigEndian, &code)
	if err != nil {
		return nil, err
	}
	command.Code = code

	var (
		languageCode byte
		remarkLen    int32
		extFieldsLen int32
	)
	err = binary.Read(buf, binary.BigEndian, &languageCode)
	if err != nil {
		return nil, err
	}
	command.Language = LanguageCode(languageCode)

	var version int16
	err = binary.Read(buf, binary.BigEndian, &version)
	if err != nil {
		return nil, err
	}
	command.Version = version

	// int opaque
	err = binary.Read(buf, binary.BigEndian, &command.Opaque)
	if err != nil {
		return nil, err
	}

	// int flag
	err = binary.Read(buf, binary.BigEndian, &command.Flag)
	if err != nil {
		return nil, err
	}

	// String remark
	err = binary.Read(buf, binary.BigEndian, &remarkLen)
	if err != nil {
		return nil, err
	}

	if remarkLen > 0 {
		var remarkData = make([]byte, remarkLen)
		err = binary.Read(buf, binary.BigEndian, &remarkData)
		if err != nil {
			return nil, err
		}
		command.Remark = string(remarkData)
	}

	err = binary.Read(buf, binary.BigEndian, &extFieldsLen)
	if err != nil {
		return nil, err
	}

	if extFieldsLen > 0 {
		extFieldsData := make([]byte, extFieldsLen)
		err = binary.Read(buf, binary.BigEndian, &extFieldsData)
		if err != nil {
			return nil, err
		}

		command.ExtFields = make(map[string]string)
		buf := bytes.NewBuffer(extFieldsData)
		var (
			kLen int16
			vLen int32
		)
		for buf.Len() > 0 {
			err = binary.Read(buf, binary.BigEndian, &kLen)
			if err != nil {
				return nil, err
			}

			key, err := getExtFieldsData(buf, int32(kLen))
			if err != nil {
				return nil, err
			}

			err = binary.Read(buf, binary.BigEndian, &vLen)
			if err != nil {
				return nil, err
			}

			value, err := getExtFieldsData(buf, vLen)
			if err != nil {
				return nil, err
			}
			command.ExtFields[key] = value
		}
	}

	return command, nil
}

func getExtFieldsData(buff *bytes.Buffer, length int32) (string, error) {
	var data = make([]byte, length)
	err := binary.Read(buff, binary.BigEndian, &data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

var (
	jsonSerializer     = &jsonCodec{}
	rocketMqSerializer = &rmqCodec{}
	codecType          byte
)

func decode(data []byte) (*RemotingCommand, error) {
	buf := bytes.NewBuffer(data)
	length := int32(len(data))
	var oriHeaderLen int32
	err := binary.Read(buf, binary.BigEndian, &oriHeaderLen)
	if err != nil {
		return nil, err
	}

	headerLength := oriHeaderLen & 0xFFFFFF
	headerData := make([]byte, headerLength)
	err = binary.Read(buf, binary.BigEndian, &headerData)
	if err != nil {
		return nil, err
	}

	var command *RemotingCommand
	switch codeType := byte((oriHeaderLen >> 24) & 0xFF); codeType {
	case JsonCodecs:
		command, err = jsonSerializer.decodeHeader(headerData)
	case RocketMQCodecs:
		command, err = rocketMqSerializer.decodeHeader(headerData)
	default:
		err = fmt.Errorf("unknown codec type: %d", codeType)
	}
	if err != nil {
		return nil, err
	}

	bodyLength := length - 4 - headerLength
	if bodyLength > 0 {
		bodyData := make([]byte, bodyLength)
		err = binary.Read(buf, binary.BigEndian, &bodyData)
		if err != nil {
			return nil, err
		}
		command.Body = bodyData
	}
	return command, nil
}
