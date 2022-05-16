package pb

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"

	"google.golang.org/protobuf/proto"
)

const SPLIT_LEN = 4

func NewRPCMsg() *Msg {
	return &Msg{}
}

func (msg *Msg) Send(writer io.Writer) error {
	data, err := proto.Marshal(msg)
	err = binary.Write(writer, binary.BigEndian, []byte(data))
	if err != nil {
		return err
	}
	return nil
}

func Read(r io.Reader) (*Msg, error) {
	msg := NewRPCMsg()
	body, err := ioutil.ReadAll(r)
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (msg *Msg) Encode() ([]byte, error) {
	return proto.Marshal(msg)
}

func (msg *Msg) Decode(b []byte) error {

	err := proto.Unmarshal(b, msg)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
