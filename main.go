package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type RyanairResp struct {
	TermsOfUse    string    `json:"termsOfUse"`
	Currency      string    `json:"currency"`
	CurrPrecision int       `json:"currPrecision"`
	RouteGroup    string    `json:"routeGroup"`
	TripType      string    `json:"tripType"`
	UpgradeType   string    `json:"upgradeType"`
	Trips         []Trips   `json:"trips"`
	ServerTimeUTC time.Time `json:"serverTimeUTC"`
}
type Fares struct {
	Type              string  `json:"type"`
	Amount            float64 `json:"amount"`
	Count             int     `json:"count"`
	HasDiscount       bool    `json:"hasDiscount"`
	PublishedFare     float64 `json:"publishedFare"`
	DiscountInPercent int     `json:"discountInPercent"`
	HasPromoDiscount  bool    `json:"hasPromoDiscount"`
	DiscountAmount    float64 `json:"discountAmount"`
	HasBogof          bool    `json:"hasBogof"`
}
type RegularFare struct {
	FareKey   string  `json:"fareKey"`
	FareClass string  `json:"fareClass"`
	Fares     []Fares `json:"fares"`
}
type Segments struct {
	SegmentNr    int         `json:"segmentNr"`
	Origin       string      `json:"origin"`
	Destination  string      `json:"destination"`
	FlightNumber string      `json:"flightNumber"`
	Time         []string    `json:"time"`
	TimeUTC      []time.Time `json:"timeUTC"`
	Duration     string      `json:"duration"`
}
type Flights struct {
	FaresLeft    int         `json:"faresLeft"`
	FlightKey    string      `json:"flightKey"`
	InfantsLeft  int         `json:"infantsLeft"`
	RegularFare  RegularFare `json:"regularFare"`
	OperatedBy   string      `json:"operatedBy"`
	Segments     []Segments  `json:"segments"`
	FlightNumber string      `json:"flightNumber"`
	Time         []string    `json:"time"`
	TimeUTC      []time.Time `json:"timeUTC"`
	Duration     string      `json:"duration"`
}
type Dates struct {
	DateOut string    `json:"dateOut"`
	Flights []Flights `json:"flights"`
}
type Trips struct {
	Origin          string  `json:"origin"`
	OriginName      string  `json:"originName"`
	Destination     string  `json:"destination"`
	DestinationName string  `json:"destinationName"`
	RouteGroup      string  `json:"routeGroup"`
	TripType        string  `json:"tripType"`
	UpgradeType     string  `json:"upgradeType"`
	Dates           []Dates `json:"dates"`
}

func doReq(buf *bytes.Buffer, origin, destination string, date time.Time) error {
	u, _ := url.Parse("https://www.ryanair.com/api/booking/v4/es-es/availability")

	q := u.Query()
	q.Set("ADT", "1")
	q.Set("TEEN", "0")
	q.Set("CHD", "0")
	q.Set("INF", "0")
	q.Set("Origin", string(origin))
	q.Set("Destination", string(destination))
	q.Set("promoCode", "")
	q.Set("IncludeConnectingFlights", "false")
	q.Set("DateOut", date.Format("2006-01-02"))
	q.Set("DateIn", "")
	q.Set("FlexDaysBeforeOut", "0")
	q.Set("FlexDaysOut", "0")
	q.Set("FlexDaysBeforeIn", "0")
	q.Set("FlexDaysIn", "0")
	q.Set("RoundTrip", "false")
	q.Set("ToUs", "AGREED")

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.6")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Client", "desktop")
	req.Header.Set("Client-Version", "3.112.0")
	req.Header.Set("Cookie", "fr-correlation-id=c0183468-3e00-498a-ac6d-ab00fca53f0f; rid=a15d9f49-2b7d-4c37-9a38-40b00e6f260e; RY_COOKIE_CONSENT=true; STORAGE_PREFERENCES={\"STRICTLY_NECESSARY\":true,\"PERFORMANCE\":false,\"FUNCTIONAL\":false,\"TARGETING\":false,\"SOCIAL_MEDIA\":false,\"PIXEL\":false,\"__VERSION\":3}; mkt=/es/es/; myRyanairID=; rid.sig=jyna6R42wntYgoTpqvxHMK7H+KyM6xLed+9I3KsvYZaVt7P36AL6zp9dGFPu5uVxaIiFpNXrszr+LfNCdY3IT3oCSYLeNv/ujtjsDqOzkY5JmUFsCdAEz3kpPbhCUwiAt/vQmRn3hREI7zVdoZEQRlIxg+JQgZr7xof7l5bqUoYk5r1E2GfQ5gCk3SrXFOOL1I22oV1G8pkY/xDePp76SRBA4UOpS4LENQLeQf2nVFfnQwxsbaxHgGE+gPjOc0ToHmSELvxsYVrnnJkatLMBrMYjh1XlLfW5ir3VuR1oTL1K8/SIDrpL+pc3IpRqe5BNahBHDK37E7OR1jDBZEFw7xMw8+5iSbyiGnZmnf1+C3KwyVlndf2IYAyQMVqB1ruwEPftpvDIiJU1j369zuEjfRWBpCA1c++9SrM/MNlrp95Hgv+EtpOgQEnfxiE1wTjq4D9ZS06Vpy9mmuiawIC5mTJK6k4hNBkc8YB0SmIMIbA9QgApjzcNkwuu+TOxdYjwC+j2sINLQ6qlT5V+swjIRauPmGsMIUGr260SnamKcBdIsX1U/OM7VYWbVaZ3U59qMo0S4L0bolEMtC7i0HuSrqkxG1jxuLCLLn0ImKO+vBV5U0k4R3vFhPDjiFijHYg4EM9GTmSpbu1QyZ151/Zx4WG4J+jEgkR94oShHkk0XouHvbG7euEGoLu1zF6iYCUuwFDlfcCXc4JFgha/dcaIr6VVrZO8kCeb919YpZANmWWexvZ4tY1zkvDPEyUpOaiFn6sscVVdpca4PXR4BewbunjgQqDMZfY48lec8E6vXUcp5gilcn/VM1sfb0j2hZGlP1RLuYcJIPBFiBqU2/qv0VujvS0/2WFMxgdzzopiQdM=; .AspNetCore.Session=CfDJ8Ag8I2KwOBpOjaDfGy4wfk1MX5BLyISMSGydfNn45W2%2F2IFIqxXNAyIqi%2BO0T%2F%2BeRbZQEQJJpTgYUNZHhIhwW7pxiyRWae31rNsZdmujbAIiD5qxi82AKIrg%2F1HKI7FaM4v8DiDb0wHwS20U6S2dRl%2Br%2FF4oAqfLJLYlDZvnszQf")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", "\"Brave\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ryanairResp RyanairResp
	err = json.NewDecoder(resp.Body).Decode(&ryanairResp)
	if err != nil {
		return err
	}

	for _, trip := range ryanairResp.Trips {
		for _, flightsOnDate := range trip.Dates {
			for _, flight := range flightsOnDate.Flights {
				flightKey := flight.FlightKey
				faresLeft := flight.FaresLeft

				for _, fare := range flight.RegularFare.Fares {
					fmt.Fprintf(buf, "%s (in price %.2f€ left %d)\n", flightKey, fare.Amount, faresLeft)
				}
			}
		}
	}

	return nil
}

