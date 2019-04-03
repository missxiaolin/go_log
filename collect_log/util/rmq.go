package util

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
)

//连接结构
type Connect struct {
	Name string `json:Name`
	Addr string `json:Addr`
}

//信道结构
type Channel struct {
	Name string `json:Name`
	Connect string `json:Connect`
	QosCount int `json:QosCount`
	QosSize int `json:QosSize`
}

//交换机绑定结构
type EBind struct {
	Destination string `json:Destination`
	Key string `json:Key`
	NoWait bool `json:NoWait`
}

//交换机结构
type Exchange struct {
	Name string `json:Name`
	Channel string `json:Channel`
	Type string `json:Type`
	Durable bool `json:Durable`
	AutoDeleted bool `json:AutoDeleted`
	Internal bool `json:Internal`
	NoWait bool `json:NoWait`
	Bind []EBind `json:Bind`
	Args map[string]interface{} `json:Args`
}

//队列绑定结构
type QBind struct {
	ExchangeName string `json:ExchangeName`
	Key string `json:Key`
	NoWait bool `json:NoWait`
}

//队列结构
type Queue struct {
	Name string `json:Name`
	Channel string `json:Channel`
	Durable bool `json:Durable`
	AutoDelete bool `json:AutoDelete`
	Exclusive bool `json:Exclusive`
	NoWait bool `json:NoWait`
	Bind []QBind `json:Bind`
	Args map[string]interface{} `json:Args`
}

//发送者配置
type Pusher struct {
	Name string `json:Name`
	Channel string `json:Channel`
	Exchange string `json:Exchange`
	Key string `json:Key`
	Mandtory bool `json:Mandtory`
	Immediate bool `json:Immediate`
	ContentType string `json:ContentType`
	DeliveryMode uint8 `json:DeliveryMode`
}

//接收者配置
type Poper struct {
	Name string `json:Name`
	QName string `json:QName`
	Channel string `json:Channel`
	Consumer string `json:Consumer`
	AutoACK bool `json:AutoACK`
	Exclusive bool `json:Exclusive`
	NoLocal bool `json:NoLocal`
	NoWait bool `json:NoWait`
}

//配置文件结构
type tCfg struct {
	Connects []Connect `json:Connects`
	Channels []Channel `json:Channels`
	Exchanges []Exchange `json:Exchanges`
	Queue []Queue `json:Queue`
	Pusher []Pusher `json:Pusher`
	Poper []Poper `json:Poper`
}

var _Cfg *tCfg = new(tCfg) //配置文件对象
var _ConnectPool map[string]*amqp.Connection = make(map[string]*amqp.Connection) //连接名称:连接对象
var _ChannelPool map[string]*amqp.Channel = make(map[string]*amqp.Channel) //信道名称:信道对象
var _ExchangePool map[string]string = make(map[string]string) //交换机名称:所属信道名称
var _QueuePool map[string]string = make(map[string]string) //队列名称:所属信道名称
var _Pusher map[string]Pusher = make(map[string]Pusher) //Pusher名称:Pusher配置
var _Poper map[string]Poper = make(map[string]Poper) //Poper名称:Poper配置

//读取配置文件
func loadCfg(path string) (err error) {
	var fp *os.File
	if fp, err = os.Open(path); err != nil {
		return err
	}
	var data []byte
	if data, err =  ioutil.ReadAll(fp); err != nil {
		return err
	}
	if err = fp.Close(); err != nil {
		return err
	}
	if err = json.Unmarshal(data, _Cfg); err != nil {
		return err
	}
	return nil
}

func CloseConnect(name string) (err error) {
	if _, ok := _ConnectPool[name]; ok {
		_ConnectPool[name].Close()
	}
	return nil
}

//创建连接
func CreateConnect(v Connect) (err error) {
	var connect *amqp.Connection
	if connect, err = amqp.Dial(v.Addr); err != nil {
		return err
	} else {
		if _, ok := _ConnectPool[v.Name]; !ok {
			_ConnectPool[v.Name] = connect
		} else {
			return errors.New("连接已存在\n")
		}
	}
	return nil
}

