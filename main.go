package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	ipv4URL  = "https://raw.githubusercontent.com/dibdot/DoH-IP-blocklists/master/doh-ipv4.txt"
	ipv6URL  = "https://raw.githubusercontent.com/dibdot/DoH-IP-blocklists/master/doh-ipv6.txt"
	cacheTTL = 1 * time.Hour
)

var (
	ipv4Cache   []string
	ipv6Cache   []string
	lastUpdated time.Time
)

func fetchAndFilter(url string, filter string) ([]string, error) {
	fmt.Println("Fetching", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(body), "\n")
	filtered := make([]string, 0)

	for _, line := range lines {
		if !strings.HasPrefix(line, filter) {
			filtered = append(filtered, line)
		}
	}

	return filtered, nil
}

func updateCache() {
	ipv4List, err := fetchAndFilter(ipv4URL, "127.0.0.1")
	if err != nil {
		fmt.Println("Error fetching IPv4 list:", err)
		return
	}
	ipv4Cache = ipv4List

	ipv6List, err := fetchAndFilter(ipv6URL, "::1")
	if err != nil {
		fmt.Println("Error fetching IPv6 list:", err)
		return
	}
	ipv6Cache = ipv6List

	lastUpdated = time.Now()
}

func serveIPv4List(w http.ResponseWriter, r *http.Request) {
	if time.Since(lastUpdated) > cacheTTL {
		updateCache()
	}

	for _, ip := range ipv4Cache {
		fmt.Fprintf(w, "%s\n", ip)
	}
}

func serveIPv6List(w http.ResponseWriter, r *http.Request) {
	if time.Since(lastUpdated) > cacheTTL {
		updateCache()
	}

	for _, ip := range ipv6Cache {
		fmt.Fprintf(w, "%s\n", ip)
	}
}

func main() {
	updateCache()

	http.HandleFunc("/ipv4", serveIPv4List)
	http.HandleFunc("/ipv6", serveIPv6List)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
