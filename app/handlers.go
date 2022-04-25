package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	ipstackLocationEndpointFmt = "%s%s?access_key=%s"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("healthcheck performed")
	w.WriteHeader(http.StatusOK)
}

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("readiness check performed")
	w.WriteHeader(http.StatusOK)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	ipaddress := r.URL.Query().Get("ip")
	if err := checkIP(ipaddress); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiResponse, err := getIPLocation(ipaddress)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	tmp, _ := json.Marshal(apiResponse)
	w.Write(tmp)
}

func getIPLocation(ipaddress string) (*IPGeoData, error) {
	// we check redis first
	cacheResp, err := getFromCache(ipaddress)
	if err == nil && cacheResp != nil {
		logrus.Infof("got cached location for IP: %s", ipaddress)
		return cacheResp, nil
	}

	locationURL := fmt.Sprintf(ipstackLocationEndpointFmt, defaultIPStackURL, ipaddress, os.Getenv("IPSTACK_API_KEY"))
	resp, err := http.Get(locationURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response obtained from ipstack api")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipstackData IPGeoData
	err = json.Unmarshal(body, &ipstackData)
	if err != nil {
		return nil, err
	}

	if ipstackData.Error == nil {
		//store in redis
		saveToCache(ipstackData)
		return &ipstackData, nil

	}

	return nil, errors.New(ipstackData.Error.Info)
}

func checkIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	return nil
}
