package version

type AzureCliVersion struct {
	AzureCli          string     `json:"azure-cli"`
	AzureCliCore      string     `json:"azure-cli-core"`
	AzureCliTelemetry string     `json:"azure-cli-telemetry"`
	Extensions        Extensions `json:"extensions"`
}
type Extensions struct {
	Account      string `json:"account"`
	AzureDevops  string `json:"azure-devops"`
	LogAnalytics string `json:"log-analytics"`
}
