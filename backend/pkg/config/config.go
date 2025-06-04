package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Storage   StorageConfig   `yaml:"storage"`
	Generator GeneratorConfig `yaml:"generator"`
	Logging   LoggingConfig   `yaml:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Environment  string `yaml:"environment"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type        string `yaml:"type"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Name        string `yaml:"name"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	SSLMode     string `yaml:"ssl_mode"`
	MaxConns    int    `yaml:"max_conns"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxLifetime int    `yaml:"max_lifetime"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	Type      string `yaml:"type"`
	BasePath  string `yaml:"base_path"`
	OutputDir string `yaml:"output_dir"`
	ImageDir  string `yaml:"image_dir"`
	TempDir   string `yaml:"temp_dir"`
}

// GeneratorConfig holds card generation configuration
type GeneratorConfig struct {
	ImageWidth     int                 `yaml:"image_width"`
	ImageHeight    int                 `yaml:"image_height"`
	DPI            int                 `yaml:"dpi"`
	Quality        int                 `yaml:"quality"`
	Format         string              `yaml:"format"`
	TemplatesPath  string              `yaml:"templates_path"`
	FontsPath      string              `yaml:"fonts_path"`
	DefaultFont    string              `yaml:"default_font"`
	TextRendering  TextRenderingConfig `yaml:"text_rendering"`
	ArtProcessing  ArtProcessingConfig `yaml:"art_processing"`
	ParallelJobs   int                 `yaml:"parallel_jobs"`
	EnableCaching  bool                `yaml:"enable_caching"`
	CacheDirectory string              `yaml:"cache_directory"`
}

// TextRenderingConfig holds text rendering configuration
type TextRenderingConfig struct {
	DefaultFontSize float64           `yaml:"default_font_size"`
	LineSpacing     float64           `yaml:"line_spacing"`
	DefaultColor    string            `yaml:"default_color"`
	AntiAliasing    bool              `yaml:"anti_aliasing"`
	FontCache       bool              `yaml:"font_cache"`
	CustomFonts     map[string]string `yaml:"custom_fonts"`
}

