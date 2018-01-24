package event

import (
	"github.com/asaskevich/EventBus"
	"it-chain/common"
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"github.com/golang/protobuf/proto"
	"errors"
)

var logger_event_publisher = common.GetLogger("message_publisher.go")

//message를 받으면 관심이 있는 subscriber에게 전달하는 역활을 한다.
type MessageHandler struct{
	bus 		EventBus.Bus
	topicMap    map[string]string
	//signer 있어야함
}

//todo signer를 받아야함
//
//topic의 일치성을 위해 처음에 topic의 list를 받고 이 list에 없는 topic은 등록불가하다.
func NewMessageHandler(messageTypes []string) *MessageHandler{

	topicMap := make(map[string]string)

	for _,messageType := range messageTypes{
		topicMap[messageType] = messageType
	}

	return &MessageHandler{
		bus: EventBus.New(),
		topicMap: topicMap,
	}
}

//subscriber를 등록한다.
func (mh *MessageHandler) AddSubscriber(topic string, subfunc func(interface{}), transactional bool) error{

	t, ok := mh.topicMap[topic]

	if ok {
		err := mh.bus.SubscribeAsync(t, subfunc, transactional)

		if err != nil {
			logger_event_publisher.Error("failed to add subscriber", err.Error())
			return err
		}

		logger_event_publisher.Infoln("new message subscriber added")
		return nil
	}

	logger_event_publisher.Error("failed to add subscriber: invalid topic")
	return errors.New("invaild topic")
}

//invaild message검증 및 전파
func (mh *MessageHandler) ReceivedMessageHandle(message comm.OutterMessage){

	//todo vaild message 검증
	//message
	m := &pb.Message{}

	err := proto.Unmarshal(message.Envelope.Payload,m)

	if err != nil{
		logger_event_publisher.Info("failed to Unmarshal message:", message)
	}

	messageType := m.GetMessageType()

	mt, ok := mh.topicMap[messageType]

	if ok{
		mh.bus.Publish(mt,m)
		logger_event_publisher.Info("message published messageType:", mt)
	}
}

