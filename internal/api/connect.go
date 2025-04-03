package api

import (
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Connect(req *desc.ConnectRequest, stream desc.ChatV1_ConnectServer) error {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetId()]
	i.mxChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat channel not found: %s", req.GetId())
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetId()]; !okChat {
		i.chats[req.GetId()] = &Chat{
			streams: make(map[string]desc.ChatV1_ConnectServer),
		}
	}
	i.mxChat.Unlock()

	i.chats[req.GetId()].m.Lock()
	i.chats[req.GetId()].streams[req.GetUsername()] = stream
	i.chats[req.GetId()].m.Unlock()

	for {
		select {
		case msg, okChan := <-chatChan:
			if !okChan {
				return nil
			}

			for user, st := range i.chats[req.GetId()].streams {
				if user == req.GetUsername() {
					continue
				}
				if err := st.Send(msg); err != nil {
					return err
				}
			}

		case <-stream.Context().Done():
			i.chats[req.GetId()].m.Lock()
			delete(i.chats[req.GetId()].streams, req.GetUsername())
			i.chats[req.GetId()].m.Unlock()
			return nil
		}
	}

}
