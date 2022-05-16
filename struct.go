package main

type Config struct {
	Server struct {
		Address string `json:"address,omitempty"`
		Port    string `json:"port,omitempty"`
	} `json:"server,omitempty"`

	Records struct {
		Domains []domain `json:"domain,omitempty"`
	}
}

type domain struct {
	Record string `json:"record,omitempty"`
	TTL    string `json:"ttl,omitempty"`
	A      string `json:"A,omitempty"`
	AAAA   string `json:"AAAA,omitempty"`
	CNAME  string `json:"CNAME,omitempty"`
	MX     mx     `json:"MX,omitempty"`
	NS     string `json:"NS,omitempty"`
	SOA    string `json:"SOA,omitempty"`
	SRV    srv    `json:"SRV,omitempty"`
	TXT    string `json:"TXT,omitempty"`
}

type srv struct {
	Priority string `json:"priority,omitempty"`
	Weight   string `json:"weight,omitempty"`
	Port     string `json:"port,omitempty"`
	Target   string `json:"target,omitempty"`
}

type mx struct {
	Destination string `json:"destination,omitempty"`
	Priority    string `json:"priority,omitempty"`
}
