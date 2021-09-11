package configs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
	"net"
	"os"
	"path/filepath"
)

// Config stores configuration.
type Configure struct {
	RunMode     string   `envconfig:"RUN_MODE" default:"debug"`
	HTTPPort    string   `envconfig:"HTTP_PORT" default:"9999"`
	HTTPAddr    string   `envconfig:"HTTP_ADDR" default:"0.0.0.0"`
	LogLevel    string   `envconfig:"LOG_LEVEL" default:"debug"`
	IPWhiteList []string `envconfig:"IP_WHITE_LIST" default:"*"`
	RootPath    string   `envconfig:"ROOT_PATH" default:""`
}

var (
	TmpPath       = os.TempDir()
	Config        Configure
	RestClient    *kubernetes.Clientset
	KuBeResConf   *rest.Config
	CoreV1Client  *coreV1.CoreV1Client
	DynamicClient dynamic.Interface
	MetricsClient metrics.Interface
	envFile       = kingpin.Flag("env_file", "Load the environment variable file").Default(".envfile").String()
	rootPath      = kingpin.Flag("root_path", "Save data directory").Default("").String()
)

const notFoundKubeConfig = `Missing or incomplete kubernetes configuration info.  Please point to an existing, complete config file:

  1. Via the KUBECONFIG environment variable
  2. In your home directory as ~/.kube/config`

func Init(version string) {
	logs.Debug("Load variable")
	kingpin.Version(version)
	kingpin.Parse()
	// load environment variables from file.
	_ = godotenv.Load(*envFile)

	// load the configuration from the environment.
	err := envconfig.Process("", &Config)
	if err != nil {
		logs.Fatal(err)
	}
	if *rootPath != "" {
		Config.RootPath = *rootPath
	}
	logs.SetLogLevel(Config.LogLevel)
	logs.SetLogFormatter(&logrus.JSONFormatter{})
	if !utils.InSliceString("*", Config.IPWhiteList) {
		for _, ip := range Config.IPWhiteList {
			if net.ParseIP(ip) != nil {
				continue
			}
			logs.Fatal(fmt.Sprint(ip, ", Invalid whitelist"))
		}
	}
	KuBeResConf, err = kConfig()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
	RestClient, err = InitRestClient()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
	DynamicClient, err = InitDynamicClient()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
	MetricsClient, err = InitMetricsClient()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
	CoreV1Client, err = InitCoreV1Client()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
}

func kConfig() (conf *rest.Config, err error) {
	if Config.RunMode == gin.DebugMode {
		home, _ := homedir.Dir()
		var kubeConfigEnv = os.Getenv("KUBECONFIG")
		var kubeConfig = filepath.Join(home, ".kube", "config")
		var kuBeConf []byte
		if utils.FileOrPathExist(kubeConfigEnv) {
			kuBeConf, err = ioutil.ReadFile(kubeConfigEnv)
		} else if utils.FileOrPathExist(kubeConfig) {
			kuBeConf, err = ioutil.ReadFile(kubeConfig)
		}
		if err != nil {
			return nil, err
		}
		conf, err = clientcmd.RESTConfigFromKubeConfig(kuBeConf)
	} else {
		conf, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	conf.Timeout = 0
	return
}

func InitRestClient() (*kubernetes.Clientset, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kConf)
}

func InitDynamicClient() (dynamic.Interface, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(kConf)
}

func InitMetricsClient() (metrics.Interface, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	return metrics.NewForConfig(kConf)
}

func InitCoreV1Client() (*coreV1.CoreV1Client, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	return coreV1.NewForConfig(kConf)
}
