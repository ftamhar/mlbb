// Code generated by github.com/darmawan01/toldata. DO NOT EDIT.
// package: services
// source: proto/hello/hello.proto

package hello

import (
	"context"
	"errors"
	"github.com/darmawan01/toldata"
	"google.golang.org/protobuf/proto"
	"io"

	nats "github.com/nats-io/nats.go"
)

// Workaround for template problem
func _eof() error {
	return io.EOF
}

type HelloServicesToldataInterface interface {
	ToldataHealthCheck(ctx context.Context, req *toldata.Empty) (*toldata.ToldataHealthCheckInfo, error)

	Greeting(ctx context.Context, req *GreetingRequest) (*GreetingResponse, error)

	Upload(stream HelloServices_UploadToldataServer)

	Upload2(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
}

type HelloServicesToldataClient struct {
	Bus *toldata.Bus
}

type HelloServicesToldataServer struct {
	Bus     *toldata.Bus
	Service HelloServicesToldataInterface
}

func NewHelloServicesToldataClient(bus *toldata.Bus) *HelloServicesToldataClient {
	s := &HelloServicesToldataClient{Bus: bus}
	return s
}

func NewHelloServicesToldataServer(bus *toldata.Bus, service HelloServicesToldataInterface) *HelloServicesToldataServer {
	s := &HelloServicesToldataServer{Bus: bus, Service: service}
	return s
}

func (service *HelloServicesToldataClient) ToldataHealthCheck(ctx context.Context, req *toldata.Empty) (*toldata.ToldataHealthCheckInfo, error) {
	functionName := "services/HelloServices/ToldataHealthCheck"

	reqRaw, err := proto.Marshal(req)
	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	result, err := service.Bus.Connection.RequestWithContext(ctx, functionName, reqRaw)
	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error
		p := &toldata.ToldataHealthCheckInfo{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return nil, errors.New(pErr.ErrorMessage)
		} else {
			return nil, err
		}
	}
}

func (service *HelloServicesToldataClient) Greeting(ctx context.Context, req *GreetingRequest) (*GreetingResponse, error) {
	functionName := "services/HelloServices/Greeting"

	if req == nil {
		return nil, errors.New("empty-request")
	}
	reqRaw, err := proto.Marshal(req)

	result, err := service.Bus.Connection.RequestWithContext(ctx, functionName, reqRaw)
	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error
		p := &GreetingResponse{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return nil, errors.New(pErr.ErrorMessage)
		} else {
			return nil, err
		}
	}
}

type HelloServices_UploadToldataServer interface {
	Receive() (*UploadRequest, error)
	OnData(*UploadRequest) error
	Done(resp *UploadResponse) error

	GetResponse() (*UploadResponse, error)

	TriggerEOF()
	Error(err error)
	OnExit(func())
	Exit()
}

type HelloServices_UploadToldataServerImpl struct {
	request         chan *UploadRequest
	isRequestClosed bool

	response chan *UploadResponse

	cancel chan struct{}
	eof    chan struct{}
	err    chan error
	done   chan struct{}

	isEOF      bool
	isCanceled bool

	streamErr error

	Context context.Context
}

func CreateHelloServices_UploadToldataServerImpl(ctx context.Context) *HelloServices_UploadToldataServerImpl {
	t := &HelloServices_UploadToldataServerImpl{}

	t.Context = ctx
	t.request = make(chan *UploadRequest)
	t.response = make(chan *UploadResponse)
	t.cancel = make(chan struct{})
	t.eof = make(chan struct{})
	t.done = make(chan struct{})
	t.err = make(chan error)
	return t
}

func (impl *HelloServices_UploadToldataServerImpl) Exit() {
	close(impl.done)
}

func (impl *HelloServices_UploadToldataServerImpl) OnExit(fn func()) {
	go func() {
		select {
		case <-impl.done:
			fn()
		}
	}()
}

func (impl *HelloServices_UploadToldataServerImpl) TriggerEOF() {
	if impl.streamErr != nil {
		return
	}
	if impl.isEOF == false {
		close(impl.eof)
		impl.isEOF = true
	}
}

