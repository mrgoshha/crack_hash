package hash

//<?xml version="1.0" encoding="utf-8"?>
//<!-- Created with Liquid Technologies Online Tools 1.0 (https://www.liquid-technologies.com) -->
//<CrackHashManagerRequest xmlns="http://ccfit.nsu.ru/schema/crack-hash-request" xsi:schemaLocation="http://ccfit.nsu.ru/schema/crack-hash-request schema.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
//	<RequestId>string</RequestId>
//	<PartNumber>-870</PartNumber>
//	<PartCount>2566</PartCount>
//	<Hash>string</Hash>
//	<MaxLength>-4590</MaxLength>
//	<Alphabet>
//		<symbols>string</symbols>
//		<symbols>string</symbols>
//	</Alphabet>
//</CrackHashManagerRequest>

type CrackHashManagerRequest struct {
	RequestId  string   `xml:"RequestId"`
	PartNumber int      `xml:"PartNumber"`
	PartCount  int      `xml:"PartCount"`
	Hash       string   `xml:"Hash"`
	MaxLength  int      `xml:"MaxLength"`
	Alphabet   *Symbols `xml:"Alphabet"`
}

type Symbols struct {
	Symbols []string `xml:"symbols>string"`
}

func (s *Symbols) String() string {
	var str string
	for _, symbol := range s.Symbols {
		str += symbol
	}
	return str
}

//<?xml version="1.0" encoding="utf-8"?>
//<!-- Created with Liquid Technologies Online Tools 1.0 (https://www.liquid-technologies.com) -->
//<CrackHashWorkerResponse xmlns="http://ccfit.nsu.ru/schema/crack-hash-response" xsi:schemaLocation="http://ccfit.nsu.ru/schema/crack-hash-response schema.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
//	<RequestId xmlns="">string</RequestId>
//	<PartNumber xmlns="">-870</PartNumber>
//	<Answers xmlns="">
//		<words>string</words>
//	</Answers>
//</CrackHashWorkerResponse>

type CrackHashWorkerResponse struct {
	RequestId  string `xml:"RequestId"`
	PartNumber int    `xml:"PartNumber"`
	Answers    *Words `xml:"Alphabet"`
}

type Words struct {
	Words []string `xml:"words>string"`
}
