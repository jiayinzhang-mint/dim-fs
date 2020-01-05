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
	"github.com/insdim/dim-fs/utils"
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
		panic(fmt.Errorf("GRPC service port undefined"))
	}

	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// register services
	protocol.RegisterCoreServiceServer(s, &rpc.CoreService{})
	reflection.Register(s)

	utils.LogInfo("DIMFs GRPC service listening at " + viper.GetString("rpc.port"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
		utils.LogError("config file error")
		os.Exit(1)
	}
}

func startRestServer() error {
	if viper.GetString("rest.port") == "" {
		panic(fmt.Errorf("REST service port undefined"))
	}

	r := mux.NewRouter().StrictSlash(true)
	rest.InitImageAPI(r)
	utils.LogInfo("DIMFs REST service listening at " + viper.GetString("rest.port"))

	err := http.ListenAndServe(":"+viper.GetString("rest.port"), r)

	return err
}

func main() {
	utils.LogInfo("Build Env: " + BuildEnv)
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
