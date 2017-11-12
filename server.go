package goteleport

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"github.com/tspn/sliceutils"
	"net/http"
	"io/ioutil"
	"log"
)

func (t *Teleporter) serverListenForInboundMessageBuffer(){
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		host := getHost(r.RemoteAddr)

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			return
		}

		var m Message
		err = json.Unmarshal(b, &m)
		if err != nil{
			log.Println(err)
			return
		}

		switch m.MType {
		case DATA:
			t.in <- m.Payload
		case PING:
			t.lock.Lock()
			if sliceutils.String(t.client).Filter(func(i interface{}) bool {
				return i.(string) == fmt.Sprintf("http://%s:%s", host, m.Payload)
			}).Len() == 0{
				t.client = append(t.client, fmt.Sprintf("http://%s:%s", host, m.Payload))
				fmt.Println("add", host)
			}
			t.lock.Unlock()
		}
		fmt.Println(t.client)
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", t.port), nil)
}


func (t *Teleporter) sendMessageToEachClient(i interface{}) {
	addrs := []string{}
	for _, addr := range t.client {
		_, _, errs := gorequest.New().Post(addr).SendStruct(i).End()
		if errs != nil {
			continue
		}
		addrs = append(addrs, addr)
	}
	t.client = addrs
}

func (t *Teleporter) serverListenForOutboundMessageBuffer(){
	for{
		v, ok := <- t.out
		if !ok {
			continue
		}

		b, err := json.Marshal(v)
		if err != nil{
			continue
		}

		d := Message{
			MType:   DATA,
			Payload: b,
		}
		t.lock.Lock()
		t.sendMessageToEachClient(d)
		t.lock.Unlock()
	}
}
