package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

/* Global variables */
var Configuration = Config {}
var initialized int

// Info from config file
type Config struct {
	Username   string
	Password   string
}

/* Only this method should be used outside class*/
func GetConfig() Config{
	if(initialized==0) {
		readConfig()
	}
	return Configuration
}

// Reads info from config file
func readConfig() {
	var configfile = "res/config.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile, err)
	}

	config_data, err := toml.LoadFile(configfile)
	
	if err != nil {
    	log.Fatal("Error ", err.Error())
	} else {
		initialized = 1
		Configuration.Username = config_data.Get("sql-database.username").(string)
    	Configuration.Password = config_data.Get("sql-database.password").(string)
	}
}