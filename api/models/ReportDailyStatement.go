package models

import (
    "encoding/json"
    "time"

    "github.com/jlaffaye/ftp"
)

// DEFINING JSON STRUCTURE
type DailyStatementReportResponse struct {
    ReportResponse  DailyStatementReportData `json:"response"`
}

type DailyStatementReportData struct {
    ReportData      []DailyStatementReportArray `json:"data"`
    DocumentName    string  `json:"document_name"`
    DocumentNumber  int     `json:"document_number"`
    SoldToNumber    string  `json:"sold_to_number"`
    SoldToName      string  `json:"sold_to_name"`
    SoldToAddress   string  `json:"sold_to_address"`
    ShipToNumber    string  `json:"ship_to_number"`
    ShipToName      string  `json:"ship_to_name"`
    ShipToAddress   string  `json:"ship_to_address"`
    Period          string  `json:"period"`
    TotalVolume     float32 `json:"total_volume"`
    TotalAmount     float32 `json:"total_amount"`
}

type DailyStatementReportArray struct {
    BatchNumber     int     `json:"batch_no"`
    TransactionType string  `json:"transaction_type"`
    TransactionDate string  `json:"transaction_date"`
    TransactionTime string  `json:"transaction_time"`
    CardId          string  `json:"card_id"`
    ReferenceNumber string  `json:"reference_no"`
    Product         string  `json:"product"`
    Volume          float32 `json:"volume"`
    Ppn             float32 `json:"ppn"`
    Amount          float32 `json:"amount"`
}

type DailyStatementReportFileParam struct {
    ReportFulldate  string
    SoldToNumber    string
    ShipToNumber    string
    FtpConn         *ftp.ServerConn
}

// UNBOXING PROCESS
func (reportData *DailyStatementReportResponse) DailyStatementFile(args DailyStatementReportFileParam) (*DailyStatementReportResponse, error) {
    var reportName = "DailyStatementReport"

    // DATE FORMATTING
    layout      := "20060102"
    t, _        := time.Parse(layout, args.ReportFulldate)
    reportYear  := t.Format("2006")
    reportMonth := t.Format("January")
    reportDate  := t.Format("02")
    outputDate  := reportDate+" "+reportMonth+" "+reportYear

    // RETRIEVE DATA FROM JSON
    FilePath := "Daily/"+outputDate+"/"+reportName+"_"+args.SoldToNumber+"_"+args.ShipToNumber+"_"+args.ReportFulldate+".json"
    
    byteValue, err := ReadJsonFile(FilePath, args.FtpConn)
    if err != nil {
        return nil, err
    }
    
    // MANAGE DATA
    json.Unmarshal(byteValue, &reportData)

    return reportData, nil
}