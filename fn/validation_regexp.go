package fn

/*

ASCII: caracteres ASCII
AlphaNumeric: caracteres alfanumericos
Alpha: caracteres alfabeticos
AlphaSpace: caracteres alfanumericos y espacios
AlphaUnicodeNumeric: caracteres alfanumericos y Unicode
AlphaUnicode: caracteres alfabeticos y Unicode
Base32: la cadena codificada en base32
Base64RawURL: una cadena codificada en base64 sin el caracter '=' en el final
Base64: una cadena codificada en base64
Base64URL:  una cadena codificada en base64 con el caracter '=' en el final
Bic: un número BIC
BtcAddressLowerBech32: Dirección Bitcoin en formato Bech32 en minúsculas (https://en.bitcoin.it/wiki/Bech32)
BtcAddress: Dirección Bitcoin
BtcAddressUpperBech32: Dirección Bitcoin en formato Bech32 en mayúsculas (https://en.bitcoin.it/wiki/Bech32)
Cron: Entrada de cron (https://en.wikipedia.org/wiki/Cron)
Cve: Un número CVE (https://cve.mitre.org/cve/identifiers/syntaxchange.html)
DataURI:
DnsRFC1035Label:
Ein:
Email:
EthAddressLower:
EthAddress:
EthAddressUpper:
FqdnRFC1123: Igual que hostnameRFC1123, pero debe contener un TLD no numérico (que puede terminar en '.').
HTMLEncoded:
HTML:
HexColor:
Hexadecimal:
HostnameRFC1123: Acepta nombres de host que comienzan con un dígito https://tools.ietf.org/html/rfc1123
HostnameRFC952: https://tools.ietf.org/html/rfc952
Hsl:
Hsla:
ISBN10:
ISBN13:
ISSN:
JWT: JSON Web Token (https://tools.ietf.org/html/rfc7519)
Latitude:
Longitude:
Md4:
Md5:
MongodbConnString:
MongodbId:
Multibyte:
Numeric:
PrintableASCII:
Rgb: Color RGB
Rgba: Color RGBA
Ripemd128:
Ripemd160:
SSN:
Semver: Versionado Semántico (https://semver.org/)
Sha256: Cadena codificada en SHA256
Sha384: Cadena codificada en SHA384
Sha512: Cadena codificada en SHA512
SplitParams:
Tiger128:
Tiger160:
Tiger192:
ULID:
URLEncoded:
UUID3RFC4122:
UUID3:
UUID4RFC4122:
UUID4:
UUID5RFC4122:
UUID5:
UUIDRFC4122:
UUID:


*/

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
)

func lazyRegexCompile(str string) func() *regexp.Regexp {
	var regex *regexp.Regexp
	var once sync.Once
	return func() *regexp.Regexp {
		once.Do(func() {
			regex = regexp.MustCompile(str)
		})
		return regex
	}
}

