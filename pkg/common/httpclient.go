package common

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

//MakePUTAPICall :
func MakePUTAPICall(server, uri string, payload []byte, payloadType string) (int, error) {
	connURL := fmt.Sprintf("%s%s", server, uri)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("PUT", connURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	retries := 3
	retcode := http.StatusOK
	for try := 1; try <= retries; try++ {
		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if try == retries {
				zap.S().Debugf("[try=%d] error making http call: %v, stopping tries ", try, err)
				return http.StatusInternalServerError,
					errors.New("error conencting to gitsecure api service")
			}
			zap.S().Debugf("[try=%d] error making http call: %v ", try, err)
			zap.S().Debugf("trying again after %d seconds", 5)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()
		retcode = resp.StatusCode
		break
	}
	return retcode, nil
}

//MakeGetAPICall :
func MakeGetAPICall(server, uri string, payload []byte) (int, []byte, error) {
	connURL := fmt.Sprintf("%s%s", server, uri)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("GET", connURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	retries := 3
	body := []byte{}
	retcode := http.StatusOK
	for try := 1; try <= retries; try++ {
		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if try == retries {
				zap.S().Debugf("[try=%d] error making http call: %v, stopping tries ", try, err)
				return http.StatusInternalServerError, nil,
					errors.New("error conencting to gitsecure api service")
			}
			zap.S().Debugf("[try=%d] error making http call: %v ", try, err)
			zap.S().Debugf("trying again after %d seconds", 5)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		retcode = resp.StatusCode
		break
	}

	return retcode, body, nil
}

//MakePOSTAPICall :
func MakePOSTAPICall(url string, payload []byte) (int, []byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	retries := 3
	body := []byte{}
	retcode := http.StatusOK
	for try := 1; try <= retries; try++ {
		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if try == retries {
				zap.S().Debugf("[try=%d] error making http call: %v, stopping tries ", try, err)
				return http.StatusInternalServerError, nil,
					errors.New("error conencting to gitsecure api service")
			}
			zap.S().Debugf("[try=%d] error making http call: %v ", try, err)
			zap.S().Debugf("trying again after %d seconds", 5)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		retcode = resp.StatusCode
		break
	}
	return retcode, body, nil
}

//MakeDeleteAPICall :
func MakeDeleteAPICall(server, uri string, payload []byte) (int, []byte, error) {
	connURL := fmt.Sprintf("%s%s", server, uri)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("DELETE", connURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	retries := 3
	body := []byte{}
	retcode := http.StatusOK
	for try := 1; try <= retries; try++ {
		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if try == retries {
				zap.S().Debugf("[try=%d] error making http call: %v, stopping tries ", try, err)
				return http.StatusInternalServerError, nil,
					errors.New("error conencting to gitsecure api service")
			}
			zap.S().Debugf("[try=%d] error making http call: %v ", try, err)
			zap.S().Debugf("trying again after %d seconds", 5)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		retcode = resp.StatusCode
		break
	}
	return retcode, body, nil
}
