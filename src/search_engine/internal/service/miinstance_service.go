package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"search_engine/config"
	"search_engine/internal/model"
	"strings"
)

type MiinstanceService struct {
	client  *http.Client
	api_url string
}

func NewMiinstanceService() *MiinstanceService {
	cfg := config.MustLoadConfig()
	return &MiinstanceService{
		client:  &http.Client{},
		api_url: cfg.MiinstanceServiceConfigs.MiinstanceServiceURL(),
	}
}

func (ms *MiinstanceService) GetMiinstances(ctx context.Context, passports []string) ([]*model.Miinstance, error) {
	if len(passports) == 0 {
		return []*model.Miinstance{}, nil
	}

	idxValue := strings.Join(passports, ",")

	base_url := ms.api_url
	fullURL := fmt.Sprintf("%s/api/miinstance/passport/?idx=%s", strings.TrimRight(base_url, "/"), idxValue)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("сервер вернул %d: %s", resp.StatusCode, string(body))
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела: %w", err)
	}

	// Пробуем разные варианты парсинга

	// Вариант 1: прямой массив
	var instances []*model.Miinstance
	if err := json.Unmarshal(body, &instances); err == nil {
		return instances, nil
	}

	// Вариант 2: объект с полем data
	var response1 struct {
		Data []*model.Miinstance `json:"data"`
	}
	if err := json.Unmarshal(body, &response1); err == nil {
		return response1.Data, nil
	}

	// Вариант 3: объект с полем result
	var response2 struct {
		Result []*model.Miinstance `json:"result"`
	}
	if err := json.Unmarshal(body, &response2); err == nil {
		return response2.Result, nil
	}

	var response3 struct {
		Items []*model.Miinstance `json:"items"`
	}
	if err := json.Unmarshal(body, &response3); err == nil {
		return response3.Items, nil
	}

	// Если ничего не подошло, возвращаем ошибку с сырым ответом
	return nil, fmt.Errorf("не удалось распарсить ответ сервера. Ожидался массив или объект с полями data/result/items. Получено: %s", string(body))
}