func (impl *HelloServices_UploadToldataServerImpl) Receive() (*UploadRequest, error) {

	if impl.streamErr != nil {
		return nil, impl.streamErr
	}
	if impl.isEOF {
		return nil, io.EOF
	}

	select {
	case data := <-impl.request:
		return data, impl.streamErr
	case <-impl.cancel:
		return nil, impl.streamErr
	case <-impl.eof:
		return nil, io.EOF
	case err := <-impl.err:

		return nil, err

	}
}

func (impl *HelloServices_UploadToldataServerImpl) OnData(req *UploadRequest) error {
	if impl.streamErr != nil {
		return impl.streamErr
	}

	select {
	case err := <-impl.err:
		return err
	case impl.request <- req:
		return nil
	}
}

func (impl *HelloServices_UploadToldataServerImpl) Done(resp *UploadResponse) error {
	if impl.streamErr != nil {
		return impl.streamErr
	}

	select {
	case impl.response <- resp:
		return nil
	case err := <-impl.err:
		return err

	}
}

func (impl *HelloServices_UploadToldataServerImpl) GetResponse() (*UploadResponse, error) {
	if impl.streamErr != nil {
		return nil, impl.streamErr
	}

	select {
	case err := <-impl.err:

		return nil, err

	case <-impl.cancel:

		return nil, errors.New("canceled")

	case response := <-impl.response:
		return response, nil

	}
}

func (impl *HelloServices_UploadToldataServerImpl) Cancel() {
	if impl.isCanceled == false {
		close(impl.cancel)
		impl.isCanceled = true
	}
}

func (impl *HelloServices_UploadToldataServerImpl) Error(err error) {
	impl.err <- err
	impl.streamErr = err
}

type HelloServicesToldataClient_Upload struct {
	Context context.Context
	Service *HelloServicesToldataClient
	ID      string
}

func (client *HelloServicesToldataClient_Upload) Send(req *UploadRequest) error {
	functionName := "services/HelloServices/Upload_Send_" + client.ID
	if req == nil {
		return errors.New("empty-request")
	}
	reqRaw, err := proto.Marshal(req)
	result, err := client.Service.Bus.Connection.RequestWithContext(client.Context, functionName, reqRaw)
	if err != nil {
		return errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error
		return nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return errors.New(pErr.ErrorMessage)
		} else {
			return err
		}
	}
}

func (client *HelloServicesToldataClient_Upload) Done() (*UploadResponse, error) {
	functionName := "services/HelloServices/Upload_Done_" + client.ID

	result, err := client.Service.Bus.Connection.RequestWithContext(client.Context, functionName, nil)

	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error
		p := &UploadResponse{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return nil, errors.New(pErr.ErrorMessage)
		} else {
			return nil, err
		}
	}
}

func (impl *HelloServices_UploadToldataServerImpl) Subscribe(service *HelloServicesToldataServer, id string) error {
	bus := service.Bus
	var sub *nats.Subscription
	var subscriptions []*nats.Subscription
	var err error

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/Upload_Send_"+id, "services/HelloServices", func(m *nats.Msg) {
		var input UploadRequest
		err := proto.Unmarshal(m.Data, &input)
		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		}

		err = impl.OnData(&input)

		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		} else {
			zero := []byte{0}
			bus.Connection.Publish(m.Reply, zero)
		}

	})

	subscriptions = append(subscriptions, sub)

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/Upload_Done_"+id, "services/HelloServices", func(m *nats.Msg) {

		defer impl.Exit()
		impl.TriggerEOF()
		result, err := impl.GetResponse()

		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		}
		raw, err := proto.Marshal(result)
		if err != nil {
			bus.HandleError(m.Reply, err)
		} else {
			zero := []byte{0}
			bus.Connection.Publish(m.Reply, append(zero, raw...))
		}

	})

	subscriptions = append(subscriptions, sub)

	impl.OnExit(func() {
		for i := range subscriptions {
			subscriptions[i].Unsubscribe()
		}
	})

	return err
}

