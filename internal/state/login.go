package state

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

type LoginState struct {
	types.BaseGameState
	modal   *ui.LoginModal
	loading bool
}

func NewLoginState() *LoginState {
	s := &LoginState{
		modal: ui.NewLoginModal(),
	}

	// Center modal
	s.modal.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	s.modal.SetOnLogin(s.handleLogin)
	s.modal.SetOnContinue(func() {
		s.SetNextState(types.GameStateTitle, nil)
	})
	return s
}

func (s *LoginState) handleLogin(username, password string) {
	if s.loading {
		return
	}

	// Basic validation
	username = strings.TrimSpace(username)
	if username == "" {
		s.modal.SetError(l.String(l.ERROR_USERNAME_REQUIRED))
		return
	}
	if password == "" {
		s.modal.SetError(l.String(l.ERROR_PASSWORD_REQUIRED))
		return
	}

	s.loading = true
	s.modal.SetError("")

	// Attempt login
	err := external.Login(username, password, true)
	if err != nil {
		// Check if user doesn't exist (404)
		if strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusNotFound)) {
			// Try to register
			err = external.Register(username, password)
			if err != nil {
				s.modal.SetError(l.String(l.ERROR_REGISTER_FAIL) + "\n" + err.Error())
				s.loading = false
				return
			}

			// Try login again
			err = external.Login(username, password, true)
			if err != nil {
				s.modal.SetError(l.String(l.ERROR_LOGIN_REGISTER_FAIL) + "\n" + err.Error())
				s.loading = false
				return
			}
		} else {
			s.modal.SetError(l.String(l.ERROR_LOGIN_FAILED) + "\n" + err.Error())
			s.loading = false
			return
		}
	}
	s.SetNextState(types.GameStateTitle, nil)
	s.loading = false
}

func (s *LoginState) Update() error {
	s.BaseGameState.Update()
	s.modal.Update()
	return nil
}

func (s *LoginState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.modal.Draw(screen, opts)
}
