package config

import (
	"github.com/todennus/x/logging"
	"github.com/todennus/x/mime"
	gormlogger "gorm.io/gorm/logger"
)

type Variable struct {
	Server       ServerVariable       `envconfig:"server"`
	Postgres     PostgresVariable     `envconfig:"postgres"`
	Redis        RedisVariable        `envconfig:"redis"`
	Minio        MinioVariable        `envconfig:"minio"`
	OAuth2       OAuth2Variable       `envconfig:"oauth2"`
	OAuth2Client OAuth2ClientVariable `envconfig:"oauth2_client"`
	Session      SessionVariable      `envconfig:"session"`
	Service      ServiceVariable      `envconfig:"service"`
	File         FileVariable         `envconfig:"file"`
	User         UserVariable         `envconfig:"user"`
}

func DefaultVariable() Variable {
	return Variable{
		Server:   DefaultServerVariable(),
		Postgres: DefaultPostgresVariable(),
		Redis:    DefaultRedisVariable(),
		Minio:    DefaultMinioVariable(),
		OAuth2:   DefaultOAuth2Variable(),
		Session:  DefaultSessionVariable(),
		Service:  DefaultServiceVariable(),
		File:     DefaultFileVariable(),
		User:     DefaultUserVariable(),
	}
}

type ServerVariable struct {
	Host           string `envconfig:"host"`
	Port           int    `envconfig:"port"`
	NodeID         int    `envconfig:"nodeid"`
	LogLevel       int    `envconfig:"loglevel"`
	RequestTimeout int    `envconfig:"request_timeout"` // The timeout of each request (in millisecond).
}

func DefaultServerVariable() ServerVariable {
	return ServerVariable{
		Host:           "0.0.0.0",
		Port:           8080,
		NodeID:         0,
		LogLevel:       int(logging.LevelDebug),
		RequestTimeout: 3000, // 3s
	}
}

type PostgresVariable struct {
	LogLevel      int `envconfig:"loglevel"`
	RetryAttempts int `envconfig:"retry_attempts"`
	RetryInterval int `envconfig:"retry_interval"` // in second
}

func DefaultPostgresVariable() PostgresVariable {
	return PostgresVariable{
		LogLevel:      int(gormlogger.Warn),
		RetryAttempts: 3,
		RetryInterval: 1,
	}
}

type RedisVariable struct {
	Addr string `envconfig:"addr"`
	DB   int    `envconfig:"db"`
}

func DefaultRedisVariable() RedisVariable {
	return RedisVariable{
		Addr: "localhost:6379",
		DB:   0,
	}
}

type MinioVariable struct {
	Endpoint string `envconfig:"endpoint"`
}

func DefaultMinioVariable() MinioVariable {
	return MinioVariable{
		Endpoint: "localhost:9000",
	}
}

type OAuth2Variable struct {
	IdPLoginURL string `envconfig:"idp_login_url"`

	TokenIssuer string `envconfig:"token_issuer"`

	AccessTokenExpiration  int `envconfig:"access_token_expiration"`  // in second
	RefreshTokenExpiration int `envconfig:"refresh_token_expiration"` // in second
	IDTokenExpiration      int `envconfig:"id_token_expiration"`      // in second

	// AuthorizationCodeFlowExpiration is the duration during which the code
	// must be exchanged.
	AuthorizationCodeFlowExpiration int `envconfig:"authorization_code_flow_expiration"` // in second

	// AuthenticationCallbackExpiration is the duration during which the IdP
	// must send the result to /auth/callback. Otherwise, the user must return
	// to the Client App to authenticate again.
	AuthenticationCallbackExpiration int `envconfig:"authentication_callback_expiration"` // in second

	// SessionUpdateExpiration is the duration during which the user must be
	// redirected to /session/update to update their session. Otherwise, the
	// user will be redirected to the IdP login page again.
	SessionUpdateExpiration int `envconfig:"session_update_expiration"` // in second

	// ConsentSessionExpiration is the duration during which the user must be
	// redirected to /oauth2/authorize to responds to Client about the consent
	// result. Otherwise, the user will be redirected to the consent page again.
	ConsentSessionExpiration int `envconfig:"consent_session_expiration"` // in second

	// ConsentExpiration is the duration during which every request to
	// /oauth2/authorize will be automatically consented to by the user. After
	// this period, the user will be redirected to the consent page again.
	ConsentExpiration int `envconfig:"consent_expiration"` // in second
}

