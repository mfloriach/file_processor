package repo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"processor/types"
	"processor/utils"
)

const (
	EndpointSearch = "/movie?search="
)

type metadataImdbApiRepo struct {
	client  utils.HttpClient
	token   string
	baseUrl string
}

func NewMetadataImdbApiRepo() types.MovieDetailsRepo {
	res := map[string]interface{}{
		"description": "from api",
		"stars":       1.4,
		"actors":      []string{"2019-03-05", "2019-03-05", "2019-03-05"},
	}
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	env := utils.GetEnv()

	url, err := url.Parse(env.ImdbBaseUrl)
	if err != nil {
		panic(err)
	}

	client := utils.NewHttpClient(utils.NewMockHttpClient(b, nil))

	return &metadataImdbApiRepo{
		client:  client,
		token:   env.ImdbToken,
		baseUrl: url.Host,
	}
}

func (mda metadataImdbApiRepo) Get(ctx context.Context, title string) (types.MovieDetails, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mda.baseUrl+EndpointSearch+title, nil)
	if err != nil {
		return types.MovieDetails{}, err
	}
	req.Header.Add("Authentication", mda.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := mda.client.Do(req)
	if err != nil {
		return types.MovieDetails{}, err
	}

	var m types.MovieDetails
	if err := json.Unmarshal(resp, &m); err != nil {
		return types.MovieDetails{}, err
	}

	return m, nil
}
