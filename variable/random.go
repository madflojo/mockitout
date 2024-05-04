package variable

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type RandomFunc func() string

var RandomMap = map[string]RandomFunc{
	// common
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
	"randomLongitude":   randomWrapper(gofakeit.Longitude),
	"randomLatitude":    randomWrapper(gofakeit.Latitude),
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
	// domains, emails and usernames
	"randomDomainName":   gofakeit.DomainName,
	"randomDomainSuffix": gofakeit.DomainSuffix,
	"randomEmail":        gofakeit.Email,
	"randomUserName":     gofakeit.Username,
	"randomUrl":          gofakeit.URL,
	// files and directories
	"randomFileExt": gofakeit.FileExtension,
	// stores
	"randomPrice":            randomPrice,
	"randomProduct":          gofakeit.ProductName,
	"randomProductMaterial":  gofakeit.ProductMaterial,
	"randomProduectCategory": gofakeit.ProductCategory,
	// grammar
	"randomNoun":      gofakeit.Noun,
	"randomVerb":      gofakeit.Verb,
	"randomIngverb":   gofakeit.VerbAction,
	"randomAdjective": gofakeit.Adjective,
	"randomWord":      gofakeit.Word,
	"randomWords":     randomSentence,
	"randomPhrase":    gofakeit.Phrase,
	// lorem ipsum
	"randomLoremWord":       gofakeit.LoremIpsumWord,
	"randomLoremWords":      randomLoremWords,
	"randomLoremSentence":   randomLoremSentence,
	"randomLoremSentences":  randomLoremSentences,
	"randomLoremParagraph":  randomLoremParagraph,
	"randomLoremParagraphs": randomLoremParagraphs,

	// server specific
	"hostname": getHostname,
	"goos":     getGoos,
	"goarch":   getGoarch,
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

func randomCreditCard() string {
	return gofakeit.CreditCard().Number
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getGoos() string {
	return runtime.GOOS
}

func getGoarch() string {
	return runtime.GOARCH
}

func randomPrice() string {
	return fmt.Sprintf("%f", gofakeit.Price(0, 1000))
}

func randomSentence() string {
	return gofakeit.Sentence(20)
}

func randomLoremWords() string {
	return gofakeit.LoremIpsumSentence(20)
}

func randomLoremSentence() string {
	return gofakeit.LoremIpsumSentence(1)
}

func randomLoremSentences() string {
	return gofakeit.LoremIpsumSentence(5)
}

func randomLoremParagraph() string {
	return gofakeit.LoremIpsumParagraph(1, 5, 12, "")
}

func randomLoremParagraphs() string {
	return gofakeit.LoremIpsumParagraph(3, 5, 12, "\n")
}
