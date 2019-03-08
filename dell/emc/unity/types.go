package unity

import "time"

type MetricRealTimeQuery struct {
	Content struct {
		Id             int       `json:"id"`
		Paths          []string  `json:paths`
		Interval       int       `json:"interval"`
		MaximumSamples int       `json:"maximumSamples"`
		Expiration     time.Time `json:"expiration"`
	} `json:"content"`
}

type Metric struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:updated`
	Links   []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Entries []struct {
		Content struct {
			QueryId   int       `json:"queryId"`
			Path      string    `json:"path"`
			Timestamp time.Time `json:"timestamp"`
			Values    struct {
				Spa interface{} `json:"spa"`
				Spb interface{} `json:"spb"`
				// Spa string `json:"spa"`
				// Spb string `json:"spb"`
			} `json:"values"`
		} `json:"content"`
	} `json:"entries"`
}
