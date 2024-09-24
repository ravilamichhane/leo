package update

import (
	"log/slog"
	"os/exec"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update code",
	Long:  `Update code`,
	Run: func(cmd *cobra.Command, args []string) {
		exec_cmd := exec.Command("bash", "-c", "go install github.com/ravilmc/leo/cli/leo@latest")

		slog.Info("INSTALL", slog.String("message", "Leo Update Started"))
		if _, err := exec_cmd.Output(); err != nil {
			slog.Error("INSTALL", slog.Any("error", err))
		} else {
			slog.Info("INSTALL", slog.String("message", "Leo Updated Successfully"))
		}
	},
}
