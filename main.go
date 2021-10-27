package main

import (
	"github.com/integrii/flaggy"
	"mt/config"
	"mt/log"
	"mt/minio"
	"mt/notify"
	"strings"
)

func main() {
	// General info
	flaggy.SetName("MinIO Toolkit")
	flaggy.SetDescription("A toolkit for various MinIO operations")
	flaggy.SetVersion("0.3.1")

	// define main flags
	configFile := "/etc/mt/config.yml"
	flaggy.String(&configFile, "c", "config", "Path of the configuration file in YAML format")

	logLevel := ""
	flaggy.String(&logLevel, "l", "log-level", "Log level (Overwrites config file. Available options are debug, info, warn and error.)")

	// define "upload" subcommand and its flags
	upload := flaggy.NewSubcommand("upload")
	upload.Description = "Upload files to MinIO"
	flaggy.AttachSubcommand(upload, 1)

	var source string
	upload.String(&source, "s", "source", "Source directory or file")

	var destination string
	upload.String(&destination, "d", "destination", "Destination for upload operation with \"server/bucket/path\" format")

	var recursive bool
	upload.Bool(&recursive, "r", "recursive", "Upload a directory recursively to MinIO")

	var md5sum bool
	upload.Bool(&md5sum, "md5", "md5-validation", "Validate uploads with md5sum calculation")

	var removeSourceFiles bool
	upload.Bool(&removeSourceFiles, "rm", "remove-source-files", "Remove source files after successful upload")

	var stopOnError bool
	upload.Bool(&stopOnError, "soe", "stop-on-error", "Stop on the first error and don't try to upload other files")

	var notifyErrors bool
	upload.Bool(&notifyErrors, "n", "notify-errors", "Notify errors using the channels in configuration file")

	// Parse flags
	flaggy.Parse()

	// Get params from config file
	params, err := config.NewParams(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Overwrite log level with flag
	if logLevel != "" {
		params.Log.Level = logLevel
	}

	// Get logger
	logger := log.NewLogger(&params.Log)

	// Get notifier
	notifier := notify.NewNotifier(&params.Notify)

	// Run "upload" subcommand
	if upload.Used {
		serverName := strings.Split(destination, "/")[0]
		serverParams, err := params.Server(serverName)
		if err != nil {
			logger.Fatal(err)
		}

		minio.Upload(
			logger,
			notifier,
			serverParams,
			minio.UploadParams{
				Source:            source,
				Destination:       strings.TrimPrefix(destination, serverName+"/"),
				Recursive:         recursive,
				RemoveSourceFiles: removeSourceFiles,
				Md5sum:            md5sum,
				StopOnError:       stopOnError,
				NotifyErrors:      notifyErrors,
			})

	} else {
		// If no subcommand used, print help and exit
		flaggy.ShowHelpAndExit("No subcommands passed")
	}
}
