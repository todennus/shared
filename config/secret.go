package config

type Secret struct {
	Postgres       PostgresSecret       `envconfig:"postgres"`
	Authentication AuthenticationSecret `envconfig:"auth"`
	OAuth2         OAuth2Secret         `envconfig:"oauth2"`
	Redis          RedisSecret          `envconfig:"redis"`
	Session        SessionSecret        `envconfig:"session"`
}

type PostgresSecret struct {
	DSN string `envconfig:"dsn"`
}

type AuthenticationSecret struct {
	// Using RSA key to sign and verify the token. If both RSAKey and SecretKey
	// are provided, RSAKey will be used.
	TokenRSAPrivateKey string `envconfig:"token_rsa_private_key"`
	TokenRSAPublicKey  string `envconfig:"token_rsa_public_key"`

	// Use HMAC to sign and verify the token. Not support verifying at client.
	TokenHMACSecretKey string `envconfig:"token_hmac_secret_key"`
}

type OAuth2Secret struct {
	IdPSecret string `envconfig:"idp_secret"`
}

type RedisSecret struct {
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
}

type SessionSecret struct {
	AuthenticationKey string `envconfig:"authentication_key"`
	EncryptionKey     string `envconfig:"encryption_key"`
}
