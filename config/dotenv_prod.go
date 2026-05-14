//go:build !dev

package config

func init() {
	// No-op: production reads from OS environment
}
