package pslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

// 默认 Logger
var (
	defaultPrefix    = "moot"
	defaultErrFormat = "-> %s: "
	Logger           = New(defaultPrefix)
	mu               sync.RWMutex
)

var loggersMap = map[string]*PSLogger{}

// PSLogger 带前缀的 Logger
type PSLogger struct {
	psFull        string      // 完整PS
	psList        []string    // PS 列表
	parent        *PSLogger   // 父级
	kids          []*PSLogger // 子
	isPsFullGened bool        // 是否已经生成？
	opts          *Option     // 选项集
}

func defaultPrefixGenerator(pslist []string) string {
	return fmt.Sprintf("[%s]", strings.Join(pslist, "."))
}

type GenFunc func(pslist []string) string

func (l *PSLogger) Stringer() string {
	return fmt.Sprintf("<PSLogger>: %s", l.genPrefix())
}

// New 新建一个 ps logger
func New(prefix string, opts ...IOption) *PSLogger {

	mu.Lock()
	defer mu.Unlock()

	if loggersMap[prefix] != nil {
		return loggersMap[prefix]
	}

	if prefix == "" {
		prefix = defaultPrefix
	}
	ps := []string{prefix}

	kidsList := make([]*PSLogger, 0)
	obj := &PSLogger{
		psList: ps,
		parent: nil,
		kids:   kidsList,
	}

	obj.opts = new(Option)
	obj.setDefaultOpts()

	for _, opt := range opts {
		opt.apply(obj.opts)
	}

	obj.psFull = obj.genPrefix()
	obj.isPsFullGened = true

	return obj
}

func (l *PSLogger) setDefaultOpts() {
	// 如果没有设置 logger，默认使用 logurs 的默认 logger
	l.opts.logger = logrus.StandardLogger()
	l.opts.errFormat = defaultErrFormat
	l.opts.psGener = defaultPrefixGenerator
}

// genPrefix 根据 pslsit 生成前缀
func (l *PSLogger) genPrefix() string {
	if l.isPsFullGened {
		return l.psFull
	}

	gener := defaultPrefixGenerator
	if l.opts.psGener != nil {
		gener = l.opts.psGener
	}

	vs := gener(l.psList)

	l.psFull = vs
	l.isPsFullGened = true

	return vs
}

// GetPrefix
func (l *PSLogger) GetPrefix() string {
	return l.genPrefix()
}

// genAppendedPs 生成 前缀列表
func (l *PSLogger) genAppendedPs(prefix string) []string {
	ps := make([]string, len(l.psList))
	for idx, p := range l.psList {
		ps[idx] = p
	}
	ps = append(ps, prefix)
	return ps
}

// p 生成前缀并组装log内容，单条日志非format时使用
func (l *PSLogger) p(args ...interface{}) []interface{} {
	vs := []interface{}{fmt.Sprintf("%s", l.genPrefix())}
	for _, arg := range args {
		vs = append(vs, arg)
	}
	return vs
}

// pf 生成前缀，单条日志 format 时使用
func (l *PSLogger) pf(format string) string {
	format = fmt.Sprintf("%s %s", l.genPrefix(), format)
	return format
}

// pson 生前前缀并组装log内容，JSON 格式专用
func (l *PSLogger) pson(ps string, args interface{}) []interface{} {
	rs, _ := json.Marshal(args)
	jstring := string(rs)
	allargs := []interface{}{
		fmt.Sprintf("%s:", ps),
		jstring,
	}
	vs := l.p(allargs...)
	return vs
}

// Sub 生成一个子 logger，如果已经有了，则返回
func (l *PSLogger) Sub(prefix string, opts ...IOption) *PSLogger {
	mu.Lock()
	defer mu.Unlock()

	newPlist := l.genAppendedPs(prefix)

	fullPs := l.opts.psGener(newPlist)
	if loggersMap[fullPs] != nil {
		return loggersMap[fullPs]
	}

	npl := &PSLogger{
		psList: newPlist,
		parent: l,
	}

	npl.opts = new(Option)
	*npl.opts = *l.opts

	l.kids = append(l.kids, npl)

	for _, opt := range opts {
		opt.apply(npl.opts)
	}

	return npl
}

