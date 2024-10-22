package cmd

import (
	"fmt"
	"github.com/nilpntr/gluster-exporter/internal/handlers"
	"github.com/nilpntr/gluster-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gluster-exporter",
	Short: "Gluster Exporter is an exporter for gluster to prometheus",
	RunE: func(cmd *cobra.Command, args []string) error {
		metricsClient, err := metrics.New()
		if err != nil {
			return err
		}

		registry := prometheus.NewRegistry()
		registry.MustRegister(metricsClient, version.NewCollector("gluster_exporter"))

		mux := http.NewServeMux()

		mux.Handle(viper.GetString("web_metrics_path"), promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		mux.HandleFunc("/healthz", handlers.Healthz)

		zap.L().Sugar().Infof("Starting exporter on: %v", viper.GetInt("web_listen_address"))

		if err := http.ListenAndServe(viper.GetString("web_listen_address"), mux); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
	rootCmd.Flags().String("log.level", "info", "Which log level to use, allowed levels: [info, error, debug]")
	rootCmd.Flags().String("web.listen-address", ":9106", "Address to listen on for web interface")
	rootCmd.Flags().String("web.metrics-path", "/metrics", "Path under which to expose metrics")
	rootCmd.Flags().String("gluster.volumes", "_all", "Comma separated volume names: vol1,vol2,vol3. Default is '_all' to scrape all metrics")
	rootCmd.Flags().String("gluster.binary", "/usr/sbin/gluster", "Path to the gluster binary")
	rootCmd.Flags().Bool("profile", false, "Enable gluster profiling reports")
	rootCmd.Flags().Bool("quota", false, "Enable gluster quota reports")
}

func initConfig() {
	_ = viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log.level"))
	_ = viper.BindPFlag("web_listen_address", rootCmd.Flags().Lookup("web.listen-address"))
	_ = viper.BindPFlag("web_metrics_path", rootCmd.Flags().Lookup("web.metrics-path"))
	_ = viper.BindPFlag("gluster_volumes", rootCmd.Flags().Lookup("gluster.volumes"))
	_ = viper.BindPFlag("gluster_binary", rootCmd.Flags().Lookup("gluster.binary"))
	_ = viper.BindPFlag("profile", rootCmd.Flags().Lookup("profile"))
	_ = viper.BindPFlag("quota", rootCmd.Flags().Lookup("quota"))

	viper.AutomaticEnv()
}

func initLogger() {
	var level zapcore.Level
	switch viper.GetString("log_level") {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "timestamp"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     enc,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	logger := zap.Must(zapCfg.Build())
	logger.Info("Logger initialized 🎉")
	zap.ReplaceGlobals(logger)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
