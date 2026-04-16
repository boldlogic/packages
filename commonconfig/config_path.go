package commonconfig

import "flag"

// GetConfigPath возвращает путь к файлу конфигурации.
// Если флаг -config не передан, возвращает defaultConfigPath.
// Флаг парсится из аргументов командной строки.
func GetConfigPath(defaultConfigPath string) string {
	path := flag.String("config", defaultConfigPath, "")
	flag.Parse()
	return *path
}
