package shaders

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/user"
)

//go:embed *.kage
var shaderFS embed.FS

type ShaderManager struct {
	noteShader2D     *ebiten.Shader
	noteShader3D     *ebiten.Shader
	holdNoteShader2D *ebiten.Shader
	holdNoteShader3D *ebiten.Shader
	holdTailShader3D *ebiten.Shader // New tail shader for 3D hold notes
	laneShader       *ebiten.Shader // Lane background shader
	markerShader     *ebiten.Shader // Measure/beat marker shader
	tunnelShader     *ebiten.Shader // Tunnel background shader
	hitEffectShader  *ebiten.Shader // Hit effect shader for on-hit shadows
}

var Manager *ShaderManager

func InitManager() error {
	Manager = &ShaderManager{}
	
	if err := Manager.loadShaders(); err != nil {
		return err
	}
	
	// Initialize lane renderer
	InitLaneRenderer()
	
	return nil
}

// ReinitManager reinitializes the shader manager (useful when settings change)
func ReinitManager() error {
	if Manager == nil {
		return InitManager()
	}
	// Reload shaders into existing manager instead of creating new one
	return Manager.loadShaders()
}

// ReinitSystem reinitializes the entire shader system (for future use if needed)
func ReinitSystem() error {
	if err := ReinitManager(); err != nil {
		return err
	}
	ReinitRenderer()
	InitLaneRenderer()
	return nil
}

func (sm *ShaderManager) loadShaders() error {
	var err error
	
	// Load 2D note shader
	noteSource2D, err := shaderFS.ReadFile("note.kage")
	if err != nil {
		return err
	}
	
	sm.noteShader2D, err = ebiten.NewShader(noteSource2D)
	if err != nil {
		return err
	}
	
	// Load 3D note shader
	noteSource3D, err := shaderFS.ReadFile("note3d.kage")
	if err != nil {
		return err
	}
	
	sm.noteShader3D, err = ebiten.NewShader(noteSource3D)
	if err != nil {
		return err
	}
	
	// Load 2D hold note shader
	holdSource2D, err := shaderFS.ReadFile("holdnote.kage")
	if err != nil {
		return err
	}
	
	sm.holdNoteShader2D, err = ebiten.NewShader(holdSource2D)
	if err != nil {
		return err
	}
	
	// Load 3D hold note shader
	holdSource3D, err := shaderFS.ReadFile("holdnote3d.kage")
	if err != nil {
		return err
	}
	
	sm.holdNoteShader3D, err = ebiten.NewShader(holdSource3D)
	if err != nil {
		return err
	}
	
	// Load 3D hold tail shader
	tailSource3D, err := shaderFS.ReadFile("holdtail3d.kage")
	if err != nil {
		return err
	}
	
	sm.holdTailShader3D, err = ebiten.NewShader(tailSource3D)
	if err != nil {
		return err
	}
	
	// Load lane background shader
	laneSource, err := shaderFS.ReadFile("lane.kage")
	if err != nil {
		return err
	}
	
	sm.laneShader, err = ebiten.NewShader(laneSource)
	if err != nil {
		return err
	}
	
	// Load marker shader
	markerSource, err := shaderFS.ReadFile("marker.kage")
	if err != nil {
		return err
	}
	
	sm.markerShader, err = ebiten.NewShader(markerSource)
	if err != nil {
		return err
	}
	
	// Load tunnel background shader
	tunnelSource, err := shaderFS.ReadFile("tunnel.kage")
	if err != nil {
		return err
	}
	
	sm.tunnelShader, err = ebiten.NewShader(tunnelSource)
	if err != nil {
		return err
	}
	
	// Load hit effect shader
	hitEffectSource, err := shaderFS.ReadFile("hiteffect.kage")
	if err != nil {
		return err
	}
	
	sm.hitEffectShader, err = ebiten.NewShader(hitEffectSource)
	if err != nil {
		return err
	}
	
	return nil
}

// GetNoteShader returns the appropriate note shader based on user settings
func (sm *ShaderManager) GetNoteShader() *ebiten.Shader {
	// Safety check - if user settings not available, default to 2D
	if user.S() == nil || !user.S().Use3DNotes {
		return sm.noteShader2D
	}
	return sm.noteShader3D
}

// GetHoldNoteShader returns the appropriate hold note shader based on user settings  
func (sm *ShaderManager) GetHoldNoteShader() *ebiten.Shader {
	// Safety check - if user settings not available, default to 2D
	if user.S() == nil || !user.S().Use3DNotes {
		return sm.holdNoteShader2D
	}
	return sm.holdNoteShader3D
}

// GetHoldTailShader returns the hold tail shader (only available in 3D mode)
func (sm *ShaderManager) GetHoldTailShader() *ebiten.Shader {
	// Only return tail shader in 3D mode, otherwise nil
	if user.S() == nil || !user.S().Use3DNotes {
		return nil
	}
	return sm.holdTailShader3D
}

// GetLaneShader returns the lane background shader
func (sm *ShaderManager) GetLaneShader() *ebiten.Shader {
	return sm.laneShader
}

// GetMarkerShader returns the marker shader
func (sm *ShaderManager) GetMarkerShader() *ebiten.Shader {
	return sm.markerShader
}

// GetTunnelShader returns the tunnel background shader
func (sm *ShaderManager) GetTunnelShader() *ebiten.Shader {
	return sm.tunnelShader
}

// GetHitEffectShader returns the hit effect shader
func (sm *ShaderManager) GetHitEffectShader() *ebiten.Shader {
	return sm.hitEffectShader
}