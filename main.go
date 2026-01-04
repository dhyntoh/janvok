package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	keySeed          = "zivpn-installer-key"
	defaultPlanMonth = "1-month"
	defaultPlanYear  = "1-year"
)

var encryptedPayloads = map[string]string{
	"zivpn":     "EUFUaf3Xnh9TdfWcoZSTDu3RCy+OHumDzrcAfQhv1c3q99AUtPPR0yi6dnCWTQItux1gt2RDdWWID4kRsRIxbFcm8vcMLCmW5ohW43sBwL7VGmlMKqks0+LyI5oJ/PDKyVLJiMSIxHieHGZPXnTTM9Zov6ohGupuvz/CCFMVqg3EStpxzXeVFroxIJDjpGyzFzgpPjRdCvfuYJUCRSpHh5RUkBHe8ICjWFjHD8IJW3oQRyc6Pa0LmPScy8TYCQIzGJIj3P9ZwaJb0VtR/4C4XhYwpqzZE2iEtlVPGhuIWW+2ldkmQvTQ7fTsuZQOrxHGFxhN7ohijuYzbf+B05K6gQefHXVP1WogWHrpUUENVLKoi3KbxifaP9J/wLnJyLTTG/gCr+U+4V8BYOfQeSPH0rmtk316D8RgpCnVa2UWChptUg1qOAAJHHiQNcvB968+jJ9Ncmq/mWmvEeJkChhogNhg+o0A7b8IuMbAxEIlw0GeBnRejUwXlu2N2Hm12t8sAldR/imRdygr11g+QvxlNDlGF/gF7Eb6txnes59d85e6SNTOnZc+L9kQQpS05Zq/WqXM7vLHUDe4qg3OfjCzKrnHWlkHfJMylRGZj8ixNZ4nKUZ2uua+WCzz1lFEtCxO4eZ6zrQMZEJOMXqns6LXsQr7aLWqBRE2HQj80qUPrGRWNb1ZpolCpf1qtgZcEAKRNoe1hbGc9l++gZS9R53ScoTXwqM5DwhBPWRey1uKEJZAkc0/slvkjU8uXClc4HvJRAdBOvvHeRjVqIG+qejfNJI9kKWRY2pBRRFTMFjvIOgcQcaukmWRUSnf2SdodBpTE7liTNF630/3UYerBS2FOHyBUUcCULpZXV5h/zOOicdyfKuM/XfLfDUG308NakjlncQAML1CzNvrX/1zGRDs4OTdMZpWQ3japbeS//ehGdvkNL9mCDE2WON0Ndxuk3vooWT21tst8ev+Fd3Q5qF9vJluMK4GWQI49IT0wv6ix3cYGaxz+cC2v2fr0V3DOVCfBkxnTqaAc9GwxxTyj4k3ZmvKTp2avWsi741pfVoHcVaegaV2CmWrwwkQfMaKQD5KpjMAdve1L/OigD5cg7dCO8k8+Cxl5m7yvMyp90q4L7vejYBu0J5qHuUjZnWR3u5TxHMvXb/P7NfPR+RqOxRc8apcHQwJ5jH8KWbZovWHS2Pp3QyAMw+enKHNFUggYeG1Q22/k4B0Uim+8Xca8RJNfX9UUBfE65S/6vT93l1nBTcP+j7bLslblXS3HE3maPj0wYUCLnJbTske2BZKCB0t8pi3l3gmtzSE9Vc/OcXXxGNj6C4int2Viygm5NCFVMRFIg6vWJfx+/4rhbOC64lt27/opWSdh+Dn6g+3qqj9hlKWifWOQinYeoMhnTEEa+5oIHQGPrrD4lMgcCxK0vcWvS+za30lMkhMI37O65SHRKoZw/R1VW6EwqMF3BJ0q6rOlBlMT8Id08M5zf8klB06f5kiKyzLI/k/ItFEVRGNV08Wp18AoSdY43bvnoPTcfBees31HuVsy+JezTI8Q+sIIMM9I0Cfs0DFLGjHYA2hXyMqtoEs1eSnlZO9geMEedCBQr2xY9ttHl8RbswhahZqCQO+TIE72GulKLONQeP2LTuj0K0cuvxLvsWIqspBRxSpkEkAV2nVoyeshB0V4VwyGQJrgvD+3jiK1jWuvIoBNnmClbLDje+PpQmWTTxUIm99EmuM82SUC7BrXVJkoNaHDpQ5czx5GBLJKafw8YcMebOoO+g7ZC/AIg==",
	"xray":      "ryww+1bzrCAYjMHW4m42DvdBgPqQvOX4NtMogWIkZKe0LvR2/IU/DmkIJDZZSExLStsep4D+nPf6T84ygaxpAwnYhfRN/8lyogkXO3knhnDVCPgshbaa/0VYbGBdI1HnNuRugNXDCw7+ZJhKujzF9PziXy5kFtU6zF3jItId5LixhwCFnbXGFQ+NpPU5hF/DXq+VSytHZid7/X7Y6daSJvebBGvWO1XrmvF83pYVm4Dxi7P7ItyugiThxN1TB/KxpfLsjvVXjCkcrfwvBPfI8p6L9xOUvL0AcKqAM3uQM0MA2cAWZFhaKHKg4P9ILGQv9vKKu2b081/NsgUmxSF8ZxEEf5WKowoBb1616gPKFxhoOdszUgQkQTLcn87/7TODSozvoG7Mld/99sekDCKydpWuy7u9wMRpIlLeti0NoJpx8nQxxUzOvex8cgUVwAzhMpyX9TNqbVwyl96enMd5VFVVSwUdNq3m7EkcQPltSZ1y/e7m1Rae9AlrH/4+AxrIrd++I7B9yzvgx9KASvclph8no7zyXVam2KeWf2ujOSCdHNowyig+TXPPhCfTNKT54s9ossPvmADqmDKybWwdMUiY7nR4gPYMfB3i0f+Hduc+1lS6m0gymn8dmeydBGqtobldxkUCyiAAKv4OpUb9RUQ4l58iTlq0yra9AgLPnTlzRPrLvYr//iufNbvA2TWNTAZmMN7Eq+jg9yjX+zT7MLhMkgtXkcwgSKgh2PcZCFVeRS6Kl3FckZepecQ4As/z7cH9goon6Ub41o1tQhHCZ1sxvXJa+kBiSE0AkwTnM32s1x6To/NVI0k6ybgdix5EfMIARLWEa10LG1qxYK5L4HncfqJkwt5PClnnzXV1isc1Bn0HjGn4W9+W93h8hx2L4YQfdPpCI9TCuTk8xIxBIKLebTPuch3Yu8C2mxQF3/xYjeZwoTB6FbHYH8ZwUDXdKnu2bVKy1+lXwYMV0o7KsLpuq4CPZdMXeu8lPedxNiVOpnODq+nFPptkPYMcaQUWboIluydcZTgisKdQMZ9Ps7940Avw9nWvXa162Tb7cTuu6/aK4DOHewC33dXBIZVhggsTFXCu9QB6ZDgOSGF7sRFeJI3jcl8rgmZrzXNDLg/UE44zj8hNTDyo0CYq80KmUbdjLUyIh4qF/3h7jEjOuCXn61jNFEEHaxG+ZBuNQgdqFxq67v4c8N8E/l55OxTSQ6j87psiflfJtSYg3xPLlVrfS/+0uLYrQss+TYJVhzyN1vbGm8d1In8QnGwzghzX28Q+a72uTHNraZrofrMhI0oscerF0ZuTHjYkwDYHuxigMGqPrIDVH4dAw59yTJx4so3BXim8IVJ/yj2WYBUPwpLKvCm2t8L/0XPO3qFbEHvIUu0tyK9LyEJEO+WsiQ5PY/0TFg4mP50vzGpjWWetAzY7iPf/rt9Qh5tzeJNbcAadIwbo0ufICsqS47YGcXieQWGkyIp7mGixocS4MCGVno1uuDC/FrtHiYPXgSN2jTew6jQkkM9790gGoiSOU+74LPswkUanHSBE+WAsj3lx2ILCndbd3R7B48znFUYZ9GywKVZY/+tmFpVr7az7sqv5HhWjKkIi54+DCZ1wCDw6cCs6r9WxYcTucRmkxzciR/suPQGhqOfOALW2WzKxAELZCuTzqD2YHnbvJzCWf1Gs7tOOCjHsz0iz",
	"hysteria2": "6uyb39P36rQkTse0Zoe0Af5ha48QVUMGp6212ucpByMd/xadFPy+6GeoKcZ+7m6p1yfl1gyTTdINCb4Mij2mBarWR6XRjhzwXAYNuvxlzOrhVquovSMgBCfCKVYD4CaM7U2GYpr3xl6sEgjTvjJm9a8Kr8gzClit827n5UyLpCirOA57awQ6iH2yYXy3BhWvy9AVINi5HQ9kS3lUIgY9BpPgdtprqyrl5f1vBtRdcxm1Oj/aKu3PgfirVNig+H835XuiKF1OhGaoqb9oE/OAM/MLE5XwxLbeUSOecbzCCfOpH2IS757zlFqjKQOw5K6T+aDXFIO9igCVC5vBZ/b7O4P07A3u4Tq8dJ24oRlYP97WCCnXjv/dv8UMvesHR8BFlIp1iB2JB+5NFPMRTz7oZfKZoyj2UcF11z7m3p8s/BPHvOZF7grVnWz1mkygOEEPQqgKCAJj/pqGYixlGjFZEFKgpMgfgkhfxKFlgbr8DfVXJ1f4RJpZ8dcfiO7U/R5oM0HZm//2TAUigesqanf0cv8R/vULD7x8O37Qc7hPu6uZmyX1a06e350KmXLtYqpOZ5a7vm2cyGCOATHu5FW7dMWQNvM5ke7/asxcMKQaz//TeD9Cxt22zDfFiFQZBdmVx3ODYhoGgedzW/jgPD76zyYtsH3J1d/J8GSZmaT7Urx13E1vLE0IblAFZZ4ibZtVjhSOeWshUfLXlGTETe8jDNcM3GUZizQr5ge9GPJaZPewgiBT73TYz+GIaa/ykUicizyRR70492W0DE1Kz0Irgnw5LfS/ANo6vAP+n5nqQRZLn1JAJvUcQCeEmlijuEcYGFf4HoAw1M+i+gt4qDMm3XSwCBILtzC7nU2oyKq8ybhjXPFG2/Pf55i/w/5Ycmr/XMRmJ7wpqmFjEeqxjjwCfCXadxZ6F6VJ0vtkGJAQ5klrhqw7v9FHwTdp7UE9JxkOEAN6ArEVSc+bWaqk/T08bv6ds7e7EFf191EgABfxYMVaAflxMvFF7oMJd6NLgXevMdUl3RPvQQWyAUWLRMzALFThHtzstlXJRoC09pezpbrxVuhZDwM0HoJXryw80BpRv6cOlhVmjdrD/blGJwNL+f4UjI4pz/hoY0vgva72WSHKLeUnXDbFSXM0PqqAo77jeSczwJToY3O+x9i3qIIYWWoaPQGNw18FKpmqRhW4Lys4JIr5kf/av8NBwLug3mUxU55WrF+D4iyE6saJJfSwiNled3+T62UAJTAnO9BNWoSkJ77ViQZRgkPTcCSZbp+T7Rp2xk5x7Y5PLH0rZZMWolksQEqq+APZXQ7wMUmdHiYt4paEcWdql/asCK0632ixYYg86COXe+XJaNx7HRxXrIAL1bIn7AeoErAERbUNjsvbk78n5tcqDHl4fLfHyqMDGk1utN+umoCGi/Ne5K4fmiHqK9M3gjpwXVtrOa53mQI6PskFJejD03d0haRiQ5o=",
}

