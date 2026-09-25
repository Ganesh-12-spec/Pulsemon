package checker

import "net/http"

func Check(url string) error {
	_, err := http.Get(url)

	if err != nil {
		return err
	}

	return nil
}
