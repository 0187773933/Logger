package types

import (

)

type ConfigFile struct {
	ServerName string `yaml:"server_name"`
	ServerBaseUrl string `yaml:"server_base_url"`
	ServerLiveUrl string `yaml:"server_live_url"`
	ServerPrivateUrl string `yaml:"server_private_url"`
	ServerPublicUrl string `yaml:"server_public_url"`
	ServerPort string `yaml:"server_port"`
	ServerAPIKey string `yaml:"server_api_key"`
	ServerLoginUrlPrefix string `yaml:"server_login_url_prefix"`
	ServerUrlPrefix string `yaml:"server_url_prefix"`
	ServerCookieName string `yaml:"server_cookie_name"`
	ServerCookieSecret string `yaml:"server_cookie_secret"`
	ServerCookieAdminSecretMessage string `yaml:"server_cookie_admin_secret_message"`
	ServerCookieSecretMessage string `yaml:"server_cookie_secret_message"`
	ServerAllowOriginsString string `yaml:"server_allow_origins"`
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
	TimeZone string `yaml:"time_zone"`
	SaveFilesPath string `yaml:"save_files_path"`
	BoltDBPath string `yaml:"bolt_db_path"`
	EncryptionKey string `yaml:"encryption_key"`
	LogLevel string `yaml:"log_level"`
}