package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"hyper/engine"
	"hyper/types"

	gflag "github.com/jessevdk/go-flags"
)

// We need to process the POD json data with the given file
func (cli *HyperClient) HyperCmdPod(args ...string) error {
	t1 := time.Now()
	var parser = gflag.NewParser(nil, gflag.Default)
	parser.Usage = "pod POD_FILE\n\nCreate a pod, initialize a pod and run it"
	args, err := parser.Parse()
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			return err
		} else {
			return nil
		}
	}
	if len(args) < 2 {
		return fmt.Errorf("\"pod\" requires a minimum of 1 argument, please provide POD spec file.\n")
	}
	jsonFile := args[1]
	if _, err := os.Stat(jsonFile); err != nil {
		return err
	}

	jsonbody, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	podId, err := cli.RunPod(string(jsonbody))
	if err != nil {
		return err
	}
	fmt.Printf("POD id is %s\n", podId)
	t2 := time.Now()
	fmt.Printf("Time to run a POD is %d ms\n", (t2.UnixNano()-t1.UnixNano())/1000000)

	return nil
}

func (cli *HyperClient) CreatePod(jsonbody string) (string, error) {
	v := url.Values{}
	v.Set("podArgs", jsonbody)
	body, _, err := readBody(cli.call("POST", "/pod/create?"+v.Encode(), nil, nil))
	if err != nil {
		return "", err
	}
	out := engine.NewOutput()
	remoteInfo, err := out.AddEnv()
	if err != nil {
		return "", err
	}

	if _, err := out.Write(body); err != nil {
		return "", fmt.Errorf("Error reading remote info: %s", err)
	}
	out.Close()
	errCode := remoteInfo.GetInt("Code")
	if errCode == types.E_OK {
		//fmt.Println("VM is successful to start!")
	} else {
		// case types.E_CONTEXT_INIT_FAIL:
		// case types.E_DEVICE_FAIL:
		// case types.E_QMP_INIT_FAIL:
		// case types.E_QMP_COMMAND_FAIL:
		if errCode != types.E_BAD_REQUEST &&
			errCode != types.E_FAILED {
			return "", fmt.Errorf("Error code is %d", errCode)
		} else {
			return "", fmt.Errorf("Cause is %s", remoteInfo.Get("Cause"))
		}
	}
	return remoteInfo.Get("ID"), nil
}

func (cli *HyperClient) HyperCmdStart(args ...string) error {
	var opts struct {
		//		OnlyVm    bool     `long:"onlyvm" default:"false" value-name:"false" description:"Only start a new VM"`
		Cpu int `short:"c" long:"cpu" default:"1" value-name:"1" description:"CPU number for the VM"`
		Mem int `short:"m" long:"memory" default:"128" value-name:"128" description:"Memory size (MB) for the VM"`
	}
	var parser = gflag.NewParser(&opts, gflag.Default)
	parser.Usage = "start [-c 1 -m 128]| POD_ID \n\nlaunch a 'pending' pod"
	args, err := parser.Parse()
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			return err
		} else {
			return nil
		}
	}
	if false {
		// Only run a new VM
		v := url.Values{}
		v.Set("cpu", fmt.Sprintf("%d", opts.Cpu))
		v.Set("mem", fmt.Sprintf("%d", opts.Mem))
		body, _, err := readBody(cli.call("POST", "/vm/create?"+v.Encode(), nil, nil))
		if err != nil {
			return err
		}
		out := engine.NewOutput()
		remoteInfo, err := out.AddEnv()
		if err != nil {
			return err
		}

		if _, err := out.Write(body); err != nil {
			return fmt.Errorf("Error reading remote info: %s", err)
		}
		out.Close()
		errCode := remoteInfo.GetInt("Code")
		if errCode == types.E_OK {
			//fmt.Println("VM is successful to start!")
		} else {
			// case types.E_CONTEXT_INIT_FAIL:
			// case types.E_DEVICE_FAIL:
			// case types.E_QMP_INIT_FAIL:
			// case types.E_QMP_COMMAND_FAIL:
			if errCode != types.E_BAD_REQUEST &&
				errCode != types.E_FAILED {
				return fmt.Errorf("Error code is %d", errCode)
			} else {
				return fmt.Errorf("Cause is %s", remoteInfo.Get("Cause"))
			}
		}
		fmt.Printf("New VM id is %s\n", remoteInfo.Get("ID"))
		return nil
	}
	if len(args) < 2 {
		return fmt.Errorf("\"start\" requires a minimum of 1 argument, please provide POD ID.\n")
	}
	var (
		podId string
		vmId  string
	)
	if len(args) == 2 {
		podId = args[1]
	}
	if len(args) == 3 {
		podId = args[1]
		vmId = args[2]
	}
	// fmt.Printf("Pod ID is %s, VM ID is %s\n", podId, vmId)
	_, err = cli.StartPod(podId, vmId)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully started the Pod(%s)\n", podId)
	return nil
}

func (cli *HyperClient) StartPod(podId, vmId string) (string, error) {
	v := url.Values{}
	v.Set("podId", podId)
	v.Set("vmId", vmId)
	body, _, err := readBody(cli.call("POST", "/pod/start?"+v.Encode(), nil, nil))
	if err != nil {
		return "", err
	}
	out := engine.NewOutput()
	remoteInfo, err := out.AddEnv()
	if err != nil {
		return "", err
	}

	if _, err := out.Write(body); err != nil {
		return "", fmt.Errorf("Error reading remote info: %s", err)
	}
	out.Close()
	errCode := remoteInfo.GetInt("Code")
	if errCode == types.E_OK {
		//fmt.Println("VM is successful to start!")
	} else {
		// case types.E_CONTEXT_INIT_FAIL:
		// case types.E_DEVICE_FAIL:
		// case types.E_QMP_INIT_FAIL:
		// case types.E_QMP_COMMAND_FAIL:
		if errCode != types.E_BAD_REQUEST &&
			errCode != types.E_FAILED {
			return "", fmt.Errorf("Error code is %d", errCode)
		} else {
			return "", fmt.Errorf("Cause is %s", remoteInfo.Get("Cause"))
		}
	}
	return remoteInfo.Get("ID"), nil
}

func (cli *HyperClient) RunPod(podstring string) (string, error) {
	v := url.Values{}
	v.Set("podArgs", podstring)
	body, _, err := readBody(cli.call("POST", "/pod/run?"+v.Encode(), nil, nil))
	if err != nil {
		return "", err
	}
	out := engine.NewOutput()
	remoteInfo, err := out.AddEnv()
	if err != nil {
		return "", err
	}

	if _, err := out.Write(body); err != nil {
		return "", fmt.Errorf("Error reading remote info: %s", err)
	}
	out.Close()
	errCode := remoteInfo.GetInt("Code")
	if errCode == types.E_OK {
		//fmt.Println("VM is successful to start!")
	} else {
		// case types.E_CONTEXT_INIT_FAIL:
		// case types.E_DEVICE_FAIL:
		// case types.E_QMP_INIT_FAIL:
		// case types.E_QMP_COMMAND_FAIL:
		if errCode != types.E_BAD_REQUEST &&
			errCode != types.E_FAILED {
			return "", fmt.Errorf("Error code is %d", errCode)
		} else {
			return "", fmt.Errorf("Cause is %s", remoteInfo.Get("Cause"))
		}
	}
	return remoteInfo.Get("ID"), nil
}