func DefaultOAuth2Variable() OAuth2Variable {
	return OAuth2Variable{
		IdPLoginURL:                      "http://localhost:7063/login",
		TokenIssuer:                      "todennus",
		AccessTokenExpiration:            60,                // 60s
		RefreshTokenExpiration:           60 * 60,           // 1h
		IDTokenExpiration:                24 * 60 * 60,      // 1d
		AuthorizationCodeFlowExpiration:  10 * 60,           // 10m
		AuthenticationCallbackExpiration: 15 * 60,           // 15m
		SessionUpdateExpiration:          15,                // 15s
		ConsentSessionExpiration:         15,                // 15s
		ConsentExpiration:                30 * 24 * 60 * 60, // 30d
	}
}

type OAuth2ClientVariable struct {
	SecretLength int `envconfig:"secret_length"`
}

func DefaultOAuth2ClientVariable() OAuth2ClientVariable {
	return OAuth2ClientVariable{
		SecretLength: 64,
	}
}

type SessionVariable struct {
	Expiration int `envconfig:"expiration"`
}

func DefaultSessionVariable() SessionVariable {
	return SessionVariable{
		Expiration: 24 * 60 * 60, // 24h
	}
}

type ServiceVariable struct {
	OAuth2TokenURL       string `envconfig:"oauth2_token_url"`
	UserGRPCAddr         string `envconfig:"user_grpc_addr"`
	OAuth2ClientGRPCAddr string `envconfig:"oauth2_client_grpc_addr"`
	FileGRPCAddr         string `envconfig:"file_grpc_addr"`
}

func DefaultServiceVariable() ServiceVariable {
	return ServiceVariable{
		OAuth2TokenURL:       "http://localhost:8080/oauth2/token",
		UserGRPCAddr:         "localhost:8081",
		OAuth2ClientGRPCAddr: "localhost:8082",
		FileGRPCAddr:         "localhost:8085",
	}
}

type FileVariable struct {
	DefaultImageAllowedTypes []string `envconfig:"default_image_allowed_types"`
	DefaultMaxSize           int      `envconfig:"default_max_size"`

	// UploadSessionExpiration is the duration during which the user can use the
	// upload_token to upload a file.
	UploadSessionExpiration int `envconfig:"upload_session_expiration"`

	// TemporaryFileExpiration is the duration during which the user can use the
	// session_token to performe a specific action on the uploaded file like
	// setting an avatar, sending an image to a chat room, etc.
	TemporaryFileExpiration int `envconfig:"temporary_file_expiration"`

	StorageImageBucket     string `envconfig:"storage_image_bucket"`
	StorageTemporaryBucket string `envconfig:"storage_temporary_bucket"`
}

func DefaultFileVariable() FileVariable {
	return FileVariable{
		DefaultImageAllowedTypes: []string{mime.ImageJPEG, mime.ImagePNG}, // support png and jpeg.
		DefaultMaxSize:           3 * 1024 * 1024,                         // 1MB
		UploadSessionExpiration:  60,                                      // 1m
		TemporaryFileExpiration:  10 * 60,                                 // 10m
		StorageImageBucket:       "images",
		StorageTemporaryBucket:   "temporary-files",
	}
}

type UserVariable struct {
	AvatarAllowedTypes []string `envconfig:"avatar_allowed_types"`
	AvatarMaxSize      int      `envconfig:"avatar_max_size"`

	// AvatarPolicyTokenExpiration is the duration during which the user can
	// request the file-service validating the policy_token along with file
	// metadata to obtain an upload_token.
	AvatarPolicyTokenExpiration int `envconfig:"avatar_policy_token_expiration"`
}

func DefaultUserVariable() UserVariable {
	return UserVariable{
		AvatarAllowedTypes:          []string{mime.ImageJPEG, mime.ImagePNG}, // support png and jpeg.
		AvatarMaxSize:               3 * 1024 * 1024,                         // 3MB
		AvatarPolicyTokenExpiration: 60,                                      // 1m
	}
}