type TokenStore struct {
	Tokens map[string]TokenEntry `json:"tokens"`
}

type TokenEntry struct {
	Plan      string    `json:"plan"`
	ExpiresAt time.Time `json:"expires_at"`
	UsedBy    string    `json:"used_by,omitempty"`
	UsedAt    time.Time `json:"used_at,omitempty"`
}

type AccountStore struct {
	Accounts map[string]AccountEntry `json:"accounts"`
}

type AccountEntry struct {
	Username    string    `json:"username"`
	ExpiresAt   time.Time `json:"expires_at"`
	DeviceLimit int       `json:"device_limit"`
}

type PortAssignments struct {
	ZivpnPort    int
	VmessPort    int
	VlessPort    int
	TrojanPort   int
	HysteriaPort int
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "mint":
		if err := mintToken(os.Args[2:]); err != nil {
			exitError(err)
		}
	case "bot":
		if err := runTelegramBot(os.Args[2:]); err != nil {
			exitError(err)
		}
	case "install":
		if err := runInstall(os.Args[2:]); err != nil {
			exitError(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: zivpn-installer <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  mint     Create activation token (1 month or 1 year)")
	fmt.Println("  bot      Run Telegram admin bot for account management")
	fmt.Println("  install  Run installer with activation token")
}

func mintToken(args []string) error {
	flags := flag.NewFlagSet("mint", flag.ContinueOnError)
	months := flags.Int("months", 1, "Validity period in months")
	years := flags.Int("years", 0, "Validity period in years")
	plan := flags.String("plan", "", "Plan label")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if *years > 0 {
		*months = *years * 12
	}
	if *months <= 0 {
		return errors.New("validity period must be positive")
	}

	store, path, err := loadTokenStore()
	if err != nil {
		return err
	}

	token, err := randomToken(24)
	if err != nil {
		return err
	}

	if *plan == "" {
		if *months >= 12 {
			*plan = defaultPlanYear
		} else {
			*plan = defaultPlanMonth
		}
	}

	store.Tokens[token] = TokenEntry{
		Plan:      *plan,
		ExpiresAt: time.Now().AddDate(0, *months, 0),
	}

	if err := saveTokenStore(path, store); err != nil {
		return err
	}

	fmt.Printf("Activation token: %s\n", token)
	notifyTelegram(fmt.Sprintf("Token dibuat (%s), berlaku sampai %s", *plan, store.Tokens[token].ExpiresAt.Format(time.RFC3339)))
	return nil
}

func runInstall(args []string) error {
	flags := flag.NewFlagSet("install", flag.ContinueOnError)
	dryRun := flags.Bool("dry-run", false, "Print steps without executing")
	if err := flags.Parse(args); err != nil {
		return err
	}

	store, path, err := loadTokenStore()
	if err != nil {
		return err
	}

	token := strings.TrimSpace(os.Getenv("ACTIVATION_TOKEN"))
	if token == "" {
		fmt.Print("Masukkan token aktivasi: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		token = strings.TrimSpace(input)
	}
	if token == "" {
		return errors.New("token aktivasi wajib diisi")
	}

	entry, ok := store.Tokens[token]
	if !ok {
		return errors.New("token tidak ditemukan")
	}
	if time.Now().After(entry.ExpiresAt) {
		return errors.New("token sudah kedaluwarsa")
	}

	machineID := getMachineID()
	if entry.UsedBy != "" && entry.UsedBy != machineID {
		return errors.New("token sudah digunakan di server lain")
	}

	entry.UsedBy = machineID
	entry.UsedAt = time.Now()
	store.Tokens[token] = entry
	if err := saveTokenStore(path, store); err != nil {
		return err
	}

	notifyTelegram(fmt.Sprintf("Token aktif di %s (plan %s, berlaku sampai %s)", machineID, entry.Plan, entry.ExpiresAt.Format(time.RFC3339)))

	ports, err := assignPorts()
	if err != nil {
		return err
	}

	installEnv := []string{
		fmt.Sprintf("ZIVPN_PORT=%d", ports.ZivpnPort),
		fmt.Sprintf("ZIVPN_PASSWORD=%s", randomPassword(12)),
		fmt.Sprintf("XRAY_VMESS_PORT=%d", ports.VmessPort),
		fmt.Sprintf("XRAY_VLESS_PORT=%d", ports.VlessPort),
		fmt.Sprintf("XRAY_TROJAN_PORT=%d", ports.TrojanPort),
		fmt.Sprintf("XRAY_UUID=%s", newUUID()),
		fmt.Sprintf("XRAY_TROJAN_PASSWORD=%s", randomPassword(14)),
		fmt.Sprintf("HYSTERIA_PORT=%d", ports.HysteriaPort),
		fmt.Sprintf("HYSTERIA_PASSWORD=%s", randomPassword(14)),
	}

	if *dryRun {
		fmt.Println("Dry run: lingkungan instalasi")
		for _, env := range installEnv {
			fmt.Println(env)
		}
		return nil
	}

	if err := runPayload("zivpn", installEnv); err != nil {
		return err
	}
	if err := runPayload("xray", installEnv); err != nil {
		return err
	}
	if err := runPayload("hysteria2", installEnv); err != nil {
		return err
	}

	fmt.Println("Instalasi selesai. Port yang digunakan:")
	fmt.Printf("  ZIVPN UDP: %d\n", ports.ZivpnPort)
	fmt.Printf("  Xray VMess: %d\n", ports.VmessPort)
	fmt.Printf("  Xray VLESS: %d\n", ports.VlessPort)
	fmt.Printf("  Xray Trojan: %d\n", ports.TrojanPort)
	fmt.Printf("  Hysteria2: %d\n", ports.HysteriaPort)
	return nil
}

func runPayload(name string, env []string) error {
	payload, err := decryptPayload(name)
	if err != nil {
		return err
	}

	tempDir := os.TempDir()
	scriptPath := filepath.Join(tempDir, fmt.Sprintf("%s-install.sh", name))
	if err := os.WriteFile(scriptPath, []byte(payload), 0o700); err != nil {
		return err
	}

	cmd := execCommand("bash", scriptPath)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func decryptPayload(name string) (string, error) {
	cipherText, ok := encryptedPayloads[name]
	if !ok {
		return "", fmt.Errorf("payload %s tidak ditemukan", name)
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	key := sha256.Sum256([]byte(keySeed))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext terlalu pendek")
	}

	nonce := data[:gcm.NonceSize()]
	cipherPayload := data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, cipherPayload, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func assignPorts() (PortAssignments, error) {
	used := map[int]bool{}
	pick := func(start, end int) (int, error) {
		for port := start; port <= end; port++ {
			if used[port] {
				continue
			}
			if isPortAvailable(port) {
				used[port] = true
				return port, nil
			}
		}
		return 0, fmt.Errorf("tidak ada port tersedia pada rentang %d-%d", start, end)
	}

	zivpnPort, err := pick(5667, 5699)
	if err != nil {
		return PortAssignments{}, err
	}
	vmessPort, err := pick(10000, 10050)
	if err != nil {
		return PortAssignments{}, err
	}
	vlessPort, err := pick(10051, 10100)
	if err != nil {
		return PortAssignments{}, err
	}
	trojanPort, err := pick(10101, 10150)
	if err != nil {
		return PortAssignments{}, err
	}
	hysteriaPort, err := pick(8443, 8499)
	if err != nil {
		return PortAssignments{}, err
	}

	return PortAssignments{
		ZivpnPort:    zivpnPort,
		VmessPort:    vmessPort,
		VlessPort:    vlessPort,
		TrojanPort:   trojanPort,
		HysteriaPort: hysteriaPort,
	}, nil
}

func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	packet, err := net.ListenPacket("udp", addr)
	if err != nil {
		return false
	}
	packet.Close()
	return true
}

func loadTokenStore() (TokenStore, string, error) {
	home := installerHome()
	if err := os.MkdirAll(home, 0o700); err != nil {
		return TokenStore{}, "", err
	}
	path := filepath.Join(home, "tokens.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TokenStore{Tokens: map[string]TokenEntry{}}, path, nil
		}
		return TokenStore{}, "", err
	}

	var store TokenStore
	if err := json.Unmarshal(data, &store); err != nil {
		return TokenStore{}, "", err
	}
	if store.Tokens == nil {
		store.Tokens = map[string]TokenEntry{}
	}
	return store, path, nil
}

func saveTokenStore(path string, store TokenStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func loadAccountStore() (AccountStore, string, error) {
	home := installerHome()
	if err := os.MkdirAll(home, 0o700); err != nil {
		return AccountStore{}, "", err
	}
	path := filepath.Join(home, "accounts.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AccountStore{Accounts: map[string]AccountEntry{}}, path, nil
		}
		return AccountStore{}, "", err
	}

	var store AccountStore
	if err := json.Unmarshal(data, &store); err != nil {
		return AccountStore{}, "", err
	}
	if store.Accounts == nil {
		store.Accounts = map[string]AccountEntry{}
	}
	return store, path, nil
}

func saveAccountStore(path string, store AccountStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func installerHome() string {
	if custom := os.Getenv("INSTALLER_HOME"); custom != "" {
		return custom
	}
	if os.Geteuid() == 0 {
		return "/etc/zivpn-installer"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".zivpn-installer")
}

func runTelegramBot(args []string) error {
	flags := flag.NewFlagSet("bot", flag.ContinueOnError)
	pollInterval := flags.Duration("poll-interval", 2*time.Second, "Polling interval for Telegram updates")
	if err := flags.Parse(args); err != nil {
		return err
	}

	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		return errors.New("TELEGRAM_BOT_TOKEN belum diisi")
	}
	adminIDs := parseAdminIDs(os.Getenv("TELEGRAM_ADMIN_IDS"))
	if len(adminIDs) == 0 {
		return errors.New("TELEGRAM_ADMIN_IDS belum diisi (pisahkan dengan koma)")
	}

	return pollTelegramUpdates(token, adminIDs, *pollInterval)
}

func parseAdminIDs(raw string) map[int64]bool {
	admins := map[int64]bool{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := parseInt64(part)
		if err != nil {
			continue
		}
		admins[id] = true
	}
	return admins
}

func parseInt64(input string) (int64, error) {
	var value int64
	_, err := fmt.Sscanf(input, "%d", &value)
	return value, err
}

type telegramUpdateResponse struct {
	OK     bool             `json:"ok"`
	Result []telegramUpdate `json:"result"`
}

type telegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *telegramMessage `json:"message"`
}

type telegramMessage struct {
	MessageID int64         `json:"message_id"`
	From      *telegramUser `json:"from"`
	Chat      telegramChat  `json:"chat"`
	Text      string        `json:"text"`
	Date      int64         `json:"date"`
}

type telegramUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type telegramChat struct {
	ID int64 `json:"id"`
}

func pollTelegramUpdates(token string, admins map[int64]bool, interval time.Duration) error {
	var offset int64
	for {
		updates, err := fetchTelegramUpdates(token, offset)
		if err != nil {
			time.Sleep(interval)
			continue
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message == nil || update.Message.From == nil {
				continue
			}
			if !admins[update.Message.From.ID] {
				continue
			}
			if err := handleAdminCommand(token, update.Message); err != nil {
				sendTelegramMessage(token, update.Message.Chat.ID, adminErrorText(fmt.Sprintf("Error: %v", err)), &telegramMessageOptions{ParseMode: "HTML"})
			}
		}
		time.Sleep(interval)
	}
}

func fetchTelegramUpdates(token string, offset int64) ([]telegramUpdate, error) {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("timeout", "10")
	if offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = query.Encode()

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result telegramUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, errors.New("telegram API error")
	}
	return result.Result, nil
}

func handleAdminCommand(token string, message *telegramMessage) error {
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return nil
	}
	text = normalizeAdminShortcut(text)
	fields := strings.Fields(text)
	command := strings.ToLower(fields[0])

	switch command {
	case "/start", "/menu":
		return sendTelegramMessage(token, message.Chat.ID, adminWelcomeText(), &telegramMessageOptions{
			ParseMode:   "HTML",
			ReplyMarkup: adminMenuKeyboard(),
		})
	case "/help":
		return sendTelegramMessage(token, message.Chat.ID, adminHelpText(), &telegramMessageOptions{
			ParseMode:   "HTML",
			ReplyMarkup: adminMenuKeyboard(),
		})
	case "/create":
		if len(fields) < 3 {
			return sendTelegramMessage(token, message.Chat.ID, adminUsageText("Create"), &telegramMessageOptions{ParseMode: "HTML"})
		}
		months, err := parseInt(fields[2])
		if err != nil || months <= 0 {
			return sendTelegramMessage(token, message.Chat.ID, adminErrorText("Bulan tidak valid."), &telegramMessageOptions{ParseMode: "HTML"})
		}
		limit := 1
		if len(fields) >= 4 {
			limit, _ = parseInt(fields[3])
		}
		return createAccount(token, message.Chat.ID, fields[1], months, limit)
	case "/delete":
		if len(fields) < 2 {
			return sendTelegramMessage(token, message.Chat.ID, adminUsageText("Delete"), &telegramMessageOptions{ParseMode: "HTML"})
		}
		return deleteAccount(token, message.Chat.ID, fields[1])
	case "/setexp":
		if len(fields) < 3 {
			return sendTelegramMessage(token, message.Chat.ID, adminUsageText("SetExp"), &telegramMessageOptions{ParseMode: "HTML"})
		}
		return setAccountExpiry(token, message.Chat.ID, fields[1], fields[2])
	case "/setlimit":
		if len(fields) < 3 {
			return sendTelegramMessage(token, message.Chat.ID, adminUsageText("SetLimit"), &telegramMessageOptions{ParseMode: "HTML"})
		}
		limit, err := parseInt(fields[2])
		if err != nil || limit <= 0 {
			return sendTelegramMessage(token, message.Chat.ID, adminErrorText("Limit tidak valid."), &telegramMessageOptions{ParseMode: "HTML"})
		}
		return setAccountLimit(token, message.Chat.ID, fields[1], limit)
	case "/list":
		return listAccounts(token, message.Chat.ID)
	default:
		return sendTelegramMessage(token, message.Chat.ID, adminErrorText("Perintah tidak dikenal. Ketik /help untuk bantuan."), &telegramMessageOptions{ParseMode: "HTML"})
	}
}

