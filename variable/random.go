package variable

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type RandomFunc func() string

var RandomMap = map[string]RandomFunc{
	// comon
	"guid":         gofakeit.UUID,
	"timestamp":    timeNowUnixString,
	"isoTimestamp": timeNowIso,
	// text, numbers and colors
	"randomAlphaNumberic": gofakeit.Letter,
	"randomBoolean":       randomWrapper(gofakeit.Bool),
	"randomInt":           randomInt,
	"randomColor":         gofakeit.Color,
	"randomHexColor":      gofakeit.HexColor,
	"randomAbbreviation":  gofakeit.HackerAbbreviation,
	// internet and ip addresses
	"randomIPV4":       gofakeit.IPv4Address,
	"randomIPV6":       gofakeit.IPv6Address,
	"randomMacAddress": gofakeit.MacAddress,
	"randomPassword":   randomPassword,
	"randomUserAgent":  gofakeit.UserAgent,
	"randomSemver":     gofakeit.AppVersion,
	// names
	"randomFirstName":  gofakeit.FirstName,
	"randomLastName":   gofakeit.LastName,
	"randomNamePrefix": gofakeit.NamePrefix,
	"randomNameSuffix": gofakeit.NameSuffix,
	// profession
	"randomJobTitle": gofakeit.JobTitle,
	"randomJobType":  gofakeit.JobDescriptor,
	// phone, address and location
	"randomPhoneNumber": gofakeit.PhoneFormatted,
	"randomCity":        gofakeit.City,
	"randomStreetName":  gofakeit.Street,
	"randomCountry":     gofakeit.Country,
	"randomCountryCode": gofakeit.CountryAbr,
	"randomLongitude":   randomLongitude,
	"randomLatitude":    randomLatitude,
	// images (TODO?)
	// finance
	"randomCreditCard":   randomCreditCard,
	"randomCurrencyCode": gofakeit.CurrencyShort,
	"randomCurrencyName": gofakeit.CurrencyLong,
	"randomBitcoin":      gofakeit.BitcoinAddress,
	// business
	"randomCompany":       gofakeit.Company,
	"randomCompanySuffix": gofakeit.CompanySuffix,
	"randomBs":            gofakeit.BS,
	// catchphrases
	"randomCatchPhrase":          gofakeit.Phrase,
	"randomCatchPhraceAdjective": gofakeit.Adjective,
	"randomCatchPhraseNoun":      gofakeit.Noun,
	// databases
}

func randomWrapper[T any](randFunc func() T) RandomFunc {
	return func() string {
		return fmt.Sprintf("%v", randFunc())
	}
}

func timeNowUnixString() string {
	return fmt.Sprint(time.Now().Unix())
}

func timeNowIso() string {
	return time.Now().Format(time.RFC3339)
}

func randomInt() string {
	return fmt.Sprintf("%v", gofakeit.Number(0, 100))
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 12)
}

func randomLongitude() string {
	// TODO:
	return ""
}

func randomLatitude() string {
	// TODO:
	return ""
}

func randomCreditCard() string {
	return gofakeit.CreditCard().Number
}