/*
以下为普通的 log 逻辑
*/

func (l *PSLogger) Debug(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.DebugLevel) {
		l.opts.logger.Debug(l.p(args...)...)
	}
}
func (l *PSLogger) Debugln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.DebugLevel) {
		l.opts.logger.Debugln(l.p(args...)...)
	}
}
func (l *PSLogger) Debugf(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.DebugLevel) {
		l.opts.logger.Debugf(l.pf(format), args...)
	}
}
func (l *PSLogger) DebugJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.DebugLevel) {
		l.opts.logger.Debugln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) Info(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.InfoLevel) {
		l.opts.logger.Info(l.p(args...)...)
	}
}
func (l *PSLogger) Infoln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.InfoLevel) {
		l.opts.logger.Infoln(l.p(args...)...)
	}
}
func (l *PSLogger) Infof(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.InfoLevel) {
		l.opts.logger.Infof(l.pf(format), args...)
	}
}
func (l *PSLogger) InfoJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.InfoLevel) {
		l.opts.logger.Infoln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) Warn(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warning(l.p(args...)...)
	}
}
func (l *PSLogger) Warnln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningln(l.p(args...)...)
	}
}
func (l *PSLogger) Warnf(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningf(l.pf(format), args...)
	}
}
func (l *PSLogger) WarnJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) Warning(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warning(l.p(args...)...)
	}
}
func (l *PSLogger) Warningln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningln(l.p(args...)...)
	}
}
func (l *PSLogger) Warningf(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningf(l.pf(format), args...)
	}
}

func (l *PSLogger) WarningE(format string, err error) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		if err == nil {
			err = errors.New("err=nil")
		}
		l.opts.logger.Warningf(l.pf(format), err.Error())
	}
}

func (l *PSLogger) WarningJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.WarnLevel) {
		l.opts.logger.Warningln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) Error(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.opts.logger.Error(l.p(args...)...)
	}
}
func (l *PSLogger) Errorln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.opts.logger.Errorln(l.p(args...)...)
	}
}
func (l *PSLogger) Errorf(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.opts.logger.Errorf(l.pf(format), args...)
	}
}

func (l *PSLogger) ErrorE(format string, err error) {
	if l.opts.logger.IsLevelEnabled(logrus.ErrorLevel) {
		if err == nil {
			err = errors.New("err=nil")
		}
		l.opts.logger.Errorf(l.pf(format), err.Error())
	}
}

func (l *PSLogger) ErrorJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.opts.logger.Errorln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) Panic(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.PanicLevel) {
		l.opts.logger.Panic(l.p(args...)...)
	}
}
func (l *PSLogger) Panicln(args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.PanicLevel) {
		l.opts.logger.Panicln(l.p(args...)...)
	}
}
func (l *PSLogger) Panicf(format string, args ...interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.PanicLevel) {
		l.opts.logger.Panicf(l.pf(format), args...)
	}
}
func (l *PSLogger) PanicJson(ps string, args interface{}) {
	if l.opts.logger.IsLevelEnabled(logrus.PanicLevel) {
		l.opts.logger.Panicln(l.pson(ps, args)...)
	}
}

func (l *PSLogger) NewErrWithMsgs(sl ...string) error {
	slist := []string{fmt.Sprintf(l.opts.errFormat, l.GetPrefix())}
	if len(sl) == 0 {
		return errors.New(strings.Join(slist, " "))
	}

	for _, s := range sl {
		slist = append(slist, s)
	}
	return errors.New(strings.Join(slist, " "))
}
func (l *PSLogger) NewErrWithFormat(format string, values ...interface{}) error {
	slist := []string{fmt.Sprintf(l.opts.errFormat, l.GetPrefix()), fmt.Sprintf(format, values...)}
	return errors.New(strings.Join(slist, " "))
}
func (l *PSLogger) NewErrWithErrs(s string, errs ...error) error {
	slist := []string{fmt.Sprintf(l.opts.errFormat, l.GetPrefix()), s}
	if len(errs) > 0 {
		for _, err := range errs {
			slist = append(slist, err.Error())
		}
	}
	return errors.New(strings.Join(slist, " "))
}
