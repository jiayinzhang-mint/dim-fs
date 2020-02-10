package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/insdim/dim-fs/protocol"
	"github.com/insdim/dim-fs/rest"
	"github.com/insdim/dim-fs/rpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// BuildEnv build mode (dev or prod)
var BuildEnv string

var (
	g errgroup.Group
)

func startRPCServer() error {
	if viper.GetString("rpc.port") == "" {
		panic(fmt.Errorf("dim-fs GRPC service port undefined"))
	}

	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		logrus.Fatalf("dim-fs GRPC service failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// register services
	protocol.RegisterCoreServiceServer(s, &rpc.CoreService{})
	reflection.Register(s)

	logrus.Info("dim-fs GRPC service listening at " + viper.GetString("rpc.port"))
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("dim-fs GRPC service failed to serve: %v", err)
	}

	return err
}

func prepareConfig() {
	// Get config from json
	if BuildEnv == "prod" {
		viper.SetConfigName("config.prod")
	} else {
		viper.SetConfigName("config.test")
	}
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("config file error")
		os.Exit(1)
	}
}

func startRestServer() error {
	if viper.GetString("rest.port") == "" {
		panic(fmt.Errorf("dim-fs REST service port undefined"))
	}

	r := mux.NewRouter().StrictSlash(true)
	rest.InitImageAPI(r)
	logrus.Info("dim-fs REST service listening at " + viper.GetString("rest.port"))

	err := http.ListenAndServe(":"+viper.GetString("rest.port"), r)

	return err
}

func main() {
	logrus.Info("Build Env: " + BuildEnv)
	prepareConfig()

	g.Go(func() error {
		return startRestServer()
	})

	g.Go(func() error {
		return startRPCServer()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
