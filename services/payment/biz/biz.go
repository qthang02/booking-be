package paymentbiz

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/services/payment/helper"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PaymentBiz struct {
}

func NewPaymentBiz() *PaymentBiz {
	return &PaymentBiz{}
}

func (biz *PaymentBiz) CreatePaymentURL(c echo.Context) error {
	amount := c.QueryParam("amount")
	orderInfo := c.QueryParam("orderInfo")

	amountInt, _ := strconv.Atoi(amount)

	params := url.Values{}
	params.Set("vnp_Version", "2.1.0")
	params.Set("vnp_Command", "pay")
	params.Set("vnp_TmnCode", helper.VnpTmnCode)
	params.Set("vnp_Amount", strconv.Itoa(amountInt*100))
	params.Set("vnp_CurrCode", "VND")
	params.Set("vnp_TxnRef", generateOrderID())
	params.Set("vnp_OrderInfo", orderInfo)
	params.Set("vnp_OrderType", "other")
	params.Set("vnp_Locale", "vn")
	params.Set("vnp_ReturnUrl", helper.VnpReturnURL)
	params.Set("vnp_IpAddr", getClientIP(c))
	params.Set("vnp_CreateDate", getVNPayDateTimeFormat(time.Now()))
	params.Set("vnp_ExpireDate", getVNPayDateTimeFormat(time.Now().Add(time.Minute*15)))

	signData := createSignData(params)
	signature := createSignature(signData)
	params.Set("vnp_SecureHash", signature)

	paymentURL := helper.VnpPayURL + "?" + params.Encode()

	return c.JSON(http.StatusOK, map[string]string{"paymentURL": paymentURL})
}

func (biz *PaymentBiz) HandleVNPayReturn(c echo.Context) error {
	queryParams := c.QueryParams()

	if !verifySignature(queryParams) {
		return c.String(http.StatusBadRequest, "Invalid signature")
	}

	vnpResponseCode := queryParams.Get("vnp_ResponseCode")
	vnpTransactionStatus := queryParams.Get("vnp_TransactionStatus")

	if vnpResponseCode == "00" && vnpTransactionStatus == "00" {
		return c.String(http.StatusOK, "Thanh toán thành công")
	} else {
		return c.String(http.StatusBadRequest, "Thanh toán thất bại")
	}
}

func getVNPayDateTimeFormat(timer time.Time) string {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return timer.Format("20060102150405")
	}
	return timer.In(loc).Format("20060102150405")
}

func generateOrderID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getClientIP(c echo.Context) string {
	return c.RealIP()
}

func hmacSHA512(key []byte, data string) string {
	h := hmac.New(sha512.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func createSignData(params url.Values) string {
	var signData string
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		signData += k + "=" + url.QueryEscape(params.Get(k)) + "&"
	}
	signData = strings.TrimSuffix(signData, "&")
	return signData
}

func createSignature(data string) string {
	signature := strings.ToLower(hmacSHA512([]byte(helper.VnpHashSecret), data))
	return signature
}

func verifySignature(params url.Values) bool {
	receivedHash := params.Get("vnp_SecureHash")
	signData := createSignData(params)
	calculatedHash := createSignature(signData)

	return receivedHash == calculatedHash
}
