package common

type Config struct {
	Environment  *uint8    `yaml:"environment" validate:"gte=1,lte=2"`
	Address      *string   `yaml:"address" validate:"required"`
	Cors         []*string `yaml:"cors" validate:"required"`
	MySqlDsn     *string   `yaml:"mysql_dsn" validate:"required"`
	MySqlMigrate *bool     `yaml:"mysql_migrate" validate:"required"`
}
