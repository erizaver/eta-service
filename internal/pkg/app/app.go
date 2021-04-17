package app

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	eta "github.com/erizaver/eta-service/internal/app/eta/v1"
	"github.com/erizaver/eta-service/internal/pkg/carsfacade"
	"github.com/erizaver/eta-service/internal/pkg/config"
	"github.com/erizaver/eta-service/internal/pkg/etaservice"
	cars "github.com/erizaver/eta-service/internal/pkg/externalclients/carservice/client"
	predict "github.com/erizaver/eta-service/internal/pkg/externalclients/predictservice/client"
	"github.com/erizaver/eta-service/internal/pkg/predictfacade"
	"github.com/erizaver/eta-service/internal/pkg/redisfacade"
	"github.com/erizaver/eta-service/pkg/api"
	"github.com/go-redis/redis"
)

type Application struct {
	//Config
	*config.Config

	//APIs
	EtaApi api.EtaServiceServer

	//Services
	EtaService eta.EtaService

	//Facades
	PredictFacade etaservice.PredictTimeFacade
	CarsFacade    etaservice.CarsFacade
	RedisFacade   etaservice.RedisFacade

	//Clients
	PredictClient *predict.PredictService
	CarsClient    *cars.CarsService

	//Cache
	RedisClient *redis.Client

	Mux        *runtime.ServeMux
	GrpcServer *grpc.Server
}

func NewApp() *Application {
	app := &Application{}

	log.Info("parsing config")
	app.Config = &config.Config{}

	envconfig.MustProcess("", app.Config)
	if app.Config.HttpPort == "" {
		log.Fatal("unable to parse config")
	}

	log.Info("initialising http clients")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = app.Config.RetryAmount

	standardHttpClient := retryClient.StandardClient()
	standardHttpClient.Timeout = app.Config.Timeout

	log.Info("initializing redis client")
	app.RedisClient = redis.NewClient(&redis.Options{
		Addr:     app.Config.RedisHost + app.Config.RedisPort,
		Password: "",
		DB:       0,
	})

	log.Info("initializing services")
	fakeEtaTransport := client.NewWithClient(app.Config.WheelyHost, app.Config.WheelyFakeEtaBasePath, nil, standardHttpClient)

	app.PredictClient = predict.New(fakeEtaTransport, strfmt.Default)
	app.CarsClient = cars.New(fakeEtaTransport, strfmt.Default)

	app.CarsFacade = carsfacade.NewCarsFacade(app.CarsClient, app.WheelyFakeCarServiceCarLimitPerRequest)
	app.PredictFacade = predictfacade.NewPredictFacade(app.PredictClient)
	app.RedisFacade = redisfacade.NewRedisFacade(app.RedisClient, app.Config.RedisCacheInvalidationTime, app.Config.RedisCacheSearchRadiusMeters)

	app.EtaService = etaservice.NewEtaService(app.CarsFacade, app.PredictFacade, app.RedisFacade)

	app.EtaApi = eta.NewEtaApi(app.EtaService)

	app.Mux = runtime.NewServeMux()
	app.GrpcServer = grpc.NewServer()

	return app
}

func (a *Application) Run(ctx context.Context) error {
	err := api.RegisterEtaServiceHandlerServer(ctx, a.Mux, a.EtaApi)
	if err != nil {
		return errors.Wrap(err, "unable to start server")
	}

	api.RegisterEtaServiceServer(a.GrpcServer, a.EtaApi)

	log.Info("listening to ", a.Config.HttpPort)
	return http.ListenAndServe(a.Config.HttpPort, a.Mux)
}
