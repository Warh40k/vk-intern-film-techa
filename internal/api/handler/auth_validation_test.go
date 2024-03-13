package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRequest_Validate(t *testing.T) {
	tests := []struct {
		Label       string
		request     AuthRequest
		expectError bool
	}{
		{
			Label: "RightValues",
			request: AuthRequest{
				Username: "nikita",
				Password: "7dofwrj22jr",
			},
			expectError: false,
		},
		{
			Label: "WithoutUsername",
			request: AuthRequest{
				Username: "",
				Password: "JDfjfd;kdskf",
			},
			expectError: true,
		},
		{
			Label: "WithoutPassword",
			request: AuthRequest{
				Username: "dsfsfjjdi32f",
				Password: "",
			},
			expectError: true,
		},
		{
			Label: "SmallPassword",
			request: AuthRequest{
				Username: "dsfsfjjdi32f",
				Password: "dsdfdsd",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Label, func(t *testing.T) {
			validate := validator.New()
			err := validate.Struct(tt.request)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
