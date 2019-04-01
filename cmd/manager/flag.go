package main

import (
	"flag"

	"github.com/giantswarm/api/service"
	"github.com/giantswarm/microerror"
	microflag "github.com/giantswarm/microkit/flag"
)

var (
	description string     = "The aws-operator handles Kubernetes clusters running on a Kubernetes cluster inside of AWS."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "aws-operator"
	source      string     = "https://github.com/giantswarm/aws-operator"
)

func todo(fs *flag.FlagSet) error {
	{
		fs.StringSlice(f.Config.Dirs, []string{"."}, "List of config file directories.")
		fs.StringSlice(f.Config.Files, []string{"config"}, "List of the config file names. All viper supported extensions can be used.")

		fs.Bool(f.Server.Enable.Debug.Server, false, "Enable debug server at http://127.0.0.1:6060/debug.")

		fs.String(f.Server.Listen.Address, "http://127.0.0.1:8000", "Address used to make the server listen to.")
		fs.String(f.Server.Listen.MetricsAddress, "", "Optional alternate address to expose metrics on at /metrics. Leave blank to use the default server (listen address above).")

		fs.Bool(f.Server.Log.Access, false, "Whether to emit logs for each requested route.")

		fs.String(f.Server.TLS.CaFile, "", "File path of the TLS root CA file, if any.")
		fs.String(f.Server.TLS.CrtFile, "", "File path of the TLS public key file, if any.")
		fs.String(f.Server.TLS.KeyFile, "", "File path of the TLS private key file, if any.")

		fs.String(f.Service.Kubernetes.Address, "http://127.0.0.1:6443", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
		fs.Bool(f.Service.Kubernetes.InCluster, false, "Whether to use the in-cluster config to authenticate with Kubernetes.")
		fs.String(f.Service.Kubernetes.KubeConfig, "", "KubeConfig used to connect to Kubernetes. When empty other settings are used.")
		fs.String(f.Service.Kubernetes.TLS.CAFile, "", "Certificate authority file path to use to authenticate with Kubernetes.")
		fs.String(f.Service.Kubernetes.TLS.CrtFile, "", "Certificate file path to use to authenticate with Kubernetes.")
		fs.String(f.Service.Kubernetes.TLS.KeyFile, "", "Key file path to use to authenticate with Kubernetes.")
	}

	{
		v := viper.New()

		// We have to parse the flags given via command line first. Only that way we
		// are able to use the flag configuration for the location of configuration
		// directories and files in the next step below.
		microflag.Parse(v, fs)

		// Merge the given command line flags with the given environment variables and
		// the given config files, if any. The merged flags will be applied to the
		// given viper.
		err := microflag.Merge(v, fs, v.GetStringSlice(f.Config.Dirs), v.GetStringSlice(f.Config.Files))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		var newService *service.Service
		{
			c := service.Config{
				Flag:   f,
				Logger: l,
				Viper:  v,

				Description: description,
				GitCommit:   gitCommit,
				ProjectName: name,
				Source:      source,
			}

			newService, err = service.New(c)
			if err != nil {
				return microerror.Mask(err)
			}

			go newService.Boot(ctx)
		}
	}

	return nil
}
