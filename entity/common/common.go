package common

// VersionResponse carries application version info.
type VersionResponse struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// ChangelogEntry represents a single changelog item.
type ChangelogEntry struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

// SettingItem represents a single app setting entry.
type SettingItem struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
