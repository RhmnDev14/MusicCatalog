package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Host        string `mapstructure:"DB_HOST"`
	Port        string `mapstructure:"DB_PORT"`
	User        string `mapstructure:"DB_USER"`
	Password    string `mapstructure:"DB_PASSWORD"`
	Name        string `mapstructure:"DB_NAME"`
	LogMode     bool   `mapstructure:"DB_LOG_MODE"`
	MaxIdle     int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpen     int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLife     int    `mapstructure:"DB_MAX_LIFE_TIME"`
	MaxIdleTime int    `mapstructure:"DB_MAX_IDLE_TIME"`
}

type ApiConfig struct {
	ApiPort string `mapstructure:"API_PORT"`
}

type TokenConfig struct {
	IssuerName              string        `mapstructure:"TOKEN_ISSUE"`
	JwtSignatureStr         string        `mapstructure:"TOKEN_SECRET"`
	JwtExpiresStr           string        `mapstructure:"TOKEN_EXPIRE"`         // ubah jadi string
	RefreshTokenExpiresStr  string        `mapstructure:"REFRESH_TOKEN_EXPIRE"` // ubah jadi string
	JwtExpiresTime          time.Duration `mapstructure:"-"`                    // parsing manual
	RefreshTokenExpiresTime time.Duration `mapstructure:"-"`                    // parsing manual
	JwtSigningMethod        *jwt.SigningMethodHMAC
	JwtSignatureKey         []byte `mapstructure:"-"`
}

type Config struct {
	DBConfig    `mapstructure:",squash"`
	TokenConfig `mapstructure:",squash"`
	ApiConfig   `mapstructure:",squash"`
}

func (c *Config) validate() error {
	expire, err := time.ParseDuration(c.TokenConfig.JwtExpiresStr)
	fmt.Println("TOKEN_EXPIRE:", viper.GetString("TOKEN_EXPIRE"))
	if err != nil || expire <= 0 {
		return fmt.Errorf("JWT expiration time must be greater than zero")
	}
	c.TokenConfig.JwtExpiresTime = expire

	c.TokenConfig.JwtSignatureKey = []byte(c.TokenConfig.JwtSignatureStr)
	if len(c.TokenConfig.JwtSignatureKey) == 0 {
		return fmt.Errorf("JWT signature key is required")
	}

	c.TokenConfig = TokenConfig{
		IssuerName:              c.TokenConfig.IssuerName,
		JwtSignatureKey:         c.TokenConfig.JwtSignatureKey,
		JwtExpiresTime:          c.TokenConfig.JwtExpiresTime,
		RefreshTokenExpiresTime: c.TokenConfig.RefreshTokenExpiresTime,
		JwtSigningMethod:        c.TokenConfig.JwtSigningMethod,
	}
	return nil
}

func NewConfig() *Config {
	var cfg Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // current directory

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil
	}

	if err := cfg.validate(); err != nil {
		panic(err)
	}
	cfg.TokenConfig.JwtSigningMethod = jwt.SigningMethodHS256

	return &cfg
}
