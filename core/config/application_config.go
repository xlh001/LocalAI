package config

import (
	"context"
	"embed"
	"encoding/json"
	"time"

	"github.com/go-skynet/LocalAI/pkg/gallery"
	"github.com/go-skynet/LocalAI/pkg/xsysinfo"
	"github.com/rs/zerolog/log"
)

type ApplicationConfig struct {
	Context                             context.Context
	ConfigFile                          string
	ModelPath                           string
	UploadLimitMB, Threads, ContextSize int
	DisableWebUI                        bool
	F16                                 bool
	Debug                               bool
	ImageDir                            string
	AudioDir                            string
	UploadDir                           string
	ConfigsDir                          string
	DynamicConfigsDir                   string
	DynamicConfigsDirPollInterval       time.Duration
	CORS                                bool
	CSRF                                bool
	PreloadJSONModels                   string
	PreloadModelsFromPath               string
	CORSAllowOrigins                    string
	ApiKeys                             []string
	OpaqueErrors                        bool

	ModelLibraryURL string

	Galleries []gallery.Gallery

	BackendAssets     embed.FS
	AssetsDestination string

	ExternalGRPCBackends map[string]string

	AutoloadGalleries bool

	SingleBackend           bool
	ParallelBackendRequests bool

	WatchDogIdle bool
	WatchDogBusy bool
	WatchDog     bool

	ModelsURL []string

	WatchDogBusyTimeout, WatchDogIdleTimeout time.Duration
}

type AppOption func(*ApplicationConfig)

func NewApplicationConfig(o ...AppOption) *ApplicationConfig {
	opt := &ApplicationConfig{
		Context:       context.Background(),
		UploadLimitMB: 15,
		ContextSize:   512,
		Debug:         true,
	}
	for _, oo := range o {
		oo(opt)
	}
	return opt
}

func WithModelsURL(urls ...string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelsURL = urls
	}
}

func WithModelPath(path string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelPath = path
	}
}

func WithCors(b bool) AppOption {
	return func(o *ApplicationConfig) {
		o.CORS = b
	}
}

func WithCsrf(b bool) AppOption {
	return func(o *ApplicationConfig) {
		o.CSRF = b
	}
}

func WithModelLibraryURL(url string) AppOption {
	return func(o *ApplicationConfig) {
		o.ModelLibraryURL = url
	}
}

var EnableWatchDog = func(o *ApplicationConfig) {
	o.WatchDog = true
}

var EnableWatchDogIdleCheck = func(o *ApplicationConfig) {
	o.WatchDog = true
	o.WatchDogIdle = true
}

var EnableWatchDogBusyCheck = func(o *ApplicationConfig) {
	o.WatchDog = true
	o.WatchDogBusy = true
}

var DisableWebUI = func(o *ApplicationConfig) {
	o.DisableWebUI = true
}

func SetWatchDogBusyTimeout(t time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.WatchDogBusyTimeout = t
	}
}

func SetWatchDogIdleTimeout(t time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.WatchDogIdleTimeout = t
	}
}

var EnableSingleBackend = func(o *ApplicationConfig) {
	o.SingleBackend = true
}

var EnableParallelBackendRequests = func(o *ApplicationConfig) {
	o.ParallelBackendRequests = true
}

var EnableGalleriesAutoload = func(o *ApplicationConfig) {
	o.AutoloadGalleries = true
}

func WithExternalBackend(name string, uri string) AppOption {
	return func(o *ApplicationConfig) {
		if o.ExternalGRPCBackends == nil {
			o.ExternalGRPCBackends = make(map[string]string)
		}
		o.ExternalGRPCBackends[name] = uri
	}
}

func WithCorsAllowOrigins(b string) AppOption {
	return func(o *ApplicationConfig) {
		o.CORSAllowOrigins = b
	}
}

func WithBackendAssetsOutput(out string) AppOption {
	return func(o *ApplicationConfig) {
		o.AssetsDestination = out
	}
}

func WithBackendAssets(f embed.FS) AppOption {
	return func(o *ApplicationConfig) {
		o.BackendAssets = f
	}
}

func WithStringGalleries(galls string) AppOption {
	return func(o *ApplicationConfig) {
		if galls == "" {
			o.Galleries = []gallery.Gallery{}
			return
		}
		var galleries []gallery.Gallery
		if err := json.Unmarshal([]byte(galls), &galleries); err != nil {
			log.Error().Err(err).Msg("failed loading galleries")
		}
		o.Galleries = append(o.Galleries, galleries...)
	}
}

func WithGalleries(galleries []gallery.Gallery) AppOption {
	return func(o *ApplicationConfig) {
		o.Galleries = append(o.Galleries, galleries...)
	}
}

func WithContext(ctx context.Context) AppOption {
	return func(o *ApplicationConfig) {
		o.Context = ctx
	}
}

func WithYAMLConfigPreload(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.PreloadModelsFromPath = configFile
	}
}

func WithJSONStringPreload(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.PreloadJSONModels = configFile
	}
}
func WithConfigFile(configFile string) AppOption {
	return func(o *ApplicationConfig) {
		o.ConfigFile = configFile
	}
}

func WithUploadLimitMB(limit int) AppOption {
	return func(o *ApplicationConfig) {
		o.UploadLimitMB = limit
	}
}

func WithThreads(threads int) AppOption {
	return func(o *ApplicationConfig) {
		if threads == 0 { // 0 is not allowed
			threads = xsysinfo.CPUPhysicalCores()
		}
		o.Threads = threads
	}
}

func WithContextSize(ctxSize int) AppOption {
	return func(o *ApplicationConfig) {
		o.ContextSize = ctxSize
	}
}

func WithF16(f16 bool) AppOption {
	return func(o *ApplicationConfig) {
		o.F16 = f16
	}
}

func WithDebug(debug bool) AppOption {
	return func(o *ApplicationConfig) {
		o.Debug = debug
	}
}

func WithAudioDir(audioDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.AudioDir = audioDir
	}
}

func WithImageDir(imageDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.ImageDir = imageDir
	}
}

func WithUploadDir(uploadDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.UploadDir = uploadDir
	}
}

func WithConfigsDir(configsDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.ConfigsDir = configsDir
	}
}

func WithDynamicConfigDir(dynamicConfigsDir string) AppOption {
	return func(o *ApplicationConfig) {
		o.DynamicConfigsDir = dynamicConfigsDir
	}
}

func WithDynamicConfigDirPollInterval(interval time.Duration) AppOption {
	return func(o *ApplicationConfig) {
		o.DynamicConfigsDirPollInterval = interval
	}
}

func WithApiKeys(apiKeys []string) AppOption {
	return func(o *ApplicationConfig) {
		o.ApiKeys = apiKeys
	}
}

func WithOpaqueErrors(opaque bool) AppOption {
	return func(o *ApplicationConfig) {
		o.OpaqueErrors = opaque
	}
}

// ToConfigLoaderOptions returns a slice of ConfigLoader Option.
// Some options defined at the application level are going to be passed as defaults for
// all the configuration for the models.
// This includes for instance the context size or the number of threads.
// If a model doesn't set configs directly to the config model file
// it will use the defaults defined here.
func (o *ApplicationConfig) ToConfigLoaderOptions() []ConfigLoaderOption {
	return []ConfigLoaderOption{
		LoadOptionContextSize(o.ContextSize),
		LoadOptionDebug(o.Debug),
		LoadOptionF16(o.F16),
		LoadOptionThreads(o.Threads),
	}
}

// func WithMetrics(meter *metrics.Metrics) AppOption {
// 	return func(o *StartupOptions) {
// 		o.Metrics = meter
// 	}
// }
