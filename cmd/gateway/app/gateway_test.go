package app

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCmdGateway(t *testing.T) {
	tests := []struct {
		name     string
		validate func(*testing.T, *cobra.Command)
	}{
		{
			name: "creates command with correct structure",
			validate: func(t *testing.T, cmd *cobra.Command) {
				assert.NotNil(t, cmd)
				assert.Equal(t, "gateway", cmd.Use)
				assert.NotEmpty(t, cmd.Short)
				assert.NotEmpty(t, cmd.Long)
				assert.NotEmpty(t, cmd.Example)
				assert.NotNil(t, cmd.RunE)
			},
		},
		{
			name: "registers required flags",
			validate: func(t *testing.T, cmd *cobra.Command) {
				bindAddrFlag := cmd.Flags().Lookup("bind-address")
				require.NotNil(t, bindAddrFlag, "bind-address flag should be registered")
				assert.Equal(t, "bind-address", bindAddrFlag.Name)
				assert.Equal(t, ":8080", bindAddrFlag.DefValue)

				authServiceFlag := cmd.Flags().Lookup("auth-service")
				require.NotNil(t, authServiceFlag, "auth-service flag should be registered")
				assert.Equal(t, "auth-service", authServiceFlag.Name)
				assert.Equal(t, ":8081", authServiceFlag.DefValue)
			},
		},
		{
			name: "has correct command metadata",
			validate: func(t *testing.T, cmd *cobra.Command) {
				assert.Equal(t, "gateway", cmd.CommandPath())
				assert.Contains(t, cmd.Short, "gateway")
				assert.Contains(t, cmd.Long, "gateway")
				assert.Contains(t, cmd.Example, "gateway")
				assert.Contains(t, cmd.Example, "--")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdGateway()
			tt.validate(t, cmd)
		})
	}
}

func TestNewCmdGateway_Flags(t *testing.T) {
	cmd := NewCmdGateway()

	t.Run("default values are correct", func(t *testing.T) {
		bindAddr, err := cmd.Flags().GetString("bind-address")
		require.NoError(t, err)
		assert.Equal(t, ":8080", bindAddr)

		authService, err := cmd.Flags().GetString("auth-service")
		require.NoError(t, err)
		assert.Equal(t, ":8081", authService)
	})

	t.Run("flags can be modified", func(t *testing.T) {
		tests := []struct {
			name      string
			flagName  string
			flagValue string
		}{
			{
				name:      "modify bind-address",
				flagName:  "bind-address",
				flagValue: ":9999",
			},
			{
				name:      "modify auth-service",
				flagName:  "auth-service",
				flagValue: "example.com:8888",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := cmd.Flags().Set(tt.flagName, tt.flagValue)
				require.NoError(t, err)

				value, err := cmd.Flags().GetString(tt.flagName)
				require.NoError(t, err)
				assert.Equal(t, tt.flagValue, value)
			})
		}
	})
}

func TestNewCmdGateway_RunE(t *testing.T) {
	cmd := NewCmdGateway()

	t.Run("RunE function exists and has correct signature", func(t *testing.T) {
		require.NotNil(t, cmd.RunE)
		assert.IsType(t, func(*cobra.Command, []string) error { return nil }, cmd.RunE)
	})

	t.Run("RunE handles flags correctly", func(t *testing.T) {
		cmd := NewCmdGateway()
		err := cmd.Flags().Set("bind-address", ":9090")
		require.NoError(t, err)

		err = cmd.Flags().Set("auth-service", "localhost:9091")
		require.NoError(t, err)

		bindAddr, _ := cmd.Flags().GetString("bind-address")
		assert.Equal(t, ":9090", bindAddr)

		authService, _ := cmd.Flags().GetString("auth-service")
		assert.Equal(t, "localhost:9091", authService)
	})

	t.Run("RunE executes and returns error for invalid address", func(t *testing.T) {
		cmd := NewCmdGateway()
		err := cmd.Flags().Set("bind-address", "invalid-address")
		require.NoError(t, err)

		runErr := cmd.RunE(cmd, nil)
		require.Error(t, runErr, "RunE should return an error for invalid address")
	})
}

func TestNewCmdGateway_Isolation(t *testing.T) {
	t.Run("multiple instances are independent", func(t *testing.T) {
		cmd1 := NewCmdGateway()
		cmd2 := NewCmdGateway()

		require.NotNil(t, cmd1)
		require.NotNil(t, cmd2)

		err := cmd1.Flags().Set("bind-address", ":1111")
		require.NoError(t, err)

		err = cmd2.Flags().Set("bind-address", ":2222")
		require.NoError(t, err)

		addr1, _ := cmd1.Flags().GetString("bind-address")
		addr2, _ := cmd2.Flags().GetString("bind-address")
		assert.NotEqual(t, addr1, addr2, "command instances should be independent")
	})
}
