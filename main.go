package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const installURL = "https://raw.githubusercontent.com/mhsanaei/3x-ui/master/install.sh"

var reader = bufio.NewReader(os.Stdin)

func input(prompt string) string {
	fmt.Print(prompt)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}
func confirm(s string) bool {
	v := strings.ToLower(input(s + " [д/Н]: "))
	return v == "д" || v == "да" || v == "y" || v == "yes"
}
func run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}
func main() {
	if os.Geteuid() != 0 { fmt.Println("Запусти от root: sudo ./goInstaller"); os.Exit(1) }
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" { fmt.Println("Поддерживаются amd64 и arm64"); os.Exit(1) }
	data, err := os.ReadFile("/etc/os-release")
	if err != nil || !strings.Contains(string(data), "ubuntu") || !strings.Contains(string(data), "24.04") {
		fmt.Println("Требуется Ubuntu 24.04 LTS"); os.Exit(1)
	}
	for {
		fmt.Println("\n=== goInstaller RU ===")
		fmt.Println("1) Установить 3x-ui (официальный скрипт)")
		fmt.Println("2) Статус службы x-ui")
		fmt.Println("3) Запустить / остановить / перезапустить x-ui")
		fmt.Println("4) Обновить 3x-ui")
		fmt.Println("5) Диагностика")
		fmt.Println("6) Сгенерировать UUID и REALITY keys")
		fmt.Println("0) Выход")
		switch input("Выбор: ") {
		case "1":
			if confirm("Загрузить и выполнить официальный install.sh 3x-ui?") {
				f, e := os.CreateTemp("", "3x-ui-*.sh")
				if e != nil { fmt.Println(e); continue }
				_ = f.Close()
				defer os.Remove(f.Name())
				if e = run("curl", "-fL", "--proto", "=https", "--tlsv1.2", "-o", f.Name(), installURL); e == nil { e = run("bash", f.Name()) }
				if e != nil { fmt.Println("Ошибка установки:", e) }
			}
		case "2":
			_ = run("systemctl", "status", "x-ui", "--no-pager")
		case "3":
			fmt.Println("1) start  2) stop  3) restart")
			switch input("Действие: ") {
			case "1": _ = run("systemctl", "start", "x-ui")
			case "2": if confirm("Остановить службу?") { _ = run("systemctl", "stop", "x-ui") }
			case "3": _ = run("systemctl", "restart", "x-ui")
			}
		case "4":
			if confirm("Открыть штатное меню обновления x-ui?") { _ = run("x-ui") }
		case "5":
			_ = run("bash", "-c", ". /etc/os-release; echo \"$PRETTY_NAME\"")
			_ = run("free", "-h"); _ = run("df", "-h", "/")
			_ = run("systemctl", "is-active", "x-ui"); _ = run("ss", "-lntp")
		case "6":
			fmt.Println("UUID:"); _ = run("uuidgen")
			fmt.Println("REALITY keypair:"); _ = run("xray", "x25519")
			fmt.Println("Создай inbound в панели: VLESS / TCP (RAW) / REALITY. Не передавай приватный ключ клиенту.")
		case "0": return
		default: fmt.Println("Неизвестный пункт")
		}
	}
}
