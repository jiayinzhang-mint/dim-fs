package main

import (
	"dim-fs/protocol"
	"dim-fs/rpc"
	"dim-fs/utils"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// BuildEnv build mode (dev or prod)
var BuildEnv string

func startRPCServer() {
	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// register services
	protocol.RegisterCoreServiceServer(s, &rpc.CoreService{})
	reflection.Register(s)

	fmt.Println("DIMFs grpc service listening at " + viper.GetString("rpc.port"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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

// func startRestServer() {
// 	fmt.Println("shit")
// 	r := mux.NewRouter().StrictSlash(true)
// 	rest.InitImageAPI(r)

// 	log.Fatal(http.ListenAndServe(":"+viper.GetString("rest.port"), r))
// 	utils.LogInfo("server listening at " + viper.GetString("rest.port"))
// }

func main() {
	utils.LogInfo("Build Env: " + BuildEnv)
	prepareConfig()

	startRPCServer()
}