func adminHelpText() string {
	lines := []string{
		"<b>Admin Panel</b>",
		"Kelola akun dengan perintah berikut:",
		"",
		"➕ <b>Create</b>  <code>/create username months [limit]</code>",
		"🗑️ <b>Delete</b>  <code>/delete username</code>",
		"📅 <b>Set Expiry</b>  <code>/setexp username YYYY-MM-DD</code>",
		"📱 <b>Set Limit</b>  <code>/setlimit username limit</code>",
		"🧾 <b>List</b>  <code>/list</code>",
		"",
		"Gunakan tombol menu untuk akses cepat.",
	}
	return strings.Join(lines, "\n")
}

func adminWelcomeText() string {
	lines := []string{
		"<b>ZIVPN Admin Console</b>",
		"Status: <b>Online</b>",
		"",
		"Silakan pilih menu di bawah atau ketik <code>/help</code>.",
	}
	return strings.Join(lines, "\n")
}

func adminUsageText(section string) string {
	switch section {
	case "Create":
		return "<b>Format Create</b>\n<code>/create username months [limit]</code>\nContoh: <code>/create user1 1 2</code>"
	case "Delete":
		return "<b>Format Delete</b>\n<code>/delete username</code>"
	case "SetExp":
		return "<b>Format Set Expiry</b>\n<code>/setexp username YYYY-MM-DD</code>\nContoh: <code>/setexp user1 2025-01-31</code>"
	case "SetLimit":
		return "<b>Format Set Limit</b>\n<code>/setlimit username limit</code>\nContoh: <code>/setlimit user1 2</code>"
	default:
		return ""
	}
}

