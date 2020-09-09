package proto

import (
		"log"
		"testing"
		"io/ioutil"
		//"github.com/golang/protobuf/proto"
	   )

func TestExample(t *testing.T) {
	var h Hello

	h.Greeting = "各位N同事"
	h.Gift = "绩效面谈"

	// encode 
	out := Encode("world", &h)

	// write file, then read file
	if err := ioutil.WriteFile("pb.out", out, 0644); err != nil {
		log.Fatalln("Failted to write pb.out")
	}
	in, err := ioutil.ReadFile("pb.out")
	if err != nil {
		log.Fatalln("Read pb.out failed: " , err)
	}

	//decode
	name, hhh := Decode(in)
	log.Printf("\n\t{\n\t\tname:  \"%s\"\n\t\tproto: %s\n\t}\n", name,hhh.(*Hello).String())
}
