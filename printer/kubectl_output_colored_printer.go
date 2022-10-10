package printer

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/kubectl"
)

// KubectlOutputColoredPrinter is a printer to print data depending on
// which kubectl subcommand is executed.
type KubectlOutputColoredPrinter struct {
	SubcommandInfo    *kubectl.SubcommandInfo
	Recursive         bool
	ObjFreshThreshold time.Duration
}

func ColorStatus(status string) (color.Color, bool) {
	switch status {
	case
		// from https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/events/event.go
		// Container event reason list
		"Failed",
		"BackOff",
		"ExceededGracePeriod",
		// Pod event reason list
		"FailedKillPod",
		"FailedCreatePodContainer",
		// "Failed",
		"NetworkNotReady",
		// Image event reason list
		// "Failed",
		"InspectFailed",
		"ErrImageNeverPull",
		// "BackOff",
		// kubelet event reason list
		"NodeNotSchedulable",
		"KubeletSetupFailed",
		"FailedAttachVolume",
		"FailedMount",
		"VolumeResizeFailed",
		"FileSystemResizeFailed",
		"FailedMapVolume",
		"ContainerGCFailed",
		"ImageGCFailed",
		"FailedNodeAllocatableEnforcement",
		"FailedCreatePodSandBox",
		"FailedPodSandBoxStatus",
		"FailedMountOnFilesystemMismatch",
		// Image manager event reason list
		"InvalidDiskCapacity",
		"FreeDiskSpaceFailed",
		// Probe event reason list
		"Unhealthy",
		// Pod worker event reason list
		"FailedSync",
		// Config event reason list
		"FailedValidation",
		// Lifecycle hooks
		"FailedPostStartHook",
		"FailedPreStopHook",

		// some other status
		"ContainerStatusUnknown",
		"CrashLoopBackOff",
		"ImagePullBackOff",
		"Evicted",
		"FailedScheduling",
		"Error",
		"ErrImagePull":
		return color.Red, true
	case
		// from https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/events/event.go
		// Container event reason list
		"Killing",
		"Preempting",
		// Pod event reason list
		// Image event reason list
		"Pulling",
		// kubelet event reason list
		"NodeNotReady",
		"NodeSchedulable",
		"Starting",
		"AlreadyMountedVolume",
		"SuccessfulAttachVolume",
		"SuccessfulMountVolume",
		"NodeAllocatableEnforced",
		// Image manager event reason list
		// Probe event reason list
		"ProbeWarning",
		// Pod worker event reason list
		// Config event reason list
		// Lifecycle hooks

		// some other status
		"Pending",
		"ContainerCreating",
		"PodInitializing",
		"Terminating",
		"Warning":
		return color.Yellow, true
	case
		"Completed":
		return color.Gray, true
	case
		"Running":
		return color.Green, true
	}
	// some ok status, not colored:
	// "Pulled",
	// "Created",
	// "Rebooted",
	// "SandboxChanged",
	// "VolumeResizeSuccessful",
	// "FileSystemResizeSuccessful",
	// "NodeReady",
	// "Started",
	// "Normal",
	return color.GrayLight, false
}

// Print reads r then write it to w, its format is based on kubectl subcommand.
// If given subcommand is not supported by the printer, it prints data in Green.
func (kp *KubectlOutputColoredPrinter) Print(r io.Reader, w io.Writer) {
	withHeader := !kp.SubcommandInfo.NoHeader

	var printer Printer = &SingleColoredPrinter{Color: color.Green} // default in green

	switch kp.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		printer = NewTablePrinter(withHeader, nil)

	case kubectl.APIVersions:
		printer = NewTablePrinter(false, nil) // api-versions always doesn't have header

	case kubectl.Get:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.None, kp.SubcommandInfo.FormatOption == kubectl.Wide:
			printer = NewTablePrinter(
				withHeader,
				func(_ int, column string) (color.Color, bool) {
					// first try to match a status
					col, matched := ColorStatus(column)
					if matched {
						return col, true
					}

					if column == "<none>" {
						return color.Gray, true
					}

					// When Readiness is "n/m" then yellow
					if strings.Count(column, "/") == 1 {
						if arr := strings.Split(column, "/"); arr[0] != arr[1] {
							_, e1 := strconv.Atoi(arr[0])
							_, e2 := strconv.Atoi(arr[1])
							if e1 == nil && e2 == nil { // check both is number
								return color.Yellow, true
							}
						}
						return color.Green, true
					}

					if strings.Contains(column, "ago)") {
						return color.Yellow, true
					}

					if isDuration, fresh := checkIfObjFresh(column, kp.ObjFreshThreshold); isDuration {
						if fresh {
							return color.Green, true
						} else {
							return color.MagentaDark, true
						}
					}

					if isIp(column) {
						return color.Blue, true
					}

					return color.GrayLight, false
				},
			)
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{}
		}

	case kubectl.Describe:
		printer = &DescribePrinter{
			TablePrinter: NewTablePrinter(false, func(_ int, column string) (color.Color, bool) {
				return ColorStatus(column)
			},
			),
		}
	case kubectl.Explain:
		printer = &ExplainPrinter{
			Recursive: kp.Recursive,
		}
	case kubectl.Version:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{}
		case kp.SubcommandInfo.Short:
			printer = &VersionShortPrinter{}
		default:
			printer = &VersionPrinter{}
		}
	case kubectl.Options:
		printer = &OptionsPrinter{}
	case kubectl.Apply:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{}
		default:
			printer = &ApplyPrinter{}
		}
	}

	if kp.SubcommandInfo.Help {
		printer = &SingleColoredPrinter{Color: color.GrayLight}
	}

	printer.Print(r, w)
}

func isIp(value string) bool {
	match, _ := regexp.MatchString("^\\d+\\.\\d+\\.\\d+\\.\\d+", value)
	return match
}

func checkIfObjFresh(value string, threshold time.Duration) (isDuration, fresh bool) {
	// decode HumanDuration from k8s.io/apimachinery/pkg/util/duration
	durationRegex := regexp.MustCompile(`^(?P<years>\d+y)?(?P<days>\d+d)?(?P<hours>\d+h)?(?P<minutes>\d+m)?(?P<seconds>\d+s)?$`)
	matches := durationRegex.FindStringSubmatch(value)
	if len(matches) > 0 {
		years := parseInt64(matches[1])
		days := parseInt64(matches[2])
		hours := parseInt64(matches[3])
		minutes := parseInt64(matches[4])
		seconds := parseInt64(matches[5])
		objAgeSeconds := years*365*24*3600 + days*24*3600 + hours*3600 + minutes*60 + seconds
		objAgeDuration := time.Duration(objAgeSeconds) * time.Second
		return true, objAgeDuration < threshold
	}
	return false, false
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}
