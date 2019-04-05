package helper

import (
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
		if cmd.Path("command").Data() == cmdArray[n] {
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

			if cmd.ExistsP("cli.execute") {
				cmdExec := cmd.Path("cli.execute").Data().(string)

				outFormat := "pretty"
				if cmd.ExistsP("cli.output") {
					outFormat = cmd.Path("cli.output").Data().(string)
				}

				if cmd.ExistsP("cli.param") {
					paramLength, err := strconv.Atoi(cmd.Path("cli.param").String())
					if err != nil {
						return nil, err
					}

					switch {
					case paramLength == 0:
						var cmdParam string

						for i := 1; i <= cmdLength-n; i++ {
							cmdParam = cmdParam + " " + cmdArray[n+i]
						}

						cmdExec = strings.Replace(cmdExec, "<0>", cmdParam, 1)
					case paramLength < cmdLength-n:
						return nil, errors.New("command: paramter ouf of bound")
					default:
						for i := 1; i <= paramLength; i++ {
							cmdExec = strings.Replace(cmdExec, "<"+strconv.Itoa(i)+">", cmdArray[n+i], 1)
						}
					}
				}

				out, err := exec.Command("sh", "-c", cmdExec).Output()
				if err != nil {
					return nil, err
				}

				outSplit := SplitAtChar(string(out), "\n", 2000, outFormat)

				if cmd.ExistsP("message") {
					var outMerge []string

					outMerge = append(outMerge, cmd.Path("message").Data().(string))
					outMerge = append(outMerge, outSplit...)

					return outMerge, nil
				}

				return outSplit, nil
			}

			if cmd.ExistsP("file") {
				out, err := ioutil.ReadFile(cmd.Path("file").Data().(string))
				if err != nil {
					return nil, err
				}

				outSplit := SplitAtChar(string(out), "\n", 2000, "normal")

				if cmd.ExistsP("message") {
					var outMerge []string

					outMerge = append(outMerge, cmd.Path("message").Data().(string))
					outMerge = append(outMerge, outSplit...)

					return outMerge, nil
				}

				return outSplit, nil
			}

			if cmd.ExistsP("message") {
				return []string{cmd.Path("message").Data().(string)}, nil
			}
		}
	}

	return nil, errors.New("command: command not found")
}
