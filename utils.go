package main

import "net/url"

func validateURLs(urls []string) ([]url.URL, error) {
	r := make([]url.URL, 0)
	for _, u := range urls {
		addr, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		r = append(r, *addr)
	}
	return r, nil
}
