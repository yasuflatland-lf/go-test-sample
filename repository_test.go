package main

import (
	"context"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
)

// Test MySQL Smoke
func TestMySQL(t *testing.T) {
	dsn, mysqlTerm := initMySQLContainer()
	defer mysqlTerm()

	t.Run("GORM Connection Open Test", func(t *testing.T) {

		ctx := context.Background()
		repository := NewRepository(ctx, dsn)

		if repository.DB() == nil {
			t.Error("DB is nil")
			return
		}
	})

	// Remove comments to generate model in the database automatically.
	//t.Run("Generate Models From Tables", func(t *testing.T) {
	//	seedDataPath, _ := os.Getwd()
	//	err = converter.NewTable2Struct().
	//		SavePath(seedDataPath + "/model.go").
	//		Dsn(dsn).
	//		Run()
	//})
}

// Generate dsn and close function for test use.
// *** DO NOT USE FOR PRODUCTION ***
func initMySQLContainer() (string, func()) {
	ctx := context.Background()
	username := "root"
	password := "password"
	seedDataPath, err := os.Getwd()
	mysqlPort, _ := nat.NewPort("tcp", "3306")
	req := testcontainers.ContainerRequest{
		Image:        "mysql:5.7",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
		},
		BindMounts: map[string]string{
			seedDataPath + "/test/db/mysql_init": "/docker-entrypoint-initdb.d",
			seedDataPath + "/test/db/my.cnf":     "/etc/mysql/conf.d/my.cnf",
		},
		WaitingFor: wait.ForListeningPort(mysqlPort),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	ip, err := mysqlC.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}

	// cloudSQL service fetch MySQL connection data from environment valuables.
	// Set here dummy server information for test purpose.
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_USERNAME", username)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_IP", ip)
	os.Setenv("DB_PORT", port.Port())

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/test", username, password, ip, port.Int())
	fmt.Println(dataSourceName)
	cTerm := func() {
		defer mysqlC.Terminate(ctx)
	}

	return dataSourceName, cTerm
}
