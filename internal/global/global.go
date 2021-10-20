package global

import "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/config/settings"

var (
	Env      string
	Settings *settings.Settings
)

func Init(env string) error {
	var err error

	// Set package level Env var
	Env = env

	// Read and parse setting file and set package level Settings var
	Settings, err = settings.Read(env, "configs/settings.yml")
	if err != nil {
		return err
	}

	return nil
}
