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
	// json.Unmarshal([]byte(`{"agreement":{"acceptsBankTransfer":true,"acceptsCreditCard":true,"agreementID":"ed356d73-70b7-485c-7cac-c1e42b49ca86","archived":false,"clientID":"096d76ac-ac5c-489c-60ca-d0bde5aac605","freelancerID":"5275dc12-2429-42ad-536d-ccfdba2fc91c","lastAction":{"date":"2014-04-11T14:21:46.867388667-04:00","name":"accepted","userID":"5275dc12-2429-42ad-536d-ccfdba2fc91c"},"lastModified":"2014-04-11T14:21:46.867389586-04:00","lastSubAction":null,"proposedServices":"<p>Create a blog</p>\n","title":"Sample Agreement","version":1,"versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"}, "payments":[{"amountDue":1000,"amountPaid":1000,"dateExpected":"0001-01-01T00:00:00Z","id":"dd44897f-0760-41b9-4405-3fa1fa29c989","isDeposit":true,"lastAction":{"date":"2014-04-11T14:21:46.838297751-04:00","name":"accepted","userID":"5275dc12-2429-42ad-536d-ccfdba2fc91c"},"paymentItems":[{"amountDue":1000,"hours":0,"rate":0,"subtaskID":"","taskID":"","title":"Deposit"}],"title":"Deposit","versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"},{"amountDue":2000,"amountPaid":0,"dateExpected":"2014-04-18T04:00:00Z","id":"c1ac7d37-91ee-46ee-6493-c6c00774a47b","isDeposit":false,"lastAction":null,"paymentItems":[],"title":"Halfway Payment","versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"},{"amountDue":2000,"amountPaid":0,"dateExpected":"2014-04-25T04:00:00Z","id":"995b2162-5f51-416b-79c8-a0d85fb59da2","isDeposit":false,"lastAction":null,"paymentItems":[],"title":"Final Payment","versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"}], "tasks":[{ "sample":true, "dateExpected": "2014-04-11T04:00:00Z", "hours": 0, "id": "77d3ea52-c969-4b98-42b7-c2938d616b36", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "7da4396b-69b8-4505-56b3-b22aa34cf158", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Ideation", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "e40c49a3-f4b1-4ca3-470a-8b80819bcdbc", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Wireframes", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "6ceb4dd8-d586-4136-605d-a4e2b85ccb13", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Design", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86"}, { "sample":true, "dateExpected": "2014-04-18T04:00:00Z", "hours": 0, "id": "98f79d60-8fe1-4f35-4f55-f6d0bf8caa47", "isPaid": false, "lastAction": null, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "11a0bd22-4dd0-4ade-6dd4-54ba5ff49f99", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Prototype in HTML/CSS", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "fe421c0d-80a9-4487-6158-cf35528ebb17", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Mockup", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86"}, { "sample":true, "dateExpected": "2014-04-18T04:00:00Z", "hours": 0, "id": "0b47680f-40ca-483b-7803-4e5258bac5bf", "isPaid": false, "lastAction": null, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "46786360-8c22-4f0d-7570-5e6042c5ae9f", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Finish coding blog", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "18960288-3cc8-4855-4fb1-758da517c95c", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Move to production", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "e5ebc27b-8a13-441a-7869-28e237f6ee62", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Final Product", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86" }], "freelancer":{"isRegistered":true, "firstName":"Matt", "lastName":"Parker", "email":"mdparker89@gmail"},"client":{"isRegistered":true, "firstName":"Marcus", "lastName":"Ellison", "email":"marcus@gmail"}}`), &data)
	json.Unmarshal(body, &data)
	agreement := data.Agreement
	payments := data.Payments
	tasks := data.Tasks
	freelancer := data.Freelancer
	var paymentTotal float64

	for _, payment := range payments {
		if payments[0].IsDeposit {
			payment.Number -= payment.Number
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
		Tasks      []*Task                `json:"tasks"`
		Freelancer map[string]interface{} `json:"freelancer"`
		Client     map[string]interface{} `json:"client"`
	}
	// json.Unmarshal([]byte(`{"agreement":{"acceptsBankTransfer":true,"acceptsCreditCard":true,"agreementID":"ed356d73-70b7-485c-7cac-c1e42b49ca86","archived":false,"clientID":"096d76ac-ac5c-489c-60ca-d0bde5aac605","freelancerID":"5275dc12-2429-42ad-536d-ccfdba2fc91c","lastAction":{"date":"2014-04-11T14:21:46.867388667-04:00","name":"accepted","userID":"5275dc12-2429-42ad-536d-ccfdba2fc91c"},"lastModified":"2014-04-11T14:21:46.867389586-04:00","lastSubAction":null,"proposedServices":"<p>Create a blog</p>\n","title":"Sample Agreement","version":1,"versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"}, "payment":{"amountDue":2000,"amountPaid":0,"number":1,"dateExpected":"2014-04-18T04:00:00Z","id":"c1ac7d37-91ee-46ee-6493-c6c00774a47b","isDeposit":false,"lastAction":null,"paymentItems":[{"amountDue":1000,"hours":10,"rate":100,"subtaskID":"345","taskID":"123","title":"subtask"}, {"amountDue":1000,"hours":0,"rate":0,"subtaskID":"","taskID":"","title":"subtask2"}],"title":"Halfway Payment","versionID":"ed356d73-70b7-485c-7cac-c1e42b49ca86"},"tasks":[{ "sample":true, "dateExpected": "2014-04-11T04:00:00Z", "hours": 0, "id": "123", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "7da4396b-69b8-4505-56b3-b22aa34cf158", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Ideation", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "e40c49a3-f4b1-4ca3-470a-8b80819bcdbc", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Wireframes", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "6ceb4dd8-d586-4136-605d-a4e2b85ccb13", "isPaid": false, "lastAction": { "name": "completed"}, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Design", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86"}, { "sample":true, "dateExpected": "2014-04-18T04:00:00Z", "hours": 0, "id": "98f79d60-8fe1-4f35-4f55-f6d0bf8caa47", "isPaid": false, "lastAction": null, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "11a0bd22-4dd0-4ade-6dd4-54ba5ff49f99", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Prototype in HTML/CSS", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "fe421c0d-80a9-4487-6158-cf35528ebb17", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Mockup", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86"}, { "sample":true, "dateExpected": "2014-04-18T04:00:00Z", "hours": 0, "id": "0b47680f-40ca-483b-7803-4e5258bac5bf", "isPaid": false, "lastAction": null, "subTasks": [ { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "46786360-8c22-4f0d-7570-5e6042c5ae9f", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Finish coding blog", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "18960288-3cc8-4855-4fb1-758da517c95c", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Move to production", "versionID": ""}, { "dateExpected": "0001-01-01T00:00:00Z", "hours": 0, "id": "e5ebc27b-8a13-441a-7869-28e237f6ee62", "isPaid": false, "lastAction": null, "subTasks": null, "title": "Review", "versionID": "" }], "title": "Final Product", "versionID": "ed356d73-70b7-485c-7cac-c1e42b49ca86" }], "freelancer":{"isRegistered":true, "firstName":"Matt", "lastName":"Parker", "email":"mdparker89@gmail"},"client":{"isRegistered":true, "firstName":"Marcus", "lastName":"Ellison", "email":"marcus@gmail"}}`), &data)
	json.Unmarshal(body, &data)
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
