package config

import (
  "os"
  "strconv"
  "strings"
)

type AmqpConfig struct {
  AMQP_EXCHANGE string;
  AMQP_USER string;
  AMQP_PASSWORD string;
  AMQP_QUEUE string;
  AMQP_HOST string;
}

type DirectusConfig struct {
  DIRECTUS_HOST string;
  ADMIN_API_KEY string;
}

type Config struct {
  Amqp AmqpConfig;
  Directus DirectusConfig;
}

func New() *Config {
  return &Config {
    Amqp: AmqpConfig{
      AMQP_EXCHANGE: getEnv("AMQP_EXCHANGE", ""),
      AMQP_USER: getEnv("AMQP_USER", ""),
      AMQP_PASSWORD: getEnv("AMQP_PASSWORD", ""),
      AMQP_QUEUE: getEnv("AMQP_QUEUE", ""),
      AMQP_HOST: getEnv("AMQP_HOST", ""),
    },
    Directus: DirectusConfig{
      DIRECTUS_HOST: getEnv("DIRECTUS_HOST", ""),
      ADMIN_API_KEY: getEnv("ADMIN_API_KEY", ""),
    },
  };
}

func getEnv(key string, defaultVal string) string {
  if value, exists := os.LookupEnv(key); exists {
return value
  }

  return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
  valueStr := getEnv(name, "")
  if value, err := strconv.Atoi(valueStr); err == nil {
return value
  }

  return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
  valStr := getEnv(name, "")
  if val, err := strconv.ParseBool(valStr); err == nil {
return val
  }

  return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
  valStr := getEnv(name, "")

  if valStr == "" {
return defaultVal
  }

  val := strings.Split(valStr, sep)

  return val
}