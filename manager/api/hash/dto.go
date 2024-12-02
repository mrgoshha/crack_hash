package hash

type CrackHashRequest struct {
	Hash      string `json:"hash"`
	MaxLength int    `json:"maxLength"`
}

type ResponseID struct {
	RequestID string `json:"request_id"`
}

type ResponseResult struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

//<?xml version="1.0" encoding="utf-8"?>
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

//<?xml version="1.0" encoding="utf-8"?>
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