func adminErrorText(message string) string {
	return fmt.Sprintf("⚠️ <b>Perhatian</b>\n%s", message)
}

func adminSuccessText(message string) string {
	return fmt.Sprintf("✅ <b>Berhasil</b>\n%s", message)
}

func adminMenuKeyboard() *telegramReplyMarkup {
	return &telegramReplyMarkup{
		Keyboard: [][]telegramKeyboardButton{
			{
				{Text: "➕ Create Account"},
				{Text: "🧾 List Accounts"},
			},
			{
				{Text: "📅 Set Expiry"},
				{Text: "📱 Set Device Limit"},
			},
			{
				{Text: "🗑️ Delete Account"},
				{Text: "ℹ️ Help"},
			},
		},
		ResizeKeyboard: true,
	}
}

func normalizeAdminShortcut(text string) string {
	shortcuts := map[string]string{
		"➕ create account":   "/create",
		"🧾 list accounts":    "/list",
		"📅 set expiry":       "/setexp",
		"📱 set device limit": "/setlimit",
		"🗑️ delete account":  "/delete",
		"ℹ️ help":            "/help",
		"start":              "/start",
		"menu":               "/menu",
		"/menu":              "/menu",
		"/start":             "/start",
	}
	key := strings.ToLower(strings.TrimSpace(text))
	if mapped, ok := shortcuts[key]; ok {
		return mapped
	}
	return text
}

