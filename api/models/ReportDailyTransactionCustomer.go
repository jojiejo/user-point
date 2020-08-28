package models

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

// DEFINING JSON STRUCTURE
type DailyTransactionCustomerReportResponse struct {
    ReportResponse  DailyTransactionCustomerReportData `json:"response"`
}

type DailyTransactionCustomerReportData struct {
    ReportData      []DailyTransactionCustomerReportArray `json:"data"`
    ReportName      string  `json:"report_name"`
    DocumentNumber  string  `json:"document_number"`
    OverallAmount   float32 `json:"overall_amount"`
    OverallVolume   float32 `json:"overall_volume"`
    Period          string  `json:"period"`
}

type DailyTransactionCustomerReportArray struct {
    PaddedMcmsId            string  `json:"padded_mcms_id"`
    PayerNumber             string  `json:"payer_number"`
    CorporateName           string  `json:"corporate_name"`
    CorporateAddress        string  `json:"corporate_address"`
    SubCorporateName        string  `json:"sub_corporate_name"`
    TransactionIndicator    string  `json:"transaction_indicator"`
    CardNumber              string  `json:"card_number"`
    CardHolderName          string  `json:"card_holder_name"`
    Vrn                     string  `json:"vrn"`
    Odometer                string  `json:"odometer"`
    CardGroupName           string  `json:"card_group_name"`
    TransactionDate         string  `json:"transaction_date"`
    TransactionTime         string  `json:"transaction_time"`
    Network                 string  `json:"network"`
    SiteCode                string  `json:"site_code"`
    SiteName                string  `json:"site_name"`
    SaleProductName         string  `json:"sale_product_name"`
    TotalDailyVolume        float32 `json:"total_daily_volume"`
    TotalDailyAmount        float32 `json:"total_daily_amount"`
    ReceiptNumber           string  `json:"receipt_number"`
    CCID                    int     `json:"cc_id"`
}

// UNBOXING PROCESS
func (daily_transaction_customer_report_data *DailyTransactionCustomerReportResponse) DailyTransactionCustomerFile(report_date string, account_name string) (*DailyTransactionCustomerReportResponse, error) {
    var err error

    // OPEN FILE
    jsonFile2, err := os.Open("report/report-daily-transaction-customer.json")

    defer jsonFile2.Close()

    // CHECK IF NULL
    if err != nil {
        return &DailyTransactionCustomerReportResponse{}, err
    }

    // READ FILE
    fmt.Println("Successfully Opened "+report_date+account_name+".json")
    byteValue, _ := ioutil.ReadAll(jsonFile2)

    // MANAGE DATA
    json.Unmarshal(byteValue, &daily_transaction_customer_report_data)

    return daily_transaction_customer_report_data, nil
}