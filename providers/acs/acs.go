package acs

import (
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
	"time"

	"waba-alertmanager/notify"
)

type Config struct {
	URL			string	`yaml:"url"`
	ApiToken	string	`yaml:"api_token"`
	Sender		SenderComponent		`yaml:"sender"`
	Template	TemplateComponent	`yaml:"template"`
}

type SenderComponent struct {
	ID	string 	`yaml:"id"`
}

type TemplateComponent struct {
	Name	string 	`yaml:"name"`
	Lang	string 	`yaml:"lang"`
}

const RequestTimeout = time.Second * 20

var _ (notify.Provider) = (*Acs)(nil)

type Acs struct {
	Config
}

func NewAcs(config Config) *Acs {
	Acs := &Acs{config}
	return Acs
}

func (c *Acs) Send(message notify.Message) error {
	for _, recipient := range message.To {
		payload := map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"messaging_product": "whatsapp",
					"recipient_type":    "individual",
					"type":              "template",
					"template": map[string]interface{}{
						"name": c.Template.Name,
						"language": map[string]interface{}{
							"code": c.Template.Lang,
						},
						"components": []interface{}{
							map[string]interface{}{
								"type": "body",
								"parameters": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": message.Text,
									},
								},
							},
						},
					},
				},
			},
			"senderId": c.Sender.ID,
			"client":   []string{recipient},
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		
		request, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		request.Header.Set("Api-Token", c.ApiToken)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("User-Agent", "AlertManager")
		
		httpClient := &http.Client{}
		httpClient.Timeout = RequestTimeout

		response, err := httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("Failed sending message. statusCode: %d", response.StatusCode)
		}

		return nil
	}

	return nil
}