func createAccount(token string, chatID int64, username string, months int, limit int) error {
	store, path, err := loadAccountStore()
	if err != nil {
		return err
	}
	entry := AccountEntry{
		Username:    username,
		ExpiresAt:   time.Now().AddDate(0, months, 0),
		DeviceLimit: limit,
	}
	store.Accounts[username] = entry
	if err := saveAccountStore(path, store); err != nil {
		return err
	}
	message := fmt.Sprintf("Akun <b>%s</b> dibuat.\nExp: <b>%s</b>\nLimit: <b>%d device</b>", username, entry.ExpiresAt.Format("2006-01-02"), limit)
	return sendTelegramMessage(token, chatID, adminSuccessText(message), &telegramMessageOptions{ParseMode: "HTML"})
}

func deleteAccount(token string, chatID int64, username string) error {
	store, path, err := loadAccountStore()
	if err != nil {
		return err
	}
	if _, ok := store.Accounts[username]; !ok {
		return sendTelegramMessage(token, chatID, adminErrorText("Akun tidak ditemukan."), &telegramMessageOptions{ParseMode: "HTML"})
	}
	delete(store.Accounts, username)
	if err := saveAccountStore(path, store); err != nil {
		return err
	}
	return sendTelegramMessage(token, chatID, adminSuccessText(fmt.Sprintf("Akun <b>%s</b> dihapus.", username)), &telegramMessageOptions{ParseMode: "HTML"})
}

