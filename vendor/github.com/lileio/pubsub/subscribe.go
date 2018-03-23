package pubsub

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	pubsubHandled = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pubsub_server_handled_total",
			Help: "Total number of PubSubs handled on the server, regardless of success or failure.",
		}, []string{"topic", "service", "success"})

	subscriberSize = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pubsub_incoming_bytes",
			Help: "A counter of pubsub subscribers incoming bytes.",
		}, []string{"topic", "service"})
)

func init() {
	prometheus.MustRegister(pubsubHandled)
	prometheus.MustRegister(subscriberSize)
}

func Subscribe(s Subscriber) {
	logrus.Info("lile pubsub: Subscribed to events")
	s.Setup(client)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func (c Client) On(topic, subscriberName string, f Handler, deadline time.Duration, autoAck bool) {
	if c.Provider == nil {
		logrus.Warnf("lile pubsub: can't register handler for topic %s, nil provider", topic)
		return
	}

	if f == nil {
		panic("lile pubsub: handler cannot be nil")
	}

	argType, numArgs := argInfo(f)
	if argType == nil {
		panic("lile pubsub: handler requires at least one argument")
	}

	if numArgs > 3 {
		panic("lile pubsub: handler has too many args")
	}

	handler := reflect.ValueOf(f)
	rawMsgFunctionType := reflect.TypeOf(func(c context.Context, m Msg) error { return nil })
	wantsRaw := (argType == rawMsgFunctionType)

	if wantsRaw {
		c.Provider.Subscribe(topic, subscriberName, f.(MsgHandler), deadline, autoAck)
		return
	}

	cb := func(ctx context.Context, m Msg) error {
		var oV []reflect.Value

		var obj reflect.Value
		if argType.Kind() != reflect.Ptr {
			obj = reflect.New(argType)
		} else {
			obj = reflect.New(argType.Elem())
		}

		err := proto.Unmarshal(m.Data, obj.Interface().(proto.Message))
		if err != nil {
			return errors.Wrap(err, "lile pubsub: could not unmarshal message")
		}

		if argType.Kind() != reflect.Ptr {
			obj = reflect.Indirect(obj)
		}

		switch numArgs {
		case 1:
			oV = []reflect.Value{obj}
		case 2:
			oV = []reflect.Value{reflect.ValueOf(ctx), obj}
		case 3:
			oV = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(m.Metadata), obj}
		}

		returnVal := handler.Call(oV)
		if len(returnVal) == 0 {
			return nil
		}

		errInterface := returnVal[0].Interface()

		pubsubHandled.WithLabelValues(topic, c.ServiceName, strconv.FormatBool(errInterface == nil)).Inc()
		subscriberSize.WithLabelValues(topic, c.ServiceName).Add(float64(len(m.Data)))
		if errInterface == nil {
			return nil
		}

		return errInterface.(error)
	}

	c.Provider.Subscribe(topic, subscriberName, cb, deadline, autoAck)
}

// Dissect the handler's signature
func argInfo(cb Handler) (reflect.Type, int) {
	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		panic("lile pubsub: handler needs to be a func")
	}

	numArgs := cbType.NumIn()
	if numArgs == 0 {
		return nil, numArgs
	}

	return cbType.In(numArgs - 1), numArgs
}
