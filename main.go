package main

import (
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
)

const (
    username = "ut"
    password = "1Xbq59yFD6toT5Y3HSLPU8kB4R88c95JHnKw0kpN3cxbML5VGSwTSiOqz6qEZuFH"
)

func handler(w http.ResponseWriter, r *http.Request) {
    auth := r.Header.Get("Authorization")
    if auth == "" || !checkBasicAuth(auth) {
        w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    data, err := ioutil.ReadFile("/temp")
    if err != nil {
        http.Error(w, "Failed to read file", http.StatusInternalServerError)
        return
    }

    valueStr := strings.TrimSpace(string(data))
    value, err := strconv.ParseFloat(valueStr, 64)
    if err != nil {
        http.Error(w, "Invalid number", http.StatusInternalServerError)
        return
    }

    result := value / 1000

    expectHeader := strings.TrimSpace(r.Header.Get("X-Temp-Expect"))

    fail := false
    if expectHeader != "" {
        fail, err = evaluateExpectation(expectHeader, result)
        if err != nil {
            http.Error(w, "Invalid Expect header", http.StatusBadRequest)
            return
        }
    }

    w.Header().Set("Content-Type", "text/plain")

    if fail {
        w.WriteHeader(417)
    } else {
        w.WriteHeader(200)
    }

    // this will always return the temperature
    fmt.Fprintf(w, "Temp: %.2f°C\n", result)
}

func evaluateExpectation(expr string, value float64) (bool, error) {
    expr = strings.TrimSpace(expr)

    var op string
    var numStr string

    switch {
    case strings.HasPrefix(expr, "<="):
        op = "<="
        numStr = expr[2:]
    case strings.HasPrefix(expr, ">="):
        op = ">="
        numStr = expr[2:]
    case strings.HasPrefix(expr, "<"):
        op = "<"
        numStr = expr[1:]
    case strings.HasPrefix(expr, ">"):
        op = ">"
        numStr = expr[1:]
    default:
        return false, fmt.Errorf("invalid operator")
    }

    threshold, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64)
    if err != nil {
        return false, err
    }

    switch op {
    case "<":
        if value < threshold {
            return false, nil
        }
        return true, nil

    case "<=":
        if value <= threshold {
            return false, nil
        }
        return true, nil

    case ">":
        if value > threshold {
            return false, nil
        }
        return true, nil

    case ">=":
        if value >= threshold {
            return false, nil
        }
        return true, nil
    }

    return false, fmt.Errorf("unknown operator")
}

func checkBasicAuth(authHeader string) bool {
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || parts[0] != "Basic" {
        return false
    }

    decoded, err := base64.StdEncoding.DecodeString(parts[1])
    if err != nil {
        return false
    }

    creds := strings.SplitN(string(decoded), ":", 2)
    if len(creds) != 2 {
        return false
    }

    return creds[0] == username && creds[1] == password
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server listening on :8080")
    http.ListenAndServe(":8080", nil)
}
