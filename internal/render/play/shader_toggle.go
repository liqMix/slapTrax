package play

import "github.com/liqmix/slaptrax/internal/logger"

// ToggleShaderRendering toggles between shader and vertex-based rendering
func ToggleShaderRendering() {
	ShaderRenderingEnabled = !ShaderRenderingEnabled
	if ShaderRenderingEnabled {
		logger.Info("Shader-based note rendering enabled")
	} else {
		logger.Info("Vertex-based note rendering enabled")
	}
}

// IsShaderRenderingEnabled returns whether shader rendering is currently enabled
func IsShaderRenderingEnabled() bool {
	return ShaderRenderingEnabled
}

// SetShaderRendering enables or disables shader rendering
func SetShaderRendering(enabled bool) {
	ShaderRenderingEnabled = enabled
	if enabled {
		logger.Info("Shader-based note rendering enabled")
	} else {
		logger.Info("Vertex-based note rendering enabled")
	}
}