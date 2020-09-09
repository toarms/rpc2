// when add new protobuf object, Need to add some related code to "Encode" and "Decode"

package proto

import (
		"log"
		"github.com/golang/protobuf/proto"
	   )

const (
		NoneStr	= "none"
	  )


// encode
func Encode(name string, x interface{}) []byte {
	var p Pb
	switch x.(type) {
		case *Hello:
			out, err := proto.Marshal(x.(*Hello))
			if err != nil {
				log.Fatalln("marshal hello fail: ", err)
			}
			p.Name	= name
			p.Data	= out
			p.Length = int32(len(p.Data))
		//TODO
		//case "NewPBObject":
			// Add new protobuf object here ...
		default:
	}

	out, err := proto.Marshal(&p)
	if err != nil {
		log.Fatalln("Failted to Marshal pb")
	}
	return out
}

// decode
func Decode(in []byte) (string, interface{}) {
	var p Pb
	if err := proto.Unmarshal(in, &p); err != nil {
		log.Fatalln("unmarshal failed: ", err)
	}

	switch p.Name {
		case "world":
			var h Hello
			if err := proto.Unmarshal(p.Data, &h); err != nil {
				break
			}
			return p.Name, &h
		//TODO
		case "NewPBName":
			// Add new protobuf object here ...
		default:
	}

	return NoneStr, nil
}

