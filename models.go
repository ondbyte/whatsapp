package whatsapp

type BusinessInfo struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	TimezoneID               string `json:"timezone_id"`
	MessageTemplateNamespace string `json:"message_template_namespace"`
}

// Valid implements util.Body.
func (bi *BusinessInfo) Valid() bool {
	return bi.ID != ""
}
