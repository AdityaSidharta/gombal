package gombal

import "github.com/kelseyhightower/envconfig"

type Env struct {
	VerifyToken     string `required:"true"`
	PageAccessToken string `required:"true"`
	Port            string `required:"true"`
}

func LoadEnv() (Env, error) {
	e := Env{}
	err := envconfig.Process("", &e)
	if err != nil {
		return e, err
	}
	return e, nil
}
