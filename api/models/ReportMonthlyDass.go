package models

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "log"
    "time"

    "github.com/jlaffaye/ftp"
)

// DEFINING JSON STRUCTURE
type MonthlyDassReportResponse struct {
    ReportResponse  MonthlyDassReportData `json:"response"`
}

type MonthlyDassReportData struct {
    ReportName      string  `json:"document_name"`
    DocumentNumber  int     `json:"document_number"`
    DateFrom        string  `json:"date_from"`
    DateTo          string  `json:"date_to"`
}

type MonthlyDassReportDataArray struct {
    // CorporateId                 string  `json:"corporate_id"`
    // SiteId                      string  `json:"site_id"`
    // TotalAmount                 string  `json:"total_amount"`
    // TotalVolume                 string  `json:"total_volume"`
    CMITRANSId                  string  `json:"CMITRANSId"`
    ReportingMonth              string  `json:"ReportingMonth"`
    DeliveryMonth               string  `json:"DeliveryMonth"`
    Currency                    string  `json:"Currency"`
    CardColcoCode               string  `json:"CardColcoCode"`
    ShellCollectingCompany      string  `json:"ShellCollectingCompany"`
    CardCustomerType            string  `json:"CardCustomerType"`
    CardProductCode             string  `json:"CardProductCode"`
    Material                    string  `json:"Material"`
    IssuerCode                  string  `json:"IssuerCode"`
    CardTypeCode                string  `json:"CardTypeCode"`
    CardAccount                 string  `json:"CardAccount"`
    CardStationCode             string  `json:"CardStationCode"`
    CardStationDescription      string  `json:"CardStationDescription"`
    CardAcceptanceLocation      string  `json:"CardAcceptanceLocation"`
    CardDelcoCode               string  `json:"CardDelcoCode"`
    ShellDeliveryCompany        string  `json:"ShellDeliveryCompany"`
    NetworkCode                 string  `json:"NetworkCode"`
    AcquirerCode                string  `json:"AcquirerCode"`
    TransactionType             string  `json:"TransactionType"`
    RecordType                  string  `json:"RecordType"`
    PriceType                   string  `json:"PriceType"`
    SettlementType              string  `json:"SettlementType"`
    Uom                         string  `json:"Uom"`
    Quantity                    float32 `json:"Quantity"`
    VolumeL                     float32 `json:"VolumeL"`
    InvoiceGrossProceeds        float32 `json:"InvoiceGrossProceeds"`
    InvoiceVatAmount            float32 `json:"InvoiceVatAmount"`
    InvoiceDiscountAmount       float32 `json:"InvoiceDiscountAmount"`
    InvoiceNetAmount            float32 `json:"InvoiceNetAmount"`
    ReferenceGrossProceeds      float32 `json:"ReferenceGrossProceeds"`
    ReferenceVatAmount          float32 `json:"ReferenceVatAmount"`
    ReferenceDiscountAmount     float32 `json:"ReferenceDiscountAmount"`
    ReferenceNetAmount          float32 `json:"ReferenceNetAmount"`
    PumpGrossProceeds           float32 `json:"PumpGrossProceeds"`
    PumpVatAmount               float32 `json:"PumpVatAmount"`
    PumpDiscountAmount          float32 `json:"PumpDiscountAmount"`
    PumpNetAmount               float32 `json:"PumpNetAmount"`
    PumpReferencePriceVariance  float32 `json:"PumpReferencePriceVariance"`
    EffectiveFuelCardDiscount   float32 `json:"EffectiveFuelCardDiscount"`
    TransactionNo               float32 `json:"TransactionNo"`
    ItemDutyExemptGiven         float32 `json:"ItemDutyExemptGiven"`
    ItemTaxExemptGiven          float32 `json:"ItemTaxExemptGiven"`
    DealerContrDelcoCur         float32 `json:"DealerContrDelcoCur"`
}

type MonthlyDassReportSuplementArray struct {
    AccountManager  string `json:"AccountManager"`
    AccountNo       string `json:"AccountNo"`
    AccountName     string `json:"AccountName"`
}

// UNBOXING PROCESS
func (monthly_dass_report_data *MonthlyDassReportResponse) MonthlyDassFile(report_fulldate string) (*MonthlyDassReportResponse, error) {
    var err error
    var report_name = "DASSReport"
    var FTP_HOST = os.Getenv("FTP_HOST")

    // INITIALIZE FTP
    conn, err := ftp.Dial(FTP_HOST)
    if err != nil {
        log.Println(err.Error())
    } else {
        log.Println("FTP Connection Success")
    }

    // FTP LOGIN
    err = conn.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
    if err != nil {
        log.Println(err.Error())
    } else {
        log.Println("FTP Login Success")
    }

    // DATE FORMATTING
    layout                  := "200601"
    t, _                    := time.Parse(layout, report_fulldate)
    report_year             := t.Format("2006")
    report_month            := t.Format("January")
    output_date             := report_month+" "+report_year

    // RETRIEVE DATA FROM JSON
    FilePath := "Monthly/"+output_date+"/"+report_name+"_"+report_month+"_"+report_year+".json"
    jsonFile, err := conn.Retr(FilePath)
    if err != nil {
        log.Println(err.Error())
        return monthly_dass_report_data, nil
    } else {
        log.Println("Successfully Retrieved "+report_name+", date: "+output_date)
    }

    // READ FILE
    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        log.Println(err.Error())
    } else {
        log.Println("Successfully Read "+report_name+", date: "+output_date)
    }
    jsonFile.Close()

    // MANAGE DATA
    json.Unmarshal(byteValue, &monthly_dass_report_data)

    return monthly_dass_report_data, nil
}