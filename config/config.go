package config

import (
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/todennus/x/logging"
	"github.com/todennus/x/session"
	"github.com/todennus/x/token"
	"github.com/xybor-x/snowflake"
)

type Config struct {
	Secret   Secret
	Variable Variable

	Logger         logging.Logger
	TokenEngine    token.Engine
	SessionManager *session.Manager
}

func (c *Config) NewSnowflakeNode() *snowflake.Node {
	result, err := snowflake.NewNode(int64(c.Variable.Server.NodeID))
	if err != nil {
		panic(err)
	}
	return result
}

func Load(paths ...string) (*Config, error) {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			return nil, err
		}
	}

	c := &Config{}
	c.Variable = DefaultVariable()

	if err := load(&c.Variable); err != nil {
		return nil, err
	}

	if err := load(&c.Secret); err != nil {
		return nil, err
	}

	if err := c.loadInfras(); err != nil {
		return nil, err
	}

	return c, nil
}

func load[T any](obj *T) error {
	sType := reflect.TypeOf(obj).Elem()
	sValue := reflect.ValueOf(obj).Elem()
	for i := range sType.NumField() {
		field := sType.Field(i)
		prefix := field.Tag.Get("envconfig")
		if prefix == "" {
			prefix = strings.ToLower(field.Name)
		}

		if err := envconfig.Process(prefix, sValue.FieldByName(field.Name).Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadInfras() error {
	// Logger
	c.Logger = logging.NewSLogger(logging.Level(c.Variable.Server.LogLevel))

	// Token engine
	tokenEngine := token.NewJWTEngine()

	authSecrets := c.Secret.Authentication
	if err := tokenEngine.WithRSA(authSecrets.TokenRSAPrivateKey, authSecrets.TokenRSAPublicKey); err != nil {
		return err
	}

	if authSecrets.TokenHMACSecretKey != "" {
		if err := tokenEngine.WithHMAC(authSecrets.TokenHMACSecretKey); err != nil {
			return err
		}
	}

	c.TokenEngine = tokenEngine
	c.SessionManager = session.NewManager("/", c.Variable.Session.Expiration)

	return nil
}