func SendMessage(msg string, token, chatId string) error {
	msgEscaped := url.QueryEscape(msg)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, chatId, msgEscaped)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func calculate(origins []string, destinations []string, dates []time.Time) (string, error) {
	var buf bytes.Buffer

	for _, date := range dates {
		fmt.Fprintf(&buf, "%s:\n", date.Format("2006-01-02"))

		for _, origin := range origins {
			for _, destination := range destinations {
				err := doReq(&buf, origin, destination, date)
				if err != nil {
					return "", err
				}
			}
		}
	}

	return buf.String(), nil
}

type Config struct {
	TelegramToken  string      `yaml:"telegramToken"`
	TelegramChatId string      `yaml:"telegramChatId"`
	Dates          []time.Time `yaml:"dates"`
	Origins        []string    `yaml:"origins"`
	Destinations   []string    `yaml:"destinations"`
}

// priority:
// 1. command line
// 2. environment variable
// 3. config file
func getConfig() (Config, error) {
	configFile := flag.String("config-file", "config.yml", "Set the yaml with the configuration")
	flag.String("telegram-token", "", "Telegram bot token")
	flag.String("telegram-chat-id", "", "Telegram chat id to send info")
	flag.String("dates", "", "Comma-separated dates")
	flag.String("origins", "", "Comma-separated origins")
	flag.String("destinations", "", "Comma-separated destinations")
	flag.Parse()

	var config Config

	// set from config file variables
	data, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("Error reading the file: %v.\n", err)
	} else {
		err := yaml.Unmarshal(data, &config)
		if err != nil {
			return config, fmt.Errorf("error parsing yaml: %v", err)
		}
	}

	var dates []string

	// set from environment variables
	if val, exists := os.LookupEnv("TELEGRAM_TOKEN"); exists {
		config.TelegramToken = val
	}
	if val, exists := os.LookupEnv("TELEGRAM_CHAT_ID"); exists {
		config.TelegramChatId = val
	}
	if val, exists := os.LookupEnv("DATES"); exists {
		dates = strings.Split(val, ",")
	}
	if val, exists := os.LookupEnv("ORIGINS"); exists {
		config.Origins = strings.Split(val, ",")
	}
	if val, exists := os.LookupEnv("DESTINATIONS"); exists {
		config.Destinations = strings.Split(val, ",")
	}

	// set from cmd
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "telegram-token":
			config.TelegramToken = f.Value.String()
		case "telegram-chat-id":
			config.TelegramChatId = f.Value.String()
		case "dates":
			dates = strings.Split(f.Value.String(), ",")
		case "origins":
			config.Origins = strings.Split(f.Value.String(), ",")
		case "destinations":
			config.Destinations = strings.Split(f.Value.String(), ",")
		}
	})

	// parse dates
	if dates != nil {
		config.Dates = make([]time.Time, 0, len(data))

		for _, date := range dates {
			d, err := time.Parse("2006-01-02", date)
			if err != nil {
				return config, fmt.Errorf("error date %s bad formated: %v", date, err)
			}

			config.Dates = append(config.Dates, d)
		}
	}

	return config, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	msg, err := calculate(config.Origins, config.Destinations, config.Dates)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg)
	if config.TelegramToken != "" && config.TelegramChatId != "" {
		err = SendMessage(msg, config.TelegramToken, config.TelegramChatId)
		if err != nil {
			panic(err)
		}
	}
}