var rexMap = map[string]func() *regexp.Regexp{
	"ASCII":               lazyRegexCompile("^[\x00-\x7F]*$"),
	"AlphaNumeric":        lazyRegexCompile("^[a-zA-Z0-9]+$"),
	"Alpha":               lazyRegexCompile("^[a-zA-Z]+$"),
	"AlphaSpace":          lazyRegexCompile("^[a-zA-Z ]+$"),
	"AlphaUnicodeNumeric": lazyRegexCompile("^[a-zA-Z0-9]+$"),
	"AlphaUnicode":        lazyRegexCompile("^[a-zA-Z]+$"),
	"Base32": lazyRegexCompile(
		"^(?:[A-Z2-7]{8})*(?:[A-Z2-7]{2}={6}|[A-Z2-7]{4}={4}|[A-Z2-7]{5}={3}|[A-Z2-7]{7}=|[A-Z2-7]{8})$",
	),
	"Base64RawURL": lazyRegexCompile("^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2,4})$"),
	"Base64": lazyRegexCompile(
		"^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$",
	),
	"Base64URL":             lazyRegexCompile("^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"),
	"Bic":                   lazyRegexCompile(`^[A-Za-z]{6}[A-Za-z0-9]{2}([A-Za-z0-9]{3})?$`),
	"BtcAddressLowerBech32": lazyRegexCompile(`^bc1[02-9ac-hj-np-z]{7,76}$`),
	"BtcAddress":            lazyRegexCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`),
	"BtcAddressUpperBech32": lazyRegexCompile(`^BC1[02-9AC-HJ-NP-Z]{7,76}$`),
	"Cron": lazyRegexCompile(
		`^(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|((\*|\d+)(\/|-)\d+)|\d+|\*) ?){5,7})`,
	),
	"CVE":             lazyRegexCompile(`^CVE-(1999|2\d{3})-(0[^0]\d{2}|0\d[^0]\d{1}|0\d{2}[^0]|[1-9]{1}\d{3,})$`),
	"DataURI":         lazyRegexCompile(`^data:((?:\w+\/(?:([^;]|;[^;]).)+)?)`),
	"DnsRFC1035Label": lazyRegexCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`),
	"Ein":             lazyRegexCompile(`^(\d{2}-\d{7})$`),
	"Email": lazyRegexCompile(
		`^(?:(?:(?:(?:[a-zA-Z]|\d|[!#\$%&'\*\+\-\/=\?\^_` + "`" + `\{\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\.(?:[a-zA-Z]|\d|[!#\$%&'\*\+\-\/=\?\^_` + "`" + `\{\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\x22)(?:(?:(?:(?:\x20|\x09)*(?:\x0d\x0a))?(?:\x20|\x09)+)?(?:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(?:(?:[\x01-\x09\x0b\x0c\x0d-\x7f]|[\x{00A0}-\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\x20|\x09)*(?:\x0d\x0a))?(\x20|\x09)+)?(?:\x22))))@(?:(?:(?:[a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(?:(?:[a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])(?:[a-zA-Z]|\d|-|\.|~|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])*(?:[a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])))\.)+(?:(?:[a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(?:(?:[a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])(?:[a-zA-Z]|\d|-|\.|~|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])*(?:[a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])))\.?$`,
	),
	"EthAddressLower": lazyRegexCompile(`^0x[0-9a-f]{40}$`),
	"EthAddress":      lazyRegexCompile(`^0x[0-9a-fA-F]{40}$`),
	"EthAddressUpper": lazyRegexCompile(`^0x[0-9A-F]{40}$`),
	"FqdnRFC1123": lazyRegexCompile(
		`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$`,
	),
	"HTMLEncoded":     lazyRegexCompile(`&#[x]?([0-9a-fA-F]{2})|(&gt)|(&lt)|(&quot)|(&amp)+[;]?`),
	"HTML":            lazyRegexCompile(`<[/]?([a-zA-Z]+).*?>`),
	"HexColor":        lazyRegexCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`),
	"Hexadecimal":     lazyRegexCompile(`^(0[xX])?[0-9a-fA-F]+$`),
	"HostnameRFC1123": lazyRegexCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`),
	"HostnameRFC952":  lazyRegexCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`),
	"HSL": lazyRegexCompile(
		`^hsl\x28\s*(?:0|[1-9]\d*)\s*,\s*(?:(?:0|[1-9]\d*)\s*%)\s*,\s*(?:(?:0|[1-9]\d*)\s*%)\s*\x29$`,
	),
	"HSLA": lazyRegexCompile(
		`^hsla\x28\s*(?:0|[1-9]\d*)\s*,\s*(?:(?:0|[1-9]\d*)\s*%)\s*,\s*(?:(?:0|[1-9]\d*)\s*%)\s*,\s*(?:(?:0.[1-9]*)|[01])\s*\x29$`,
	),
	"ISBN10":    lazyRegexCompile(`^(?:[0-9]{9}X|[0-9]{10})$`),
	"ISBN13":    lazyRegexCompile(`^(?:(?:97(?:8|9))[0-9]{10})$`),
	"ISSN":      lazyRegexCompile(`^(?:[0-9]{4}-[0-9]{3}[0-9X])$`),
	"JWT":       lazyRegexCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`),
	"Latitude":  lazyRegexCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`),
	"Longitude": lazyRegexCompile(`^[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`),
	"Md4":       lazyRegexCompile(`^[0-9a-f]{32}$`),
	"Md5":       lazyRegexCompile(`^[0-9a-f]{32}$`),
	"MongodbConnString": lazyRegexCompile(
		`^mongodb(\+srv)?:\/\/(([a-zA-Z\d]+):([a-zA-Z\d$:\/?#\[\]@]+)@)?(([a-z\d.-]+)(:[\d]+)?)((,(([a-z\d.-]+)(:(\d+))?))*)?(\/([a-zA-Z-_]{1,64})?(\?(([a-zA-Z]+)=([a-zA-Z\d]+))(&(([a-zA-Z\d]+)=([a-zA-Z\d]+))?)*)?)?$`,
	),
	"MongodbId":      lazyRegexCompile(`^[a-f\d]{24}$`),
	"Multibyte":      lazyRegexCompile(`[^\x00-\x7F]`),
	"Numeric":        lazyRegexCompile(`^[-+]?[0-9]+(?:\.[0-9]+)?$`),
	"PrintableASCII": lazyRegexCompile(`^[\x20-\x7E]*$`),
	"RGB": lazyRegexCompile(
		`^rgb\x28\s*(?:(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)|(?:0|[1-9]\d*)\s*%\s*,\s*(?:0|[1-9]\d*)\s*%\s*,\s*(?:0|[1-9]\d*)\s*%)|(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*%\s*,\s*(?:(?:0.[1-9]*)|[01])\s*%\s*\)$`,
	),
	"RGBA": lazyRegexCompile(
		`^rgba\x28\s*(?:(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)|(	0|[1-9]\d*)\s*%\s*,\s*(?:0|[1-9]\d*)\s*%\s*,\s*(?:0|[1-9]\d*)\s*%)|(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*,\s*(?:0|[1-9]\d*)\s*%\s*,\s*(?:(?:0.[1-9]*)|[01])\s*%\s*\)$`,
	),
	"Ripemd128": lazyRegexCompile(`^[0-9a-f]{32}$`),
	"Ripemd160": lazyRegexCompile(`^[0-9a-f]{40}$`),
	"SSN": lazyRegexCompile(
		`^[0-9]{3}[ -]?(0[1-9]|[1-9][0-9])[ -]?([1-9][0-9]{3}|[0-9][1-9][0-9]{2}|[0-9]{2}[1-9][0-9]|[0-9]{3}[1-9])$`,
	),
	"Semver": lazyRegexCompile(
		`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
	), // numbered capture groups https://semver.org/
	"SHA256":       lazyRegexCompile(`^[0-9a-f]{64}$`),
	"SHA384":       lazyRegexCompile(`^[0-9a-f]{96}$`),
	"SHA512":       lazyRegexCompile(`^[0-9a-f]{128}$`),
	"SpicedbID":    lazyRegexCompile(`^(([a-zA-Z0-9/_|\-=+]{1,})|\*)$`),
	"ULID":         lazyRegexCompile(`^(?i)[A-HJKMNP-TV-Z0-9]{26}$`),
	"URLEncoded":   lazyRegexCompile(`^(?:[^%]|%[0-9A-Fa-f]{2})*$`),
	"UUID3RFC4122": lazyRegexCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-3[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
	"UUID3":        lazyRegexCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`),
	"UUID4RFC4122": lazyRegexCompile(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
	),
	"UUID4": lazyRegexCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`),
	"UUID5RFC4122": lazyRegexCompile(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-5[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
	),
	"UUID5":       lazyRegexCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`),
	"UUIDRFC4122": lazyRegexCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
	"UUID":        lazyRegexCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
}

// ValidateRegexp valida un `value` con una expresión regular indicada por su valor `name`, si `value` no es correcto,
// devuelve un error con el mensaje `errMsg`
func ValidateRegexp(name, value, errMsg string) error {
	if rexMap[name] == nil {
		return fmt.Errorf("validación regexp `%s` no es válido", name)
	}
	re := rexMap[name]()
	if !re.MatchString(value) {
		return errors.New(errMsg)
	}

	return nil
}
