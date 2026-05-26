package configs

import (
	"strings"
	"testing"
	"time"
)

func TestNewJwtConf(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantErrSub string
	}{
		{name: "missing secret", env: map[string]string{JWT_ACCESS_EXPIRY_MINUTES: "60"}, wantErrSub: "failed to get jwt secret"},
		{name: "missing expiry", env: map[string]string{JWT_SECRET_KEY: "secret"}, wantErrSub: "failed to get jwt accessExp"},
		{name: "bad expiry", env: map[string]string{JWT_SECRET_KEY: "secret", JWT_ACCESS_EXPIRY_MINUTES: "bad"}, wantErrSub: "invalid JWT_SECRET_KEY"},
		{name: "success", env: map[string]string{JWT_SECRET_KEY: "secret", JWT_ACCESS_EXPIRY_MINUTES: "60"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, JWT_SECRET_KEY, JWT_ACCESS_EXPIRY_MINUTES)
			setEnv(t, tt.env)

			conf, err := NewJwtConf()
			assertErrContains(t, err, tt.wantErrSub)
			if tt.wantErrSub == "" && (conf.Secret() != "secret" || conf.AccessExp() != time.Hour) {
				t.Fatalf("conf = secret:%q exp:%s, want secret and 1h", conf.Secret(), conf.AccessExp())
			}
		})
	}
}

func TestNewRedisConf(t *testing.T) {
	base := map[string]string{
		REDIS_HOST:        "localhost",
		REDIS_PORT:        "6379",
		REDIS_PASSWORD:    "pass",
		REDIS_DB:          "1",
		REDIS_TTL_MINUTES: "10",
	}

	tests := []struct {
		name       string
		mutate     func(map[string]string)
		wantErrSub string
	}{
		{name: "missing host", mutate: func(env map[string]string) { delete(env, REDIS_HOST) }, wantErrSub: "failed to get redis host"},
		{name: "missing port", mutate: func(env map[string]string) { delete(env, REDIS_PORT) }, wantErrSub: "failed to get redis port"},
		{name: "missing password", mutate: func(env map[string]string) { delete(env, REDIS_PASSWORD) }, wantErrSub: "failed to get redis password"},
		{name: "bad db", mutate: func(env map[string]string) { env[REDIS_DB] = "bad" }, wantErrSub: "failed to parse redis db"},
		{name: "bad ttl", mutate: func(env map[string]string) { env[REDIS_TTL_MINUTES] = "bad" }, wantErrSub: "failed to parse redis ttl"},
		{name: "success"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB, REDIS_TTL_MINUTES)
			env := cloneEnv(base)
			if tt.mutate != nil {
				tt.mutate(env)
			}
			setEnv(t, env)

			conf, err := NewRedisConf()
			assertErrContains(t, err, tt.wantErrSub)
			if tt.wantErrSub == "" && (conf.Address() != "localhost:6379" || conf.DB() != 1 || conf.TTL() != 10*time.Minute) {
				t.Fatalf("unexpected redis conf: addr=%s db=%d ttl=%s", conf.Address(), conf.DB(), conf.TTL())
			}
		})
	}
}

func TestNewNatsConf(t *testing.T) {
	base := map[string]string{
		NATS_URL:            "nats://localhost:4222",
		NATS_STREAM:         "LINK_EVENTS",
		NATS_SUBJECT_PREFIX: "links",
	}

	tests := []struct {
		name       string
		mutate     func(map[string]string)
		wantErrSub string
	}{
		{name: "missing url and host", mutate: func(env map[string]string) { delete(env, NATS_URL) }, wantErrSub: "failed to get nats host"},
		{name: "missing url and port", mutate: func(env map[string]string) {
			delete(env, NATS_URL)
			env[NATS_HOST] = "localhost"
		}, wantErrSub: "failed to get nats port"},
		{name: "bad batch", mutate: func(env map[string]string) { env[NATS_CONSUMER_BATCH] = "bad" }, wantErrSub: "failed to get valid nats consumer batch"},
		{name: "zero batch", mutate: func(env map[string]string) { env[NATS_CONSUMER_BATCH] = "0" }, wantErrSub: "failed to get valid nats consumer batch"},
		{name: "bad wait", mutate: func(env map[string]string) { env[NATS_CONSUMER_WAIT] = "bad" }, wantErrSub: "failed to get valid nats consumer wait"},
		{name: "zero wait", mutate: func(env map[string]string) { env[NATS_CONSUMER_WAIT] = "0" }, wantErrSub: "failed to get valid nats consumer wait"},
		{name: "success"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, NATS_URL, NATS_HOST, NATS_PORT, NATS_USER, NATS_PASSWORD, NATS_STREAM, NATS_SUBJECT_PREFIX, NATS_CONSUMER_NAME, NATS_CONSUMER_BATCH, NATS_CONSUMER_WAIT)
			env := cloneEnv(base)
			if tt.mutate != nil {
				tt.mutate(env)
			}
			setEnv(t, env)

			conf, err := NewNatsConf()
			assertErrContains(t, err, tt.wantErrSub)
			if tt.wantErrSub == "" && (conf.URL() != base[NATS_URL] || conf.Stream() != base[NATS_STREAM] || conf.SubjectPrefix() != base[NATS_SUBJECT_PREFIX]) {
				t.Fatalf("unexpected nats conf: url=%s stream=%s subject=%s", conf.URL(), conf.Stream(), conf.SubjectPrefix())
			}
		})
	}
}

