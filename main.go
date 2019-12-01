package main

import (
	"dim-fs/protocol"
	"dim-fs/service"
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

// StartServer start grpc server
func StartServer() {
	lis, err := net.Listen("tcp", ":"+viper.GetString("core.port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// register services
	protocol.RegisterCoreServiceServer(s, &service.CoreService{})
	reflection.Register(s)

	fmt.Println("DIMFs service listening at " + viper.GetString("core.port"))
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

func main() {
	utils.LogInfo("Build Env: " + BuildEnv)
	prepareConfig()
	StartServer()
}
