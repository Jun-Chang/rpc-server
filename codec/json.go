package codec

import "github.com/CyberAgent/car-golib/encode"

type CodecJson struct{}

type JsonBinder struct {
	Seq int "`json:seq`"
}

func (c CodecJson) Marshal(v interface{}) ([]byte, error) {
	return encode.EncodeJson(&v)
}

func (c CodecJson) Unmarshal(data []byte, v interface{}) error {
	return encode.DecodeJson(data, &v)
}

func (c CodecJson) String() string {
	return "json"
}
