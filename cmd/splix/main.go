package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_pgx_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool/pgx"
    core_http_middleware "github.com/korshunvladislav/splix/internal/core/transport/http/middleware"
    core_http_server "github.com/korshunvladislav/splix/internal/core/transport/http/server"
    groups_postgres_repository "github.com/korshunvladislav/splix/internal/features/groups/repository/postgres"
    groups_service "github.com/korshunvladislav/splix/internal/features/groups/service"
    groups_transport_http "github.com/korshunvladislav/splix/internal/features/groups/transport/http"
    users_postgres_repository "github.com/korshunvladislav/splix/internal/features/users/repository/postgres"
    users_service "github.com/korshunvladislav/splix/internal/features/users/service"
    users_transport_http "github.com/korshunvladislav/splix/internal/features/users/transport/http"
    "go.uber.org/zap"
)

var (
    timeZone = time.UTC
)

func main() {
    time.Local = timeZone

    ctx, cancel := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT, syscall.SIGTERM,
    )
    defer cancel()

    logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
    if err != nil {
        fmt.Println("failed to init application logger:", err)
        os.Exit(1)
    }
    defer logger.Close()

    logger.Debug("application time zone", zap.Any("zone", timeZone))

    logger.Debug("initializing postgres connection pool")
    pool, err := core_pgx_pool.NewPool(
        ctx,
        core_pgx_pool.NewConfigMust(),
    )
    if err != nil {
        logger.Fatal("failed to init postgres connection pool", zap.Error(err))
    }
    defer pool.Close()

    logger.Debug("initializing feature", zap.String("feature", "users"))
    usersRepository := users_postgres_repository.NewUsersRepository(pool)
    usersService := users_service.NewUsersService(usersRepository)
    usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

    logger.Debug("initializing feature", zap.String("feature", "groups"))
    groupsRepository := groups_postgres_repository.NewGroupsRepository(pool)
    groupsService := groups_service.NewGroupsService(groupsRepository)
    groupsTransportHTTP := groups_transport_http.NewGroupsHTTPHandler(groupsService)

    logger.Debug("initializing HTTP server")
    httpServer := core_http_server.NewHTTPServer(
        core_http_server.NewConfigMust(),
        logger,
        core_http_middleware.RequestID(),
        core_http_middleware.Logger(logger),
        core_http_middleware.Trace(),
        core_http_middleware.Panic(),
    )
    apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
    apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
    apiVersionRouterV1.RegisterRoutes(groupsTransportHTTP.Routes()...)

    /*
       Example of usage apiVersionRouterV2 with separate Middleware

       apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
          core_http_server.ApiVersion2,
          core_http_middleware.Dummy("api v2 middleware"),
       )
       apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)
    */

    httpServer.RegisterAPIRouters(
        apiVersionRouterV1,
        //apiVersionRouterV2,
    )

    if err := httpServer.Run(ctx); err != nil {
        logger.Error("HTTP server run error", zap.Error(err))
    }
}
