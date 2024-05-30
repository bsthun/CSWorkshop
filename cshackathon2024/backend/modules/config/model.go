package iconfig

type Config struct {
	//Environment uint8  `yaml:"environment" validate:"gte=1,lte=2"`
	LogLevel uint32 `yaml:"log_level" validate:"required"`

	Address string   `yaml:"address" validate:"required"`
	Cors    []string `yaml:"cors" validate:"required"`

	MysqlDsn string `yaml:"mysql_dsn" validate:"required"`

	JwtSecret string `yaml:"jwt_secret" validate:"required"`

	SpotifyAuthorization string `yaml:"spotify_authorization" validate:"required"`
	SpotifyRefreshToken  string `yaml:"spotify_refresh_token" validate:"required"`

	WeatherApiKey string `yaml:"weather_api_key" validate:"required"`
}