func setAccountExpiry(token string, chatID int64, username string, date string) error {
	store, path, err := loadAccountStore()
	if err != nil {
		return err
	}
	entry, ok := store.Accounts[username]
	if !ok {
		return sendTelegramMessage(token, chatID, adminErrorText("Akun tidak ditemukan."), &telegramMessageOptions{ParseMode: "HTML"})
	}
	exp, err := time.Parse("2006-01-02", date)
	if err != nil {
		return sendTelegramMessage(token, chatID, adminErrorText("Format tanggal salah. Gunakan YYYY-MM-DD."), &telegramMessageOptions{ParseMode: "HTML"})
	}
	entry.ExpiresAt = exp
	store.Accounts[username] = entry
	if err := saveAccountStore(path, store); err != nil {
		return err
	}
	message := fmt.Sprintf("Exp <b>%s</b> diupdate ke <b>%s</b>.", username, exp.Format("2006-01-02"))
	return sendTelegramMessage(token, chatID, adminSuccessText(message), &telegramMessageOptions{ParseMode: "HTML"})
}

func setAccountLimit(token string, chatID int64, username string, limit int) error {
	store, path, err := loadAccountStore()
	if err != nil {
		return err
	}
	entry, ok := store.Accounts[username]
	if !ok {
		return sendTelegramMessage(token, chatID, adminErrorText("Akun tidak ditemukan."), &telegramMessageOptions{ParseMode: "HTML"})
	}
	entry.DeviceLimit = limit
	store.Accounts[username] = entry
	if err := saveAccountStore(path, store); err != nil {
		return err
	}
	message := fmt.Sprintf("Limit <b>%s</b> diupdate ke <b>%d device</b>.", username, limit)
	return sendTelegramMessage(token, chatID, adminSuccessText(message), &telegramMessageOptions{ParseMode: "HTML"})
}

