// Code generated by github.com/ecordell/optgen. DO NOT EDIT.
package config

import (
	defaults "github.com/creasty/defaults"
	helpers "github.com/ecordell/optgen/helpers"
	"time"
)

type ConfigOption func(c *Config)

// NewConfigWithOptions creates a new Config with the passed in options set
func NewConfigWithOptions(opts ...ConfigOption) *Config {
	c := &Config{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// NewConfigWithOptionsAndDefaults creates a new Config with the passed in options set starting from the defaults
func NewConfigWithOptionsAndDefaults(opts ...ConfigOption) *Config {
	c := &Config{}
	defaults.MustSet(c)
	for _, o := range opts {
		o(c)
	}
	return c
}

// ToOption returns a new ConfigOption that sets the values from the passed in Config
func (c *Config) ToOption() ConfigOption {
	return func(to *Config) {}
}

// DebugMap returns a map form of Config for debugging
func (c Config) DebugMap() map[string]any {
	debugMap := map[string]any{}
	return debugMap
}

// ConfigWithOptions configures an existing Config with the passed in options set
func ConfigWithOptions(c *Config, opts ...ConfigOption) *Config {
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithOptions configures the receiver Config with the passed in options set
func (c *Config) WithOptions(opts ...ConfigOption) *Config {
	for _, o := range opts {
		o(c)
	}
	return c
}

type AppOption func(a *App)

// NewAppWithOptions creates a new App with the passed in options set
func NewAppWithOptions(opts ...AppOption) *App {
	a := &App{}
	for _, o := range opts {
		o(a)
	}
	return a
}

// NewAppWithOptionsAndDefaults creates a new App with the passed in options set starting from the defaults
func NewAppWithOptionsAndDefaults(opts ...AppOption) *App {
	a := &App{}
	defaults.MustSet(a)
	for _, o := range opts {
		o(a)
	}
	return a
}

// ToOption returns a new AppOption that sets the values from the passed in App
func (a *App) ToOption() AppOption {
	return func(to *App) {
		to.Name = a.Name
		to.Version = a.Version
		to.RunMode = a.RunMode
	}
}

// DebugMap returns a map form of App for debugging
func (a App) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["Name"] = helpers.DebugValue(a.Name, false)
	debugMap["Version"] = helpers.DebugValue(a.Version, false)
	debugMap["RunMode"] = helpers.DebugValue(a.RunMode, false)
	return debugMap
}

// AppWithOptions configures an existing App with the passed in options set
func AppWithOptions(a *App, opts ...AppOption) *App {
	for _, o := range opts {
		o(a)
	}
	return a
}

// WithOptions configures the receiver App with the passed in options set
func (a *App) WithOptions(opts ...AppOption) *App {
	for _, o := range opts {
		o(a)
	}
	return a
}

// WithName returns an option that can set Name on a App
func WithName(name string) AppOption {
	return func(a *App) {
		a.Name = name
	}
}

// WithVersion returns an option that can set Version on a App
func WithVersion(version string) AppOption {
	return func(a *App) {
		a.Version = version
	}
}

// WithRunMode returns an option that can set RunMode on a App
func WithRunMode(runMode string) AppOption {
	return func(a *App) {
		a.RunMode = runMode
	}
}

type LogOption func(l *Log)

// NewLogWithOptions creates a new Log with the passed in options set
func NewLogWithOptions(opts ...LogOption) *Log {
	l := &Log{}
	for _, o := range opts {
		o(l)
	}
	return l
}

// NewLogWithOptionsAndDefaults creates a new Log with the passed in options set starting from the defaults
func NewLogWithOptionsAndDefaults(opts ...LogOption) *Log {
	l := &Log{}
	defaults.MustSet(l)
	for _, o := range opts {
		o(l)
	}
	return l
}

// ToOption returns a new LogOption that sets the values from the passed in Log
func (l *Log) ToOption() LogOption {
	return func(to *Log) {
		to.Level = l.Level
	}
}

// DebugMap returns a map form of Log for debugging
func (l Log) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["Level"] = helpers.DebugValue(l.Level, false)
	return debugMap
}

// LogWithOptions configures an existing Log with the passed in options set
func LogWithOptions(l *Log, opts ...LogOption) *Log {
	for _, o := range opts {
		o(l)
	}
	return l
}

// WithOptions configures the receiver Log with the passed in options set
func (l *Log) WithOptions(opts ...LogOption) *Log {
	for _, o := range opts {
		o(l)
	}
	return l
}

// WithLevel returns an option that can set Level on a Log
func WithLevel(level string) LogOption {
	return func(l *Log) {
		l.Level = level
	}
}

type FeatureOption func(f *Feature)

// NewFeatureWithOptions creates a new Feature with the passed in options set
func NewFeatureWithOptions(opts ...FeatureOption) *Feature {
	f := &Feature{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// NewFeatureWithOptionsAndDefaults creates a new Feature with the passed in options set starting from the defaults
func NewFeatureWithOptionsAndDefaults(opts ...FeatureOption) *Feature {
	f := &Feature{}
	defaults.MustSet(f)
	for _, o := range opts {
		o(f)
	}
	return f
}

// ToOption returns a new FeatureOption that sets the values from the passed in Feature
func (f *Feature) ToOption() FeatureOption {
	return func(to *Feature) {
		to.ShutdownGracePeriod = f.ShutdownGracePeriod
	}
}

// DebugMap returns a map form of Feature for debugging
func (f Feature) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["ShutdownGracePeriod"] = helpers.DebugValue(f.ShutdownGracePeriod, false)
	return debugMap
}

// FeatureWithOptions configures an existing Feature with the passed in options set
func FeatureWithOptions(f *Feature, opts ...FeatureOption) *Feature {
	for _, o := range opts {
		o(f)
	}
	return f
}

// WithOptions configures the receiver Feature with the passed in options set
func (f *Feature) WithOptions(opts ...FeatureOption) *Feature {
	for _, o := range opts {
		o(f)
	}
	return f
}

// WithShutdownGracePeriod returns an option that can set ShutdownGracePeriod on a Feature
func WithShutdownGracePeriod(shutdownGracePeriod time.Duration) FeatureOption {
	return func(f *Feature) {
		f.ShutdownGracePeriod = shutdownGracePeriod
	}
}

type DataStoreOption func(d *DataStore)

// NewDataStoreWithOptions creates a new DataStore with the passed in options set
func NewDataStoreWithOptions(opts ...DataStoreOption) *DataStore {
	d := &DataStore{}
	for _, o := range opts {
		o(d)
	}
	return d
}

// NewDataStoreWithOptionsAndDefaults creates a new DataStore with the passed in options set starting from the defaults
func NewDataStoreWithOptionsAndDefaults(opts ...DataStoreOption) *DataStore {
	d := &DataStore{}
	defaults.MustSet(d)
	for _, o := range opts {
		o(d)
	}
	return d
}

// ToOption returns a new DataStoreOption that sets the values from the passed in DataStore
func (d *DataStore) ToOption() DataStoreOption {
	return func(to *DataStore) {
		to.Engine = d.Engine
		to.GcWindows = d.GcWindows
		to.GcMaxOperationTime = d.GcMaxOperationTime
		to.MigrationPhase = d.MigrationPhase
	}
}

// DebugMap returns a map form of DataStore for debugging
func (d DataStore) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["Engine"] = helpers.DebugValue(d.Engine, false)
	debugMap["GcWindows"] = helpers.DebugValue(d.GcWindows, false)
	debugMap["GcMaxOperationTime"] = helpers.DebugValue(d.GcMaxOperationTime, false)
	debugMap["MigrationPhase"] = helpers.DebugValue(d.MigrationPhase, false)
	return debugMap
}

// DataStoreWithOptions configures an existing DataStore with the passed in options set
func DataStoreWithOptions(d *DataStore, opts ...DataStoreOption) *DataStore {
	for _, o := range opts {
		o(d)
	}
	return d
}

// WithOptions configures the receiver DataStore with the passed in options set
func (d *DataStore) WithOptions(opts ...DataStoreOption) *DataStore {
	for _, o := range opts {
		o(d)
	}
	return d
}

// WithEngine returns an option that can set Engine on a DataStore
func WithEngine(engine string) DataStoreOption {
	return func(d *DataStore) {
		d.Engine = engine
	}
}

// WithGcWindows returns an option that can set GcWindows on a DataStore
func WithGcWindows(gcWindows time.Duration) DataStoreOption {
	return func(d *DataStore) {
		d.GcWindows = gcWindows
	}
}

// WithGcMaxOperationTime returns an option that can set GcMaxOperationTime on a DataStore
func WithGcMaxOperationTime(gcMaxOperationTime time.Duration) DataStoreOption {
	return func(d *DataStore) {
		d.GcMaxOperationTime = gcMaxOperationTime
	}
}

// WithMigrationPhase returns an option that can set MigrationPhase on a DataStore
func WithMigrationPhase(migrationPhase string) DataStoreOption {
	return func(d *DataStore) {
		d.MigrationPhase = migrationPhase
	}
}

type MysqlOption func(m *Mysql)

// NewMysqlWithOptions creates a new Mysql with the passed in options set
func NewMysqlWithOptions(opts ...MysqlOption) *Mysql {
	m := &Mysql{}
	for _, o := range opts {
		o(m)
	}
	return m
}

// NewMysqlWithOptionsAndDefaults creates a new Mysql with the passed in options set starting from the defaults
func NewMysqlWithOptionsAndDefaults(opts ...MysqlOption) *Mysql {
	m := &Mysql{}
	defaults.MustSet(m)
	for _, o := range opts {
		o(m)
	}
	return m
}

// ToOption returns a new MysqlOption that sets the values from the passed in Mysql
func (m *Mysql) ToOption() MysqlOption {
	return func(to *Mysql) {
		to.Host = m.Host
		to.Username = m.Username
		to.Password = m.Password
		to.Database = m.Database
		to.MaxIdleConnections = m.MaxIdleConnections
		to.MaxOpenConnections = m.MaxOpenConnections
		to.MaxConnectionLifeTime = m.MaxConnectionLifeTime
	}
}

// DebugMap returns a map form of Mysql for debugging
func (m Mysql) DebugMap() map[string]any {
	debugMap := map[string]any{}
	debugMap["Host"] = helpers.DebugValue(m.Host, false)
	debugMap["Username"] = helpers.DebugValue(m.Username, false)
	debugMap["Password"] = helpers.SensitiveDebugValue(m.Password)
	debugMap["Database"] = helpers.DebugValue(m.Database, false)
	debugMap["MaxIdleConnections"] = helpers.DebugValue(m.MaxIdleConnections, false)
	debugMap["MaxOpenConnections"] = helpers.DebugValue(m.MaxOpenConnections, false)
	debugMap["MaxConnectionLifeTime"] = helpers.DebugValue(m.MaxConnectionLifeTime, false)
	return debugMap
}

// MysqlWithOptions configures an existing Mysql with the passed in options set
func MysqlWithOptions(m *Mysql, opts ...MysqlOption) *Mysql {
	for _, o := range opts {
		o(m)
	}
	return m
}

// WithOptions configures the receiver Mysql with the passed in options set
func (m *Mysql) WithOptions(opts ...MysqlOption) *Mysql {
	for _, o := range opts {
		o(m)
	}
	return m
}

// WithHost returns an option that can set Host on a Mysql
func WithHost(host string) MysqlOption {
	return func(m *Mysql) {
		m.Host = host
	}
}

// WithUsername returns an option that can set Username on a Mysql
func WithUsername(username string) MysqlOption {
	return func(m *Mysql) {
		m.Username = username
	}
}

// WithPassword returns an option that can set Password on a Mysql
func WithPassword(password string) MysqlOption {
	return func(m *Mysql) {
		m.Password = password
	}
}

// WithDatabase returns an option that can set Database on a Mysql
func WithDatabase(database string) MysqlOption {
	return func(m *Mysql) {
		m.Database = database
	}
}

// WithMaxIdleConnections returns an option that can set MaxIdleConnections on a Mysql
func WithMaxIdleConnections(maxIdleConnections int) MysqlOption {
	return func(m *Mysql) {
		m.MaxIdleConnections = maxIdleConnections
	}
}

// WithMaxOpenConnections returns an option that can set MaxOpenConnections on a Mysql
func WithMaxOpenConnections(maxOpenConnections int) MysqlOption {
	return func(m *Mysql) {
		m.MaxOpenConnections = maxOpenConnections
	}
}

// WithMaxConnectionLifeTime returns an option that can set MaxConnectionLifeTime on a Mysql
func WithMaxConnectionLifeTime(maxConnectionLifeTime time.Duration) MysqlOption {
	return func(m *Mysql) {
		m.MaxConnectionLifeTime = maxConnectionLifeTime
	}
}
