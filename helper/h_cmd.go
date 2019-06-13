package helper

import (
	"bytes"
	"errors"
	"io/ioutil"
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

func CMDExec(cmdList []*gabs.Container, cmdArray []string, n int) ([]string, error) {
	if cmdList == nil {
		return nil, errors.New("command: empty command list")
	}

	cmdLength := len(cmdArray) - 1
	if n > cmdLength {
		return nil, errors.New("command: index out of bound")
	}

	for _, cmd := range cmdList {
		cmdTypeFound := false

		if cmd.ExistsP("type") {
			if cmd.Path("type").Data().(string) == "" {
				return nil, errors.New("command: invalid command type")
			} else {
				cmdTypeLists, err := cmd.S("command").Children()
				if err != nil {
					return nil, err
				}

				for _, cmd := range cmdTypeLists {
					if cmd.Path("name").Data() == cmdArray[n] {
						cmdTypeFound = true
					}
				}

				if !cmdTypeFound {
					return nil, errors.New("command: command not found in type")
				}
			}
		}

		if cmd.Path("command").Data() == cmdArray[n] || cmdTypeFound {
			if n < cmdLength && !cmd.ExistsP("cli.param") {
				if cmd.ExistsP("data") {
					cmds, err := cmd.S("data").Children()
					if err != nil {
						return nil, err
					}

					return CMDExec(cmds, cmdArray, n+1)
				}

				return nil, errors.New("command: command not found")
			}

			if cmd.ExistsP("curl.url") {
				cmdURL := cmd.Path("curl.url").Data().(string)

				cmdMethod := "-X GET "
				if cmd.ExistsP("curl.method") {
					cmdMethod = "-X " + cmd.Path("curl.method").Data().(string) + " "
				}

				cmdHeader := "-H 'cache-control: no-cache' "
				if cmd.ExistsP("curl.header") {
					headers, err := cmd.S("curl.header").Children()
					if err != nil {
						return nil, err
					}

					for _, header := range headers {
						cmdHeader = cmdHeader + "-H '" + header.Data().(string) + "' "
					}
				}

				cmdForm := ""
				if cmd.ExistsP("curl.form") {
					forms, err := cmd.S("curl.form").Children()
					if err != nil {
						return nil, err
					}

					for _, form := range forms {
						cmdForm = cmdForm + "-F '" + form.Data().(string) + "' "
					}
				}

				cmdBody := ""
				if cmd.ExistsP("curl.body") {
					cmdBody = "-d " + cmd.Path("curl.body").Data().(string) + " "
				}

				cmdTrim := true
				if cmd.ExistsP("curl.trim") {
					cmdTrim = cmd.Path("curl.trim").Data().(bool)
				}

				cmdOutput := true
				if cmd.ExistsP("curl.pretty") {
					cmdOutput = cmd.Path("curl.pretty").Data().(bool)
				}

				cmdCURL := "curl " + cmdMethod + cmdHeader + cmdForm + cmdBody + cmdURL

				cmdExecSplit := SplitWithEscapeN(cmdCURL, " ", -1, cmdTrim)
				execOutput, err := exec.Command(cmdExecSplit[0], cmdExecSplit[1:]...).Output()
				if err != nil {
					return nil, err
				}

				outReturn := []string{"There is nothing here, but the request is success 😆"}
				if len(string(execOutput)) != 0 {
					outReturn = SplitAfterCharN(string(execOutput), "\n", 2000, -1, cmdOutput, cmdTrim)
				}

				if cmd.ExistsP("message") {
					outReturn[0] = cmd.Path("message").Data().(string) + "\n" + outReturn[0]
				}

				return outReturn, nil
			}

			if cmd.ExistsP("cli.execute") {
				cmdExec := cmd.Path("cli.execute").Data().(string)

				if cmd.ExistsP("cli.param") {
					cmdParamLength, err := strconv.Atoi(cmd.Path("cli.param").String())
					if err != nil {
						return nil, err
					}

					switch {
					case cmdParamLength == 0:
						cmdParam := ""

						for i := 1; i <= cmdLength-n; i++ {
							cmdParam = cmdParam + " " + cmdArray[n+i]
						}

						cmdExec = strings.Replace(cmdExec, "<0>", cmdParam, 1)
					case cmdParamLength < cmdLength-n:
						return nil, errors.New("command: paramter ouf of bound")
					default:
						for i := 1; i <= cmdParamLength; i++ {
							cmdExec = strings.Replace(cmdExec, "<"+strconv.Itoa(i)+">", cmdArray[n+i], 1)
						}
					}
				}

				cmdTrim := true
				if cmd.ExistsP("cli.trim") {
					cmdTrim = cmd.Path("cli.trim").Data().(bool)
				}

				cmdOutput := true
				if cmd.ExistsP("cli.pretty") {
					cmdOutput = cmd.Path("cli.pretty").Data().(bool)
				}

				cmdExecSplit := SplitWithEscapeN(cmdExec, " ", -1, cmdTrim)
				cmdRun := exec.Command(cmdExecSplit[0], cmdExecSplit[1:]...)

				var cmdStdout bytes.Buffer
				var cmdStderr bytes.Buffer

				cmdRun.Stdout = &cmdStdout
				cmdRun.Stderr = &cmdStderr

				err := cmdRun.Run()
				if err != nil {
					outReturn := SplitAfterCharN(string(cmdStderr.String()), "\n", 2000, -1, cmdOutput, cmdTrim)
					return outReturn, err
				}
				execOutput := cmdStdout.String()

				outReturn := []string{"There is nothing here, but the execution is success 😆"}
				if len(string(execOutput)) != 0 {
					outReturn = SplitAfterCharN(string(execOutput), "\n", 2000, -1, cmdOutput, cmdTrim)
				}

				if len(string(execOutput)) != 0 && cmd.ExistsP("message") {
					outReturn[0] = cmd.Path("message").Data().(string) + "\n" + outReturn[0]
				}

				return outReturn, nil
			}

			if cmd.ExistsP("file") {
				execOutput, err := ioutil.ReadFile(cmd.Path("file").Data().(string))
				if err != nil {
					return nil, err
				}

				outReturn := []string{"Sorry, i got nothing from the file 😔"}
				if len(string(execOutput)) != 0 {
					outReturn = SplitAfterCharN(string(execOutput), "\n", 2000, -1, false, true)
				}

				if len(string(execOutput)) != 0 && cmd.ExistsP("message") {
					outReturn[0] = cmd.Path("message").Data().(string) + "\n" + outReturn[0]
				}

				return outReturn, nil
			}

			if cmd.ExistsP("message") {
				return []string{cmd.Path("message").Data().(string)}, nil
			}
		}
	}

	return nil, errors.New("command: command not found")
}
