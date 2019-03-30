package helper

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
)

var CMDList []*gabs.Container

func CMDParse(file string) ([]*gabs.Container, error) {
	json, err := gabs.ParseJSONFile(file)
	if err != nil {
		return nil, err
	}

	cmds, err := json.S("data").Children()
	if err != nil {
		return nil, err
	}

	return cmds, nil
}

func CMDExec(cmdList []*gabs.Container, cmdArray []string, n int) (interface{}, error) {
	cmdLength := len(cmdArray) - 1
	if n > cmdLength {
		return nil, errors.New("command exec: index out of bound")
	}

	for _, cmd := range cmdList {
		if cmd.Path("command").Data() == cmdArray[n] {
			if n < cmdLength && !cmd.ExistsP("param") {
				if cmd.ExistsP("data") {
					cmds, err := cmd.S("data").Children()
					if err != nil {
						return nil, err
					}

					return CMDExec(cmds, cmdArray, n+1)
				}

				return nil, errors.New("command exec: command not found")
			}

			if cmd.ExistsP("execute") {
				cmdExec := cmd.Path("execute").Data().(string)

				if cmd.ExistsP("param") {
					paramLength, err := strconv.Atoi(cmd.Path("param").String())
					if err != nil {
						return nil, err
					}

					switch {
					case paramLength == 0:
						paramLength = cmdLength - n
					case paramLength < cmdLength-n:
						return nil, errors.New("command exec: paramter ouf of bound")
					default:
					}

					for i := 1; i <= paramLength; i++ {
						cmdExec = strings.Replace(cmdExec, "<"+strconv.Itoa(i)+">", cmdArray[n+i], 1)
					}
				}

				out, err := exec.Command("sh", "-c", cmdExec).Output()
				if err != nil {
					return nil, err
				}

				if cmd.ExistsP("message") {
					return fmt.Sprintf("%v\n%v", cmd.Path("message").Data(), string(out)), nil
				}

				return string(out), nil
			}

			if cmd.ExistsP("message") {
				return cmd.Path("message").Data(), nil
			}
		}
	}

	return nil, errors.New("command exec: command not found")
}
