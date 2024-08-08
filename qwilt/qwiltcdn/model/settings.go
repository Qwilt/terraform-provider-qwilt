package model

type Settings struct {
	// CDN
	EnvType   string `tfsdk:"env_type"`
	Username  string `tfsdk:"username"`
	Password  string `tfsdk:"password"`
	XApiToken string `tfsdk:"token"`
}
