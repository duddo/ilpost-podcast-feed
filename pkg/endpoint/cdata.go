package endpoint

import "encoding/xml"

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		XMLName xml.Name
		Content string `xml:",cdata"`
	}{
		XMLName: start.Name,
		Content: string(c),
	}, start)
}
