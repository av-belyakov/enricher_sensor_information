package sensorinformationapi

import "errors"

//******************* функции настройки опций ***********************

// WithHost имя или ip адрес хоста API
func WithHost(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		sic.settings.host = v

		return nil
	}
}

// WithPort порт API
func WithPort(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		sic.settings.port = v

		return nil
	}
}

// WithUser имя пользователя
func WithUser(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'user' cannot be empty")
		}

		sic.settings.user = v

		return nil
	}
}

// WithPasswd пароль пользователя пользователя
func WithPasswd(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'passwd' cannot be empty")
		}

		sic.settings.passwd = v

		return nil
	}
}

// WithNCIRCCURL URL API НКЦКИ
func WithNCIRCCURL(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'ncirccURL' cannot be empty")
		}

		sic.settings.ncirccURL = v

		return nil
	}
}

// WithNCIRCCToken токен API НКЦКИ
func WithNCIRCCToken(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'ncirccToken' cannot be empty")
		}

		sic.settings.ncirccToken = v

		return nil
	}
}

// WithRequestTimeout ограничение времени выполнения запроса от 1 до 60 сек.
func WithRequestTimeout(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 1 || v > 60 {
			return errors.New("the request execution time should be in the range from 1 to 60 seconds")
		}

		sic.settings.requestTimeout = v

		return nil
	}
}