func (service *HelloServicesToldataClient) Upload(ctx context.Context) (*HelloServicesToldataClient_Upload, error) {
	functionName := "services/HelloServices/Upload"

	result, err := service.Bus.Connection.RequestWithContext(ctx, functionName, nil)

	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error

		p := &toldata.StreamInfo{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return nil, err
		}
		return &HelloServicesToldataClient_Upload{
			ID:      p.ID,
			Context: ctx,
			Service: service,
		}, nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return nil, errors.New(pErr.ErrorMessage)
		} else {
			return nil, err
		}
	}
}

func (service *HelloServicesToldataClient) Upload2(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	functionName := "services/HelloServices/Upload2"

	if req == nil {
		return nil, errors.New("empty-request")
	}
	reqRaw, err := proto.Marshal(req)

	result, err := service.Bus.Connection.RequestWithContext(ctx, functionName, reqRaw)
	if err != nil {
		return nil, errors.New(functionName + ":" + err.Error())
	}

	if result.Data[0] == 0 {
		// 0 means no error
		p := &UploadResponse{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return nil, errors.New(pErr.ErrorMessage)
		} else {
			return nil, err
		}
	}
}

func (service *HelloServicesToldataServer) SubscribeHelloServices() (<-chan struct{}, error) {
	bus := service.Bus

	var err error
	var sub *nats.Subscription
	var subscriptions []*nats.Subscription

	done := make(chan struct{})

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/Greeting", "services/HelloServices", func(m *nats.Msg) {
		var input GreetingRequest
		err := proto.Unmarshal(m.Data, &input)
		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		}
		result, err := service.Service.Greeting(bus.Context, &input)

		if m.Reply != "" {
			if err != nil {
				bus.HandleError(m.Reply, err)
			} else {
				raw, err := proto.Marshal(result)
				if err != nil {
					bus.HandleError(m.Reply, err)
				} else {
					zero := []byte{0}
					bus.Connection.Publish(m.Reply, append(zero, raw...))
				}
			}
		}

	})

	subscriptions = append(subscriptions, sub)

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/Upload", "services/HelloServices", func(m *nats.Msg) {
		stream := CreateHelloServices_UploadToldataServerImpl(bus.Context)

		stream.Subscribe(service, m.Reply)

		raw, err := proto.Marshal(&toldata.StreamInfo{
			ID: m.Reply,
		})
		if err != nil {
			bus.HandleError(m.Reply, err)
		} else {
			zero := []byte{0}
			bus.Connection.Publish(m.Reply, append(zero, raw...))
		}

		service.Service.Upload(stream)

	})

	subscriptions = append(subscriptions, sub)

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/Upload2", "services/HelloServices", func(m *nats.Msg) {
		var input UploadRequest
		err := proto.Unmarshal(m.Data, &input)
		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		}
		result, err := service.Service.Upload2(bus.Context, &input)

		if m.Reply != "" {
			if err != nil {
				bus.HandleError(m.Reply, err)
			} else {
				raw, err := proto.Marshal(result)
				if err != nil {
					bus.HandleError(m.Reply, err)
				} else {
					zero := []byte{0}
					bus.Connection.Publish(m.Reply, append(zero, raw...))
				}
			}
		}

	})

	subscriptions = append(subscriptions, sub)

	sub, err = bus.Connection.QueueSubscribe("services/HelloServices/ToldataHealthCheck", "services/HelloServices", func(m *nats.Msg) {
		var input toldata.Empty
		err := proto.Unmarshal(m.Data, &input)
		if err != nil {
			bus.HandleError(m.Reply, err)
			return
		}
		result, err := service.Service.ToldataHealthCheck(bus.Context, &input)

		if m.Reply != "" {
			if err != nil {
				bus.HandleError(m.Reply, err)
			} else {
				raw, err := proto.Marshal(result)
				if err != nil {
					bus.HandleError(m.Reply, err)
				} else {
					zero := []byte{0}
					bus.Connection.Publish(m.Reply, append(zero, raw...))
				}
			}
		}

	})

	subscriptions = append(subscriptions, sub)

	go func() {
		defer close(done)

		select {
		case <-bus.Context.Done():
			for i := range subscriptions {
				subscriptions[i].Unsubscribe()
			}
		}
	}()

	return done, err
}