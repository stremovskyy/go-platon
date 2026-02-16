/*
 * MIT License
 *
 * Copyright (c) 2026 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package platon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Result string

func (r Result) String() string {
	return string(r)
}

const (
	ResultAccepted Result = "ACCEPTED"
	ResultDeclined Result = "DECLINED"
	ResultError    Result = "ERROR"
)

type Response struct {
	Status        *string       `json:"status,omitempty"`
	Action        *string       `json:"action"`
	Result        *Result       `json:"result"`
	OrderId       *string       `json:"order_id"`
	TransId       *string       `json:"trans_id"`
	TransDate     *string       `json:"trans_date"`
	ResponseData  *ResponseData `json:"response,omitempty"`
	ErrorMessage  string        `json:"error_message"`
	DeclineReason string        `json:"decline_reason"`
}

type ResponseData struct {
	SubmerchantID       *string `json:"submerchant_id,omitempty"`
	SubmerchantIDStatus *string `json:"submerchant_id_status,omitempty"`
	Hash                *string `json:"hash,omitempty"`
}

func (p *Response) PrettyPrint() {
	if p == nil {
		fmt.Println("Error: Response is nil")
		return
	}

	fmt.Println("\nPlaton response:")
	fmt.Println("------------------------------------------------------")
	if p.Status != nil {
		fmt.Printf("status: %s\n", *p.Status)
	}
	if p.Action != nil {
		fmt.Printf("action: %s\n", *p.Action)
	}
	if p.Result != nil {
		fmt.Printf("result: %s\n", p.Result.String())
	}
	if p.OrderId != nil {
		fmt.Printf("order_id: %s\n", *p.OrderId)
	}
	if p.TransId != nil {
		fmt.Printf("trans_id: %s\n", *p.TransId)
	}
	if p.TransDate != nil {
		fmt.Printf("trans_date: %s\n", *p.TransDate)
	}
	if p.ResponseData != nil && p.ResponseData.SubmerchantID != nil {
		fmt.Printf("submerchant_id: %s\n", *p.ResponseData.SubmerchantID)
	}
	if p.ResponseData != nil && p.ResponseData.SubmerchantIDStatus != nil {
		fmt.Printf("submerchant_id_status: %s\n", *p.ResponseData.SubmerchantIDStatus)
	}
	if p.ErrorMessage != "" {
		fmt.Printf("error_message: %s\n", p.ErrorMessage)
	}
	if p.DeclineReason != "" {
		fmt.Printf("decline_reason: %s\n", p.DeclineReason)
	}
	fmt.Println("------------------------------------------------------")
}

func (p *Response) GetError() error {
	if p == nil {
		return nil
	}

	if msg := strings.TrimSpace(p.ErrorMessage); msg != "" {
		return fmt.Errorf("platon api error: %s", msg)
	}

	if declineReason := strings.TrimSpace(p.DeclineReason); declineReason != "" {
		return fmt.Errorf("platon api declined: %s", declineReason)
	}

	if p.Result == nil {
		return nil
	}

	switch strings.ToUpper(strings.TrimSpace(p.Result.String())) {
	case ResultError.String():
		return fmt.Errorf("unknown platon api error")
	case ResultDeclined.String():
		return fmt.Errorf("unknown platon api decline")
	}

	return nil
}

func (p *Response) SubmerchantIDStatus() (string, bool) {
	if p == nil || p.ResponseData == nil || p.ResponseData.SubmerchantIDStatus == nil {
		return "", false
	}

	return *p.ResponseData.SubmerchantIDStatus, true
}

func UnmarshalJSONResponse(data []byte) (*Response, error) {
	var resp Response

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	return &resp, nil
}

func (p *Response) UnmarshalJSON(data []byte) error {
	type responseJSON struct {
		Status              *string         `json:"status,omitempty"`
		Action              *string         `json:"action"`
		Result              *Result         `json:"result"`
		OrderId             *string         `json:"order_id"`
		TransId             *string         `json:"trans_id"`
		TransDate           *string         `json:"trans_date"`
		ResponseData        *ResponseData   `json:"response,omitempty"`
		SubmerchantID       *string         `json:"submerchant_id,omitempty"`
		SubmerchantIDStatus *string         `json:"submerchant_id_status,omitempty"`
		Hash                *string         `json:"hash,omitempty"`
		ErrorMessage        json.RawMessage `json:"error_message"`
		DeclineReason       json.RawMessage `json:"decline_reason"`
	}

	var raw responseJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	errorMessage, err := normalizeOptionalResponseString(raw.ErrorMessage)
	if err != nil {
		return fmt.Errorf("decode error_message: %w", err)
	}
	declineReason, err := normalizeOptionalResponseString(raw.DeclineReason)
	if err != nil {
		return fmt.Errorf("decode decline_reason: %w", err)
	}

	p.Status = raw.Status
	p.Action = raw.Action
	p.Result = raw.Result
	p.OrderId = raw.OrderId
	p.TransId = raw.TransId
	p.TransDate = raw.TransDate
	responseData := raw.ResponseData
	if responseData == nil {
		if raw.SubmerchantID != nil || raw.SubmerchantIDStatus != nil || raw.Hash != nil {
			responseData = &ResponseData{
				SubmerchantID:       raw.SubmerchantID,
				SubmerchantIDStatus: raw.SubmerchantIDStatus,
				Hash:                raw.Hash,
			}
		}
	} else {
		if responseData.SubmerchantID == nil {
			responseData.SubmerchantID = raw.SubmerchantID
		}
		if responseData.SubmerchantIDStatus == nil {
			responseData.SubmerchantIDStatus = raw.SubmerchantIDStatus
		}
		if responseData.Hash == nil {
			responseData.Hash = raw.Hash
		}
	}

	p.ResponseData = responseData
	p.ErrorMessage = errorMessage
	p.DeclineReason = declineReason

	return nil
}

func normalizeOptionalResponseString(raw json.RawMessage) (string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return "", nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text), nil
	}

	var decoded interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", err
	}

	normalized, err := json.Marshal(decoded)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(normalized)), nil
}
