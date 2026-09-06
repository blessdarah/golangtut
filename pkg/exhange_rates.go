package pkg

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// {"result":"success","documentation":"https://www.exchangerate-api.com/docs","terms_of_use":"https://www.exchangerate-api.com/terms","time_last_update_unix":1788566402,"time_last_update_utc":"Sat, 05 Sep 2026 00:00:02 +0000","time_next_update_unix":1788652802,"time_next_update_utc":"Sun, 06 Sep 2026 00:00:02 +0000","base_code":"USD","target_code":"XAF","conversion_rate":564.8552}

type exchangeRateByCurrency struct {
	Result             string  `json:"result"`
	Documentation      string  `json:"documentation"`
	TermsOfUse         string  `json:"terms_of_use"`
	TimeLastUpdateUnix int     `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string  `json:"time_last_update_utc"`
	TimeNextUpdateUnix int     `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string  `json:"time_next_update_utc"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
}

type exchangeRateService struct {
	baseURL string
	logger  *slog.Logger
}

type ExchangeRateService interface {
	GetExchangeRateByPair(from, to string) (*exchangeRateByCurrency, error)
}

func NewExchangeRateService(baseURL string, logger *slog.Logger) *exchangeRateService {
	return &exchangeRateService{
		baseURL,
		logger,
	}
}

func (s *exchangeRateService) GetExchangeRateByPair(
	from, to string,
) (*exchangeRateByCurrency, error) {
	url := fmt.Sprintf("%s/pair/%s/%s", s.baseURL, from, to)
	s.logger.Info("get exchange rate", "from", from, "to", to, "url", url)
	data, err := http.Get(url)
	if err != nil {
		s.logger.Error("failed to get exchange rate", "error", err)
		return nil, err
	}
	defer data.Body.Close()

	s.logger.Info("exr service:", "res data", data)

	var res exchangeRateByCurrency
	err = json.NewDecoder(data.Body).Decode(&res)
	if err != nil {
		s.logger.Error("failed to decode exchange rate", "error", err)
		return nil, err
	}

	return &res, nil
}