//初始化连接
func initConnect() (err error) {
	for _, v := range _Cfg.Connects {
		if err = CreateConnect(v); err != nil {
			return err
		}
	}
	return nil
}

//关闭信道
func CloseChannel(name string) (err error) {
	if _, ok := _ChannelPool[name]; ok {
		_ChannelPool[name].Close()
	}
	return nil
}

//创建信道
func CreateChannel(v Channel) (err error) {
	if _, ok := _ConnectPool[v.Connect]; !ok {
		return errors.New("连接不存在\n")
	}
	var channel *amqp.Channel
	if channel, err = _ConnectPool[v.Connect].Channel(); err != nil {
		return err
	} else {
		if _, ok := _ChannelPool[v.Name]; !ok {
			_ChannelPool[v.Name] = channel
		} else {
			return errors.New("信道已存在\n")
		}

	}
	if err = channel.Qos(v.QosCount, v.QosSize, false); err != nil {
		return nil
	}
	return nil
}

//初始化信道
func initChannel() (err error) {
	for _, v := range _Cfg.Channels {
		if err = CreateChannel(v); err != nil {
			return err
		}
	}
	return nil
}

//删除交换机
func DeleteExchange(name string, ifUnused bool, noWait bool) (err error) {
	if _, ok := _ExchangePool[name]; ok {
		if _, ok := _ChannelPool[_ExchangePool[name]]; ok {
			if err = _ChannelPool[_ExchangePool[name]].ExchangeDelete(
				name, ifUnused, noWait); err != nil {
				return err
			}
		}
		delete(_ExchangePool, name)
	}
	return nil
}

//创建交换机
func CreateExchange(v Exchange) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}
	if err = _ChannelPool[v.Channel].ExchangeDeclare(v.Name, v.Type,
		v.Durable, v.AutoDeleted, v.Internal, v.NoWait, v.Args); err != nil {
		return err
	} else {
		if _, ok := _ExchangePool[v.Name]; !ok {
			_ExchangePool[v.Name] = v.Channel
		} else {
			return errors.New("交换机已存在")
		}
	}
	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].ExchangeBind(b.Destination, b.Key, v.Name, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}

//初始化交换机
func initExchange() (err error) {
	for _, v := range _Cfg.Exchanges {
		if err = CreateExchange(v); err != nil {
			return err
		}
	}
	return nil
}

//删除队列
func DeleteQueue(name string, ifUnused bool, ifEmpty bool, noWait bool) (err error) {
	if _, ok := _QueuePool[name]; ok {
		if _, ok := _ChannelPool[_QueuePool[name]]; ok {
			if _, err = _ChannelPool[_QueuePool[name]].QueueDelete(
				name, ifUnused, ifEmpty, noWait); err != nil {
				return err
			}
		}
		delete(_QueuePool, name)
	}
	return nil
}

//创建队列
func CreateQueue(v Queue) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}

	//处理x-message-ttl的类型，json里面写的是int，go读出来的是double
	if _, ok := v.Args["x-message-ttl"]; ok {
		t := int32(v.Args["x-message-ttl"].(float64))
		delete(v.Args, "x-message-ttl")
		v.Args["x-message-ttl"] = t
	}

	if _, err = _ChannelPool[v.Channel].QueueDeclare(v.Name, v.Durable,
		v.AutoDelete, v.Exclusive, v.NoWait, v.Args); err != nil {
		return err
	} else {
		if _, ok := _QueuePool[v.Name]; !ok {
			_QueuePool[v.Name] = v.Channel
		} else {
			return errors.New("队列已存在")
		}
	}
	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].QueueBind(v.Name, b.Key, b.ExchangeName, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}


//初始化队列
func initQueue() (err error) {
	for _, v := range _Cfg.Queue {
		if err = CreateQueue(v); err != nil {
			return err
		}
	}
	return nil
}

//创建Pusher
func CreatePusher(v Pusher) (err error) {
	if _, ok := _Pusher[v.Name]; !ok {
		_Pusher[v.Name] = v
	} else {
		return errors.New("Pusher已存在")
	}
	return nil
}

