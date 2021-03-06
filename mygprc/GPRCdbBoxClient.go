package mygprc


import (
	pb_device "github.com/onezerobinary/db-box/proto/device"
	"github.com/goinggo/tracelog"
	"os"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"fmt"
	"github.com/spf13/viper"
)


func StartGRPCConnection() (connection *grpc.ClientConn){

	// Get info from production or local
	DBaddress := os.Getenv("DB_ADDRESS")

	if len(DBaddress) == 0 {
		DBaddress = viper.GetString("service.db-box")
		tracelog.Warning("GRPCdbBoxClient", "StartGRPCConnection", "####### Development #########")
	}

	// set up connection to the gRPC server
	conn, err := grpc.Dial(DBaddress, grpc.WithInsecure())
	if err != nil {
		tracelog.Errorf(err, "GRPCdbBoxClient", "StartGRPCConnection", "Did not open the connection")
		os.Exit(1)
	}

	return conn
}

func StopGRPCConnection(connection *grpc.ClientConn){
	// set up connection to the gRPC server
	err := connection.Close()

	if err != nil {
		tracelog.Errorf(err, "GRPCdbBoxClient", "StopGRPCConnection", "Did not close the connection")
		os.Exit(1)
	}
}


func GetExpoPushTokensByGeoHash(geoHash *pb_device.GeoHash) (*pb_device.ExpoPushTokens) {
	tracelog.Trace("GPRCdbBoxClient","GetExpoPushTokensByGeoHash","Search devices in a nearest area")

	conn := StartGRPCConnection()
	defer StopGRPCConnection(conn)

	client := pb_device.NewDeviceServiceClient(conn)

	expoPushTokens, _ := client.GetExpoPushTokensByGeoHash(context.Background(), geoHash)

	fmt.Println("Token: " + string(len(expoPushTokens.Token)))

	tracelog.Warning("GPRCdbBoxClient","GetExpoPushTokensByGeoHash","Empty")

	if len(expoPushTokens.Token) == 0 {
		expoPushTokens.Token = []string{}
	}

	return expoPushTokens
}