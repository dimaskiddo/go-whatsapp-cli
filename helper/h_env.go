package helper

import (
	"log"
	"os"
	"strconv"
)

func GetEnv(name string, datatype string, required bool) interface{} {
	var val interface{}
	var err error

	env := os.Getenv(name)

	if len(env) == 0 && required {
		log.Fatalln("environment: variable '" + name + "' is required")
	}

	if len(env) != 0 {
		switch datatype {
		case "string":
			val = env
		case "bool":
			val, err = strconv.ParseBool(env)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid bool datatype")
			}
		case "int":
			val, err = strconv.ParseInt(env, 0, 64)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid integer datatype")
			}

			val = int(val.(int64))
		case "int32":
			val, err = strconv.ParseInt(env, 0, 64)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid integer datatype")
			}

			val = int32(val.(int64))
		case "int64":
			val, err = strconv.ParseInt(env, 0, 64)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid integer datatype")
			}
		case "float32":
			val, err = strconv.ParseFloat(env, 64)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid float datatype")
			}

			val = float32(val.(float64))
		case "float64":
			val, err = strconv.ParseFloat(env, 64)
			if err != nil {
				log.Fatalln("environment: variable '" + name + "' has invalid float datatype")
			}
		default:
			log.Fatalln("environment: requested unknown datatype")
		}
	} else {
		return nil
	}

	return val
}
