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
type MonthlySalesReportResponse struct {
    ReportResponse  MonthlySalesReportData `json:"response"`
}

type MonthlySalesReportData struct {
    ReportData      []MonthlySalesReportArray `json:"data"`
    ReportName      string  `json:"document_name"`
    DocumentNumber  int     `json:"document_number"`
    OverallAmount   float32 `json:"total_amount"`
    OverallVolume   float32 `json:"total_volume"`
    DateFrom        string  `json:"date_from"`
    DateTo          string  `json:"date_to"`
}

type MonthlySalesReportArray struct {
    CCID                int     `json:"cc_id"`
    SaleProductCode     string  `json:"sale_product_code"`
    TotalMonthlyAmount  float32 `json:"total_monthly_amount"`
    TotalMonthlyVolume  float32 `json:"total_monthly_volume"`
    CorporateName       string  `json:"corporate_name"`
    PaddedMcmsId        string  `json:"padded_mcms_id"`
    PayerNumber         string  `json:"payer_number"`
    SaleProdutName      string  `json:"sale_product_name"`
    SiteCode            string  `json:"site_code"`
    SiteName            string  `json:"site_name"`
    AccontManager       string  `json:"account_manager"`
    CustomerBanding     string  `json:"customer_banding"`
}

// UNBOXING PROCESS
func (monthly_sales_report_data *MonthlySalesReportResponse) MonthlySalesFile(report_fulldate string) (*MonthlySalesReportResponse, error) {
    var err error
    var report_name = "MonthlySalesReport"
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
    json.Unmarshal(byteValue, &monthly_sales_report_data)

    return monthly_sales_report_data, nil
}