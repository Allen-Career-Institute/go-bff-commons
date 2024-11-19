package commons

import (
	"bff-service/internal/models/commons"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetAmountInr(amount float32) string {
	ac := accounting.Accounting{Symbol: commons.Rs, Precision: 0}
	return ac.FormatMoney(amount)
}

func GetAmount(amount float32) float32 {
	if amount >= -1 && amount < 1 {
		amount = 0
	}
	return float32(roundFloat(float64(amount), 0))
}

func GetOfferDetail(old, new float32) string {
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	diff := old - new
	delta := int((diff / old) * 100)
	if delta > 5.0 {
		return strconv.Itoa(delta) + "% OFF"
	}
	return ac.FormatMoney(old-new) + " OFF"
}

func FormatDate(fromLayout, toLayout, date string) string {
	dt, err := time.Parse(fromLayout, date)
	if err != nil {
		return date
	}
	return dt.Format(toLayout)
}

func MaskUserInfo(field string) string {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	phoneRegex := `^[0-9]{10}$`
	if regexp.MustCompile(emailRegex).MatchString(field) {
		parts := strings.Split(field, "@")
		username := parts[0]
		domain := parts[1]
		maskedUsername := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
		return maskedUsername + "@" + domain
	} else if regexp.MustCompile(phoneRegex).MatchString(field) {
		re := regexp.MustCompile(`(\d)(\d{5})(\d{3})`)
		return re.ReplaceAllString(field, "$1******$3")
	}
	return field
}

func ValidateField(field string, fieldType string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	phoneRegex := `^[0-9]{10}$`
	nanoIDRegex := `^[0-9a-zA-Z]{21}$`

	switch fieldType {
	case commons.EmailID:
		if regexp.MustCompile(emailRegex).MatchString(field) {
			return true
		}
		return false
	case commons.PhoneNumber:
		if regexp.MustCompile(phoneRegex).MatchString(field) {
			return true
		}
		return false
	case commons.ID:
		if regexp.MustCompile(nanoIDRegex).MatchString(field) {
			return true
		}
		return false
	default:
		return false
	}
}

func RemoveEmptyStrings(slice []string) []string {
	var result []string
	for _, str := range slice {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func GetBatchID(e echo.Context) string {
	batchIDStr := e.QueryParam(commons.BatchIDParam)
	batchIDs := RemoveEmptyStrings(strings.Split(batchIDStr, ","))

	if len(batchIDs) > 0 {
		return batchIDs[0]
	}

	return commons.EmptyStr
}

func GetSelectedBatchIds(e echo.Context) []string {
	batchIds := e.QueryParam(commons.SelectedBatchList)
	if batchIds == "" {
		return []string{}
	}

	return RemoveEmptyStrings(strings.Split(batchIds, ","))
}
