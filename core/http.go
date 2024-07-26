package core

import (
	"time"
	"net/http"
	"fmt"
    "crypto/tls"
)


func sendRequest(protocol, ip, port string, timeout time.Duration, userAgent string) (bool, string, string) {
    url := fmt.Sprintf("%s://%s:%s", protocol, ip, port)
    client := &http.Client{
        Timeout: timeout,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, "", ""
    }
    req.Header.Set("User-Agent", userAgent)
    
    resp, err := client.Do(req)
    if err != nil {
        return false, "", ""
    }
    defer resp.Body.Close()
    
    hostname := ""
    if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
        hostname = resp.TLS.PeerCertificates[0].Subject.CommonName
    }
    if hostname == "" {
        hostname = resp.Request.Host
    }
    
    return true, fmt.Sprintf("%d", resp.StatusCode), hostname
}
