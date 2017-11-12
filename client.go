package goteleport

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func (t *Teleporter) clientListenForOutboundMessageBuffer(){
	for {
		v, ok := <- t.out
		if !ok {
			continue
		}

		b, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}

		m := Message{
			MType:DATA,
			Payload:b,
		}
		url := fmt.Sprintf("http://%s", t.master)
		_, _, errs := gorequest.New().Post(url).SendStruct(m).End()
		if errs != nil {
			log.Println(errs)
			fmt.Println(errs)
		}
	}
}

func (t *Teleporter) clientListenForInboundMessageBuffer(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return
		}

		var m Message
		err = json.Unmarshal(b, &m)
		if err != nil{
			fmt.Println(err)
			return
		}

		switch m.MType{
		case DATA:
			t.in <- m.Payload
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", t.port), nil)
}

func (t *Teleporter) clientConnectToMaster(){

	m := Message{
		MType: PING,
		Payload:[]byte(strconv.Itoa(t.port)),
	}

	url := fmt.Sprintf("http://%s", t.master)
	_,_, errs := gorequest.New().Post(url).SendStruct(m).End()
	if errs != nil {
		log.Fatal(errs)
	}
}
