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
type DailySalesReportResponse struct {
    ReportResponse  DailySalesReportData `json:"response"`
}

type DailySalesReportData struct {
    ReportData      []DailySalesReportArray `json:"data"`
    ReportName      string  `json:"document_name"`
    DocumentNumber  int     `json:"document_number"`
    OverallAmount   float32 `json:"total_amount"`
    OverallVolume   float32 `json:"total_volume"`
    Period          string  `json:"period"`
}

type DailySalesReportArray struct {
    PaddedMcmsId        string  `json:"padded_mcms_id"`
    CorporateName       string  `json:"corporate_name"`
    SiteCode            string  `json:"site_code"`
    SiteName            string  `json:"site_name"`
    SaleProductCode     string  `json:"sale_product_code"`
    SaleProdutName      string  `json:"sale_product_name"`
    TotalDailyAmount    float32 `json:"total_daily_amount"`
    TotalDailyVolume    float32 `json:"total_daily_volume"`
    AccontManager       string  `json:"account_manager"`
    CustomerBanding     string  `json:"customer_banding"`
    CCID                int     `json:"cc_id"`
}

// UNBOXING PROCESS
func (daily_sales_report_data *DailySalesReportResponse) DailySalesFile(report_fulldate string) (*DailySalesReportResponse, error) {
    var err error
    var report_name = "DailySalesReport"
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
    layout                  := "20060102"
    t, _                    := time.Parse(layout, report_fulldate)
    report_year             := t.Format("2006")
    report_month            := t.Format("January")
    report_date             := t.Format("02")
    output_date             := report_date+" "+report_month+" "+report_year

    // RETRIEVE DATA FROM JSON
    FilePath := "Daily/"+output_date+"/"+report_name+"_"+report_fulldate+".json"
    jsonFile, err := conn.Retr(FilePath)
    if err != nil {
        log.Println(err.Error())
    } else {
        log.Println("Successfully Retrieved "+report_name+", date: "+report_fulldate)
    }

    // READ FILE
    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        log.Println(err.Error())
    } else {
        log.Println("Successfully Read "+report_name+", date: "+report_fulldate)
    }
    jsonFile.Close()

    // MANAGE DATA
    json.Unmarshal(byteValue, &daily_sales_report_data)

    return daily_sales_report_data, nil
}