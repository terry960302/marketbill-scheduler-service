package test

import "os"

type EnvSetter struct {
	Env string
}

func NewEnvSetter(env string) *EnvSetter {
	return &EnvSetter{
		Env: env,
	}
}

func (e *EnvSetter) SetEnv() {
	os.Setenv("PORT", "8080")
	os.Setenv("API_KEY", "4DC6A10B4F5D43D5977F364FC0DFE81C")

	switch e.Env {
	case "local":
		os.Setenv("PROFILE", "local")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PW", "postgres")
		os.Setenv("DB_NET", "tcp")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_NAME", "marketbill-test")
	case "dev":
		os.Setenv("PROFILE", "dev")
		os.Setenv("DB_USER", "marketbill")
		os.Setenv("DB_PW", "marketbill1234!")
		os.Setenv("DB_NET", "tcp")
		os.Setenv("DB_HOST", "marketbill-db.ciegftzvpg1l.ap-northeast-2.rds.amazonaws.com")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_NAME", "dev-db")
	case "prod":
		os.Setenv("PROFILE", "prod")
		os.Setenv("DB_USER", "marketbill")
		os.Setenv("DB_PW", "marketbill1234!")
		os.Setenv("DB_NET", "tcp")
		os.Setenv("DB_HOST", "marketbill-db.ciegftzvpg1l.ap-northeast-2.rds.amazonaws.com")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_NAME", "prod-db")
	}

}