//删除Pusher
func DeletePusher(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Pusher
func initPusher() (err error) {
	for _, v := range _Cfg.Pusher {
		if err = CreatePusher(v); err != nil {
			return err
		}
	}
	return err
}

//创建Poper
func CreatePoper(v Poper) (err error) {
	if _, ok := _Poper[v.Name]; !ok {
		_Poper[v.Name] = v
	} else {
		return errors.New("Poper已存在")
	}
	return err
}

//删除Poper
func DeletePoper(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Poper
func initPoper() (err error) {
	for _, v := range _Cfg.Poper {
		if err = CreatePoper(v); err != nil {
			return err
		}
	}
	return nil
}

//关闭
func Fini() (err error) {
	for _, conn := range _ConnectPool {
		for _, ch := range _ChannelPool {
			if err = ch.Close(); err != nil {
				return err
			}
		}
		if err = conn.Close(); err != nil {
			return err
		}
	}
	//清空所有缓存
	_Cfg = new(tCfg) //配置文件对象
	_ConnectPool = make(map[string]*amqp.Connection) //连接名称:连接对象
	_ChannelPool = make(map[string]*amqp.Channel) //信道名称:信道对象
	_ExchangePool = make(map[string]string) //交换机名称:所属信道名称
	_QueuePool = make(map[string]string) //队列名称:所属信道名称
	_Pusher = make(map[string]Pusher) //Pusher名称:Pusher配置
	_Poper = make(map[string]Poper) //Poper名称:Poper配置

	return nil
}

//初始化
func Init(path string) (err error) {
	if err = loadCfg(path); err != nil {
		return err
	}
	if err = initConnect(); err != nil {
		return err
	}
	if err = initChannel(); err != nil {
		return err
	}
	if err = initExchange(); err != nil {
		return err
	}
	if err = initQueue(); err != nil {
		return err
	}
	if err = initPusher(); err != nil {
		return err
	}
	if err = initPoper(); err != nil {
		return err
	}
	return err
}

//向交换机推送一条消息
func Push(name string, key string, msg []byte) (err error) {
	if _, ok := _Pusher[name]; !ok {
		return errors.New("Pusher不存在")
	}

	cfg := _Pusher[name]
	if key != "" {
		cfg.Key = key
	}
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("Channel不存在")
	}
	if err = _ChannelPool[cfg.Channel].Publish(cfg.Exchange, cfg.Key, cfg.Mandtory, cfg.Immediate,
		amqp.Publishing{ContentType:cfg.ContentType,Body:msg}); err != nil {
		return err
	}
	return nil
}

type MSG struct {
	Body []byte
	Tag uint64
	Channel string
	Poper string
}

func (m MSG)Ack(multiple bool) (err error) {
	if _, ok := _ChannelPool[m.Channel]; !ok {
		return errors.New("Ack失败,Channel无效")
	} else {
		_ChannelPool[m.Channel].Ack(m.Tag, multiple)
	}
	return nil
}

//处理消息(顺序处理,如果需要多线程可以在回调函数中做手脚)
func handleMsg(msgs <-chan amqp.Delivery, callback func(MSG), channel string, poperName string) {
	for d := range msgs {
		var msg MSG = MSG {
			Body : d.Body,
			Tag : d.DeliveryTag,
			Channel : channel,
			Poper : poperName,
		}
		callback(msg)
	}
}


//从队列获取消息 -- 推模式
func Pop(name string, callback func(MSG)) (err error) {
	if _, ok := _Poper[name]; !ok {
		return errors.New("Poper不存在")
	}
	cfg := _Poper[name]
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("Channel不存在")
	}
	var msgs <-chan amqp.Delivery
	if msgs, err = _ChannelPool[cfg.Channel].Consume(cfg.QName, cfg.Consumer,
		cfg.AutoACK, cfg.Exclusive, cfg.NoLocal, cfg.NoWait, nil); err != nil {
		return err
	}
	go handleMsg(msgs, callback, cfg.Channel, name)

	return nil
}

