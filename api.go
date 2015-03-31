package qapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// processResponse takes the body of a response, and either returns
// the error code, or unmarshalls the JSON response and places it into the
// output parameter. This function closes the response body after reading it.
func processResponse(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return newQuestradeError(res, body)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}
