package pslog

import "github.com/sirupsen/logrus"

type Option struct {
	logger    *logrus.Logger
	errFormat string
	psGener   GenFunc
}

type IOption interface {
	apply(option *Option)
}

type loggerOption struct {
	logger *logrus.Logger
}

func (o loggerOption) apply(option *Option) {
	option.logger = o.logger
}

func WithLogger(logger *logrus.Logger) IOption {
	return loggerOption{logger: logger}
}

type ErrFormatOption struct {
	errFormat string
}

func (o ErrFormatOption) apply(option *Option) {
	option.errFormat = o.errFormat
}

func WithErrorFormat(errFormat string) IOption {
	return ErrFormatOption{errFormat: errFormat}
}

type PSGenerOption struct {
	psGener GenFunc
}

func (o PSGenerOption) apply(option *Option) {
	option.psGener = o.psGener
}

func WithPSGenerOption(psGener GenFunc) IOption {
	return PSGenerOption{psGener: psGener}
}