func listAccounts(token string, chatID int64) error {
	store, _, err := loadAccountStore()
	if err != nil {
		return err
	}
	if len(store.Accounts) == 0 {
		return sendTelegramMessage(token, chatID, adminErrorText("Belum ada akun."), &telegramMessageOptions{ParseMode: "HTML"})
	}
	return sendTelegramMessage(token, chatID, formatAccountListHTML(store.Accounts), &telegramMessageOptions{ParseMode: "HTML"})
}

func parseInt(input string) (int, error) {
	var value int
	_, err := fmt.Sscanf(input, "%d", &value)
	return value, err
}

func randomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func randomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "zivpnpass"
	}
	return hex.EncodeToString(bytes)[:length]
}

func newUUID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

func getMachineID() string {
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(data))
	}
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func notifyTelegram(message string) {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	chatID := strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID"))
	if token == "" || chatID == "" {
		return
	}

	chatIDValue, err := parseInt64(chatID)
	if err != nil {
		return
	}
	_ = sendTelegramMessage(token, chatIDValue, message, nil)
}

type telegramMessageOptions struct {
	ParseMode   string               `json:"parse_mode,omitempty"`
	ReplyMarkup *telegramReplyMarkup `json:"reply_markup,omitempty"`
}

type telegramReplyMarkup struct {
	Keyboard        [][]telegramKeyboardButton `json:"keyboard,omitempty"`
	ResizeKeyboard  bool                       `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool                       `json:"one_time_keyboard,omitempty"`
	RemoveKeyboard  bool                       `json:"remove_keyboard,omitempty"`
}

type telegramKeyboardButton struct {
	Text string `json:"text"`
}

func sendTelegramMessage(token string, chatID int64, message string, opts *telegramMessageOptions) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := map[string]any{
		"chat_id": fmt.Sprintf("%d", chatID),
		"text":    message,
	}
	if opts != nil {
		if opts.ParseMode != "" {
			payload["parse_mode"] = opts.ParseMode
		}
		if opts.ReplyMarkup != nil {
			payload["reply_markup"] = opts.ReplyMarkup
		}
	}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func formatAccountListHTML(accounts map[string]AccountEntry) string {
	lines := []string{
		"<b>Daftar Akun</b>",
		"<pre>USERNAME           EXPIRY     LIMIT",
		"------------------------------------</pre>",
	}
	for _, entry := range accounts {
		lines = append(lines, fmt.Sprintf("<pre>%-18s %-10s %5d</pre>", entry.Username, entry.ExpiresAt.Format("2006-01-02"), entry.DeviceLimit))
	}
	lines = append(lines, "", "Gunakan <code>/help</code> untuk format perintah.")
	return strings.Join(lines, "\n")
}

func execCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
