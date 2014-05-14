package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"time"
)

type Agreement struct {
	Name             string `json:"title"`
	ProposedServices string `json:"proposedServices"`
}

func unescaped(x string) interface{} { return template.HTML(x) }

func agreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var data struct {
		Agreement  *Agreement               `json:"agreement"`
		Payments   []*Payment               `json:"payments"`
		Tasks      []map[string]interface{} `json:"tasks"`
		Freelancer map[string]interface{}   `json:"freelancer"`
	}
	json.Unmarshal(body, &data)
	agreement := data.Agreement
	payments := data.Payments
	tasks := data.Tasks
	freelancer := data.Freelancer
	var paymentTotal float64

	for _, payment := range payments {
		if !payments[0].IsDeposit {
			payment.Number += 1
		}
		paymentTotal += payment.AmountDue
	}

	m := map[string]interface{}{
		"PAYMENT_NUMBER":   1,
		"Agreement":        agreement,
		"ProposedServices": unescaped(agreement.ProposedServices),
		"Tasks":            tasks,
		"DateNow":          time.Now().Format("January 2, 2006"),
		"Payments":         payments,
		"Total":            paymentTotal,
		"Freelancer":       freelancer,
	}
	format := func(date time.Time) string {
		return date.Format("Jan 2, 2006")
	}
	formatDateString := func(date string) string {
		t, _ := time.Parse(time.RFC3339, date)
		return t.Format("Jan 2, 2006")
	}
	numberSuffix := func(number int64) string {
		if number == 1 {
			return "st"
		}
		if number == 2 {
			return "nd"
		}
		if number == 3 {
			return "rd"
		}
		return "th"
	}

	var html bytes.Buffer
	newAgreementTpl, _ := template.New("agreement_summary.html").Funcs(template.FuncMap{"formatDate": format, "formatDateString": formatDateString, "numberSuffix": numberSuffix}).ParseFiles(
		"templates/agreement_summary.html",
	)
	newAgreementTpl.Execute(&html, m)
	return html.Bytes(), nil, 200
}

type Payment struct {
	ID           string         `json:"id"`
	VersionID    string         `json:"versionID"`
	Title        string         `json:"title"`
	DateExpected time.Time      `json:"dateExpected"`
	PaymentItems []*PaymentItem `json:"paymentItems"`
	IsDeposit    bool           `json:"isDeposit"`
	AmountDue    float64        `json:"amountDue"`
	AmountPaid   float64        `json:"amountPaid"`
	Number       int64          `json:"number"`
}

type PaymentItem struct {
	AmountDue float64 `json:"amountDue"`
	Hours     float64 `json:"hours,omitempty"`
	Rate      float64 `json:"rate"`
	SubTaskID string  `json:"subTaskID"`
	TaskID    string  `json:"taskID"`
	Title     string  `json:"title"`
}
type InvoiceItem struct {
	Title     string
	ID        string
	AmountDue float64
	Hours     float64
	Rate      float64
	Items     []*PaymentItem
}
type Task struct {
	ID           string    `json:"id"`
	VersionID    string    `json:"versionID"`
	IsPaid       bool      `json:"isPaid"`
	Hours        float64   `json:"hours"`
	SubTasks     []*Task   `json:"subTasks"`
	Title        string    `json:"title"`
	DateExpected time.Time `json:"dateExpected"`
}

func invoice(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var data struct {
		Agreement  *Agreement             `json:"agreement"`
		Payment    *Payment               `json:"payment"`
		Payments   []*Payment             `json:"payments"`
		Tasks      []*Task                `json:"tasks"`
		Freelancer map[string]interface{} `json:"freelancer"`
		Client     map[string]interface{} `json:"client"`
	}
	json.Unmarshal(body, &data)
	data.Payment.Number += 1
	agreement := data.Agreement
	tasks := data.Tasks
	freelancer := data.Freelancer
	paymentItems := make(map[string]*InvoiceItem)
	for _, item := range data.Payment.PaymentItems {
		invoiceItem := &InvoiceItem{}
		if t, ok := paymentItems[item.TaskID]; ok {
			invoiceItem = t
		}
		if item.SubTaskID == "" {
			invoiceItem.AmountDue = item.AmountDue
			invoiceItem.Hours = item.Hours
			invoiceItem.Rate = item.Rate
			for _, task := range tasks {
				if task.ID == item.TaskID {
					for _, subTask := range task.SubTasks {
						newItem := &PaymentItem{}
						newItem.Title = subTask.Title
						invoiceItem.Items = append(invoiceItem.Items, newItem)
					}
				}
			}
			if len(invoiceItem.Items) == 0 {
				invoiceItem.Items = append(invoiceItem.Items, item)
			}
		} else {
			invoiceItem.AmountDue += item.AmountDue
			invoiceItem.Hours += item.Hours
			invoiceItem.Rate = invoiceItem.AmountDue / invoiceItem.Hours
			invoiceItem.Items = append(invoiceItem.Items, item)
		}
		paymentItems[item.TaskID] = invoiceItem
	}
	for key, value := range paymentItems {
		for _, task := range tasks {
			if task.ID == key {
				paymentItems[task.Title] = value
				delete(paymentItems, key)
				break
			}
		}
	}

	m := map[string]interface{}{
		"PAYMENT_NUMBER":   1,
		"Agreement":        agreement,
		"ProposedServices": unescaped(agreement.ProposedServices),
		"Tasks":            tasks,
		"PaymentDate":      time.Now().Add(24 * 7 * time.Hour).Format("January 2, 2006"),
		"DateNow":          time.Now().Format("January 2, 2006"),
		"paymentItems":     paymentItems,
		"Payment":          data.Payment,
		"Freelancer":       freelancer,
		"Client":           data.Client,
	}
	format := func(date string) string {
		t, _ := time.Parse(time.RFC3339, date)
		return t.Format("Jan 2, 2006")
	}

	newAgreementTpl, _ := template.New("invoice.html").Funcs(template.FuncMap{"formatDate": format}).ParseFiles(
		"templates/invoice.html",
	)
	var html bytes.Buffer
	newAgreementTpl.Execute(&html, m)
	return html.Bytes(), nil, 200
}
