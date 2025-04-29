package enrich

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AgifyResponse struct {
	Name  string `json:"name"`
	Age   *int   `json:"age"`
	Count int    `json:"count"`
}

type GenderizeResponse struct {
	Name   string  `json:"name"`
	Gender *string `json:"gender"`
	Prob   float64 `json:"probability"`
	Count  int     `json:"count"`
}

type NationalizeResponse struct {
	Name    string `json:"name"`
	Country []struct {
		CountryID string  `json:"country_id"`
		Prob      float64 `json:"probability"`
	} `json:"country"`
}

var httpClient = &http.Client{Timeout: 3 * time.Second}

func GetAge(ctx context.Context, name string) (*int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error getting age from agify: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	var agifyResp AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&agifyResp); err != nil {
		log.Printf("Error decoding agify response: %v", err)
		return nil, err
	}
	return agifyResp.Age, nil
}

func GetGender(ctx context.Context, name string) (*string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error getting gender from genderize: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	var genderize GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderize); err != nil {
		log.Printf("Error decoding genderize response: %v", err)
		return nil, err
	}
	return genderize.Gender, nil
}

func GetNationality(ctx context.Context, name string) (*string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error getting nationality from nationalize: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	var nationalize NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalize); err != nil {
		log.Printf("Error decoding genderize response: %v", err)
		return nil, err
	}
	if len(nationalize.Country) > 0 {
		return &nationalize.Country[0].CountryID, nil
	}
	return nil, nil
}

type EnrichDate struct {
	Age         *int
	Gender      *string
	Nationality *string
}

func EnrichPerson(ctx context.Context, name string) (*EnrichDate, error) {
	type result struct {
		age         *int
		gender      *string
		nationality *string
		err         error
	}
	ch := make(chan result, 3)

	go func() {
		age, err := GetAge(ctx, name)
		ch <- result{age: age, err: err}
	}()
	go func() {
		gender, err := GetGender(ctx, name)
		ch <- result{gender: gender, err: err}
	}()
	go func() {
		nationality, err := GetNationality(ctx, name)
		ch <- result{nationality: nationality, err: err}
	}()

	var data EnrichDate
	for i := 0; i < 3; i++ {
		res := <-ch
		if res.err != nil {
			log.Printf("Error enriching person: %v", res.err)
		}
		if res.age != nil {
			data.Age = res.age
		}
		if res.gender != nil {
			data.Gender = res.gender
		}
		if res.nationality != nil {
			data.Nationality = res.nationality
		}
	}
	return &data, nil
}