func TestNewPgConf(t *testing.T) {
	base := map[string]string{
		PG_AUTH_DSN:      "auth-dsn",
		PG_LINKS_DSN:     "links-dsn",
		PG_ANALYTICS_DSN: "analytics-dsn",
	}

	tests := []struct {
		name       string
		mutate     func(map[string]string)
		wantErrSub string
	}{
		{name: "missing auth dsn", mutate: func(env map[string]string) { delete(env, PG_AUTH_DSN) }, wantErrSub: "failed to get pg authDsn"},
		{name: "missing links dsn", mutate: func(env map[string]string) { delete(env, PG_LINKS_DSN) }, wantErrSub: "failed to get pg linksDsn"},
		{name: "missing analytics dsn", mutate: func(env map[string]string) { delete(env, PG_ANALYTICS_DSN) }, wantErrSub: "failed to get pg analyticsDsn"},
		{name: "success"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, PG_AUTH_DSN, PG_LINKS_DSN, PG_ANALYTICS_DSN)
			env := cloneEnv(base)
			if tt.mutate != nil {
				tt.mutate(env)
			}
			setEnv(t, env)

			conf, err := NewPgConf()
			assertErrContains(t, err, tt.wantErrSub)
			if tt.wantErrSub == "" && (conf.AuthDSN() != "auth-dsn" || conf.LinksDSN() != "links-dsn" || conf.AnalyticsDSN() != "analytics-dsn") {
				t.Fatalf("unexpected pg conf: %#v", conf)
			}
		})
	}
}

func TestNewHttpConfMissingEnv(t *testing.T) {
	base := map[string]string{
		HTTP_HOST:           "0.0.0.0",
		HTTP_AUTH_PORT:      "8083",
		HTTP_LINKS_PORT:     "8084",
		HTTP_ANALYTICS_PORT: "8085",
	}

	tests := []struct {
		name       string
		missing    string
		wantErrSub string
	}{
		{name: "missing host", missing: HTTP_HOST, wantErrSub: "failed to get http host"},
		{name: "missing auth port", missing: HTTP_AUTH_PORT, wantErrSub: "failed to get http authPort"},
		{name: "missing links port", missing: HTTP_LINKS_PORT, wantErrSub: "failed to get http linksPort"},
		{name: "missing analytics port", missing: HTTP_ANALYTICS_PORT, wantErrSub: "failed to get http analyticsPort"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, HTTP_HOST, HTTP_AUTH_PORT, HTTP_LINKS_PORT, HTTP_ANALYTICS_PORT)
			env := cloneEnv(base)
			delete(env, tt.missing)
			setEnv(t, env)

			_, err := NewHttpConf()
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func TestNewGRPCConfMissingEnv(t *testing.T) {
	base := map[string]string{
		GRPC_HOST:             "0.0.0.0",
		GRPC_AUTH_DOCKER_HOST: "url_shortener_auth",
		GRPC_AUTH_PORT:        "8080",
		GRPC_LINKS_PORT:       "8081",
		GRPC_ANALYTICS_PORT:   "8082",
	}

	tests := []struct {
		name       string
		missing    string
		wantErrSub string
	}{
		{name: "missing host", missing: GRPC_HOST, wantErrSub: "failed to get grpc host"},
		{name: "missing auth docker host", missing: GRPC_AUTH_DOCKER_HOST, wantErrSub: "failed to get auth docker host"},
		{name: "missing auth port", missing: GRPC_AUTH_PORT, wantErrSub: "failed to get grpc auth_port"},
		{name: "missing links port", missing: GRPC_LINKS_PORT, wantErrSub: "failed to get grpc links_port"},
		{name: "missing analytics port", missing: GRPC_ANALYTICS_PORT, wantErrSub: "failed to get grpc analytics_port"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, GRPC_HOST, GRPC_AUTH_DOCKER_HOST, GRPC_AUTH_PORT, GRPC_LINKS_PORT, GRPC_ANALYTICS_PORT)
			env := cloneEnv(base)
			delete(env, tt.missing)
			setEnv(t, env)

			_, err := NewGRPCConf()
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func TestNewMetricsConfMissingEnv(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantErrSub string
	}{
		{name: "missing host", env: map[string]string{METRICS_ANALYTICS_PORT: "9090"}, wantErrSub: "failed to get metrics host"},
		{name: "missing analytics port", env: map[string]string{METRICS_HOST: "0.0.0.0"}, wantErrSub: "failed to get metrics analytics port"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t, METRICS_HOST, METRICS_ANALYTICS_PORT)
			setEnv(t, tt.env)

			_, err := NewMetricsConf()
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func cloneEnv(env map[string]string) map[string]string {
	cp := make(map[string]string, len(env))
	for key, value := range env {
		cp[key] = value
	}
	return cp
}

func setEnv(t *testing.T, env map[string]string) {
	t.Helper()

	for key, value := range env {
		t.Setenv(key, value)
	}
}

func clearEnv(t *testing.T, keys ...string) {
	t.Helper()

	for _, key := range keys {
		t.Setenv(key, "")
	}
}

func assertErrContains(t *testing.T, err error, wantSub string) {
	t.Helper()

	if wantSub == "" {
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("err = nil, want containing %q", wantSub)
	}
	if !strings.Contains(err.Error(), wantSub) {
		t.Fatalf("err = %q, want containing %q", err.Error(), wantSub)
	}
}
