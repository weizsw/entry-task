package rpc

import (
	"log"

	"github.com/weizsw/entry-task/pb"
)

func SendRequestForResponse(session *Session, messageProto *pb.Msg) (*pb.Msg, error) {
	data, err := messageProto.Encode()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = session.Write(data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	res, err := session.Read()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	ret := &pb.Msg{}
	err = ret.Decode(res)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return ret, nil
}

func SendRequest(session *Session, msg *pb.Msg) error {
	res, err := msg.Encode()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = session.Write(res)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
