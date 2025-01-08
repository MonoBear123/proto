// Пакет для чтения конфигурационного файла и загрузки настроек.
package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Storage  Storage       `yaml:"storage"`
	TokenTTL time.Duration `yaml:"token_ttl"`
	GRPC     GRPCConfig    `yaml:"proto_gen"`
}
type Storage struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// MustLoadConfig - загружает конфигурацию из файла и возвращает структуру Config.
// При возникновении ошибки вызывается panic.
func MustLoadConfig() *Config {
	path := ReadPath()

	if path == "" {
		panic("config path is empty")
	}
	// Проверка существования конфигурационного файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file does not exist")
	}

	var config Config
	// Загрузка конфигурации из файла
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic(err)
	}

	return &config
}

// ReadPath - читает путь к конфигурационному файлу из аргументов командной строки или переменной окружения.
func ReadPath() string {
	var path string

	// Определение флага для указания пути к файлу конфигурации
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path

}
