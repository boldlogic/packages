package commonconfig

import "flag"

// GetConfigPath возвращает путь к файлу конфигурации.
// Если флаг -config не передан, возвращает defaultConfigPath.
// Флаг парсится из аргументов командной строки.
func GetConfigPath(defaultConfigPath string) string {
	if path := flag.Lookup("config"); path != nil {
		if !flag.Parsed() {
			flag.Parse()
		}
		if value := path.Value.String(); value != "" {
			return value
		}
		return defaultConfigPath
	}

	path := flag.String("config", defaultConfigPath, "")
	if !flag.Parsed() {
		flag.Parse()
	}
	return *path
}
