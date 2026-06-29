package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	Reset  = "\033[0m"
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Bold   = "\033[1m"
)

func printHelp() {
	fmt.Println(Cyan + Bold + "\n=== tg-keys CLI ===" + Reset)
	fmt.Println("Extract Telegram App ID & App Hash.")
	fmt.Println("\n" + Yellow + "Usage:" + Reset)
	fmt.Println("  tg-keys <phone_number>   Start extraction")
	fmt.Println("  tg-keys help             Show help message")
	fmt.Println("  tg-keys info             Show developer info")
	fmt.Println("\n" + Yellow + "Example:" + Reset)
	fmt.Println("  tg-keys +1234567890\n")
}

func printInfo() {
	fmt.Println(Cyan + Bold + "\n=== Developer Info ===" + Reset)
	fmt.Println("Developer : " + Green + "#\U0001D674\U0001D681\U0001D681\U0001D67E\U0001D681" + Reset)
	fmt.Println("GitHub    : https://github.com/error-404-h")
	fmt.Println("Telegram  : https://t.me/error404_h")
	fmt.Println("Version   : 1.0.0\n")
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func sendAuth(phone string) string {
	data := url.Values{}
	data.Set("phone", phone)
	req, _ := http.NewRequest("POST", "https://my.telegram.org/auth/send_password", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(Red + "Connection error." + Reset)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if hash, ok := result["random_hash"].(string); ok {
		return hash
	}
	fmt.Println(Red + "Error sending code. Check phone number." + Reset)
	os.Exit(1)
	return ""
}

func login(phone, hash, code string) string {
	data := url.Values{}
	data.Set("phone", phone)
	data.Set("random_hash", hash)
	data.Set("password", code)
	req, _ := http.NewRequest("POST", "https://my.telegram.org/auth/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(Red + "Connection error." + Reset)
		os.Exit(1)
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "stel_token" {
			return cookie.Value
		}
	}
	fmt.Println(Red + "Invalid verification code." + Reset)
	os.Exit(1)
	return ""
}

func getApi(token string) (string, string) {
	req, _ := http.NewRequest("GET", "https://my.telegram.org/apps", nil)
	req.AddCookie(&http.Cookie{Name: "stel_token", Value: token})
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`(?s)span class="form-control input-xlarge uneditable-input".*?>(.*?)</span>`)
	matches := re.FindAllStringSubmatch(html, -1)

	if len(matches) >= 2 {
		appID := strings.TrimSpace(strings.ReplaceAll(matches[0][1], "<strong>", ""))
		appID = strings.ReplaceAll(appID, "</strong>", "")
		appHash := strings.TrimSpace(matches[1][1])
		return appID, appHash
	}
	return "", ""
}

func start(phone string) {
	fmt.Printf(Cyan+"[*] Extracting for %s...\n"+Reset, phone)
	hash := sendAuth(phone)
	fmt.Println(Green + "[+] Verification code sent." + Reset)
	code := getInput(Yellow + "[?] Enter the 5-digit code: " + Reset)
	fmt.Printf(Cyan+"[*] Verifying...\n"+Reset)
	token := login(phone, hash, code)
	id, hashAPI := getApi(token)
	if id != "" && hashAPI != "" {
		fmt.Println(Green + "\n✔ Extraction Successful!" + Reset)
		fmt.Println("-------------------------------------------------")
		fmt.Printf(Bold+" App ID   : "+Reset+"%s\n", id)
		fmt.Printf(Bold+" App Hash : "+Reset+"%s\n", hashAPI)
		fmt.Println("-------------------------------------------------\n")
	} else {
		fmt.Println(Red + "\n✖ Failed to extract. Account may not have an app created yet." + Reset)
	}
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "help", "--help", "-h":
		printHelp()
	case "info", "about", "--info":
		printInfo()
	default:
		if strings.HasPrefix(cmd, "+") || isNumeric(cmd) {
			start(cmd)
		} else {
			printHelp()
		}
	}
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