// ArtProcessingConfig holds art processing configuration
type ArtProcessingConfig struct {
	EnablePlaceholder bool   `yaml:"enable_placeholder"`
	PlaceholderColor  string `yaml:"placeholder_color"`
	ResizeAlgorithm   string `yaml:"resize_algorithm"`
	EnableFilters     bool   `yaml:"enable_filters"`
	DefaultFilter     string `yaml:"default_filter"`
	MaxImageSize      int64  `yaml:"max_image_size"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// LoadConfig loads configuration from multiple sources
func LoadConfig(env string) (*Config, error) {
	// Start with default configuration
	config := getDefaultConfig()

	// Load from YAML file if it exists
	configPath := getConfigPath(env)
	if _, err := os.Stat(configPath); err == nil {
		if err := loadFromYAML(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from YAML: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnvironment(config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "localhost",
			ReadTimeout:  30,
			WriteTimeout: 30,
			Environment:  "development",
		},
		Database: DatabaseConfig{
			Type:        "memory",
			Host:        "localhost",
			Port:        5432,
			Name:        "cardgen",
			User:        "cardgen",
			SSLMode:     "disable",
			MaxConns:    10,
			MaxIdle:     5,
			MaxLifetime: 300,
		},
		Storage: StorageConfig{
			Type:      "file",
			BasePath:  "./data",
			OutputDir: "./output",
			ImageDir:  "./output/images",
			TempDir:   "./tmp",
		},
		Generator: GeneratorConfig{
			ImageWidth:    750,
			ImageHeight:   1050,
			DPI:           300,
			Quality:       95,
			Format:        "png",
			TemplatesPath: "./templates",
			FontsPath:     "./fonts",
			DefaultFont:   "arial.ttf",
			TextRendering: TextRenderingConfig{
				DefaultFontSize: 12.0,
				LineSpacing:     1.2,
				DefaultColor:    "#000000",
				AntiAliasing:    true,
				FontCache:       true,
				CustomFonts:     make(map[string]string),
			},
			ArtProcessing: ArtProcessingConfig{
				EnablePlaceholder: true,
				PlaceholderColor:  "#CCCCCC",
				ResizeAlgorithm:   "lanczos",
				EnableFilters:     false,
				DefaultFilter:     "none",
				MaxImageSize:      10 * 1024 * 1024, // 10MB
			},
			ParallelJobs:   4,
			EnableCaching:  true,
			CacheDirectory: "./cache",
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
	}
}

// getConfigPath returns the configuration file path based on environment
func getConfigPath(env string) string {
	if env == "" {
		env = getEnv("APP_ENV", "development")
	}

	configDir := getEnv("CONFIG_DIR", "./config")
	return filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
}

// loadFromYAML loads configuration from a YAML file
func loadFromYAML(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return nil
}

// loadFromEnvironment loads configuration from environment variables
func loadFromEnvironment(config *Config) {
	// Server configuration
	if port := getEnvInt("SERVER_PORT", 0); port > 0 {
		config.Server.Port = port
	}
	if host := getEnv("SERVER_HOST", ""); host != "" {
		config.Server.Host = host
	}
	if env := getEnv("APP_ENV", ""); env != "" {
		config.Server.Environment = env
	}

	// Database configuration
	if dbType := getEnv("DB_TYPE", ""); dbType != "" {
		config.Database.Type = dbType
	}
	if dbHost := getEnv("DB_HOST", ""); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := getEnvInt("DB_PORT", 0); dbPort > 0 {
		config.Database.Port = dbPort
	}
	if dbName := getEnv("DB_NAME", ""); dbName != "" {
		config.Database.Name = dbName
	}
	if dbUser := getEnv("DB_USER", ""); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := getEnv("DB_PASSWORD", ""); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if sslMode := getEnv("DB_SSL_MODE", ""); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Storage configuration
	if storageType := getEnv("STORAGE_TYPE", ""); storageType != "" {
		config.Storage.Type = storageType
	}
	if basePath := getEnv("STORAGE_BASE_PATH", ""); basePath != "" {
		config.Storage.BasePath = basePath
	}
	if outputDir := getEnv("OUTPUT_DIR", ""); outputDir != "" {
		config.Storage.OutputDir = outputDir
	}
	if imageDir := getEnv("IMAGE_DIR", ""); imageDir != "" {
		config.Storage.ImageDir = imageDir
	}

	// Generator configuration
	if width := getEnvInt("GENERATOR_IMAGE_WIDTH", 0); width > 0 {
		config.Generator.ImageWidth = width
	}
	if height := getEnvInt("GENERATOR_IMAGE_HEIGHT", 0); height > 0 {
		config.Generator.ImageHeight = height
	}
	if format := getEnv("GENERATOR_FORMAT", ""); format != "" {
		config.Generator.Format = format
	}
	if templatesPath := getEnv("TEMPLATES_PATH", ""); templatesPath != "" {
		config.Generator.TemplatesPath = templatesPath
	}
	if fontsPath := getEnv("FONTS_PATH", ""); fontsPath != "" {
		config.Generator.FontsPath = fontsPath
	}

	// Logging configuration
	if level := getEnv("LOG_LEVEL", ""); level != "" {
		config.Logging.Level = level
	}
	if format := getEnv("LOG_FORMAT", ""); format != "" {
		config.Logging.Format = format
	}
	if output := getEnv("LOG_OUTPUT", ""); output != "" {
		config.Logging.Output = output
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	if config.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	// Validate database configuration
	validDBTypes := []string{"memory", "postgres", "sqlite", "file"}
	if !contains(validDBTypes, config.Database.Type) {
		return fmt.Errorf("invalid database type: %s", config.Database.Type)
	}

	// Validate storage configuration
	validStorageTypes := []string{"memory", "file", "s3", "gcs"}
	if !contains(validStorageTypes, config.Storage.Type) {
		return fmt.Errorf("invalid storage type: %s", config.Storage.Type)
	}

	// Validate generator configuration
	if config.Generator.ImageWidth <= 0 {
		return fmt.Errorf("invalid image width: %d", config.Generator.ImageWidth)
	}
	if config.Generator.ImageHeight <= 0 {
		return fmt.Errorf("invalid image height: %d", config.Generator.ImageHeight)
	}
	if config.Generator.ParallelJobs <= 0 {
		return fmt.Errorf("parallel jobs must be greater than 0")
	}

	// Validate logging configuration
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, config.Logging.Level) {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// GetConnectionString builds a database connection string
func (db *DatabaseConfig) GetConnectionString() string {
	switch strings.ToLower(db.Type) {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode)
	case "sqlite":
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", db.Name)
	default:
		return ""
	}
}
