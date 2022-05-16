package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/miekg/dns"
)

var config Config

func Query(m *dns.Msg) {
	var q dns.Question
	var rr dns.RR
	var err error

	// Query the DNS server
	for _, q = range m.Question {
		record := getRecord(q.Name)
		switch q.Qtype {
		case dns.TypeA:
			if len(record.A) > 0 {
				rr, err = dns.NewRR(fmt.Sprintf("%s A %s", q.Name, record.A))
			}

		case dns.TypeMX:
			if len(record.MX.Destination) > 0 {
				rr, err = dns.NewRR(fmt.Sprintf("%s MX %s %s", q.Name, "10", record.MX))
			}

			// case dns.TypeAAAA:
			// case dns.TypeNS:
			// case dns.TypeCNAME:
			// case dns.TypeTXT:
			// case dns.TypePTR:
			// case dns.TypeSRV:
		}

		fmt.Printf("hostname: %s\n", q.Name)

		if rr != nil {
			if err == nil {
				ttl, _ := strconv.Atoi(record.TTL)
				rr.Header().Ttl = uint32(ttl)
				fmt.Println("Header Type: ", rr.Header().Rrtype)
			} else {
				fmt.Println("Query Error: ", err)
			}

			m.Answer = append(m.Answer, rr)
		}
	}
}

func getRecord(hostname string) domain {
	var record domain

	for i := 0; i < len(config.Records.Domains); i++ {
		if config.Records.Domains[i].Record == hostname {
			record = config.Records.Domains[i]
			break
		}
	}

	return record
}

func HandleRequest(responce dns.ResponseWriter, rm *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(rm)

	switch rm.Opcode {
	case dns.OpcodeQuery:
		Query(m)
	}

	fmt.Println("Remote IP: ", responce.RemoteAddr())

	if m.Answer == nil {
		m.SetRcode(rm, dns.RcodeNameError)
	}

	responce.WriteMsg(m)
}

// Save config to file
func saveConfig(config Config) {
	file, _ := json.MarshalIndent(config, "", "\t")
	_ = ioutil.WriteFile("prvdns.config", file, 0644)
}

func loadConfig(config Config) Config {
	data, _ := ioutil.ReadFile("prvdns.config")
	json.Unmarshal(data, &config)

	return config
}

func record(dnstype uint16, hostname, ip, ttl, param string) domain {
	var record domain

	record.TTL = ttl
	record.Record = hostname

	switch dnstype {
	case dns.TypeA:
		record.A = ip
	case dns.TypeMX:
		record.MX.Destination = ip
		record.MX.Priority = param
	}

	return record
}

func addNewRecord(record domain) {
	var found bool = false

	for i := 0; i < len(config.Records.Domains); i++ {
		if (config.Records.Domains[i].Record == record.Record) && (config.Records.Domains[i].A == record.A) {
			found = true
		}
	}

	if !found {
		config.Records.Domains = append(config.Records.Domains, record)
		fmt.Println("New Record: ", record.Record)
	} else {
		fmt.Println("Duplicate Record: ", record.Record)
	}
}

func main() {
	//Load configuration
	config = loadConfig(config)

	//Catch ctrl+c and save before exiting
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		saveConfig(config)
		os.Exit(1)
	}()

	// Some Default Records
	// config.Records.Domains = append(config.Records.Domains, record(dns.TypeA, "gw.ear.pm", "192.168.178.1", "60", ""))
	// config.Records.Domains = append(config.Records.Domains, record(dns.TypeA, "store.ear.pm", "192.168.178.7", "60", ""))

	//Add Test Record
	addNewRecord(record(dns.TypeA, "atlas.ear.pm.", "192.168.178.30", "60", ""))

	//handle DNS requests
	dns.HandleFunc(".", HandleRequest) //Make patern configurable

	// Start server
	dnsserver := &dns.Server{Addr: ":5053", Net: "udp"}
	err := dnsserver.ListenAndServe()
	defer dnsserver.Shutdown()

	if err != nil {
		fmt.Println("Cant start server: ", err)
	}
}
