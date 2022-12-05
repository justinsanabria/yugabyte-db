/*
* Copyright (c) YugaByte, Inc.
 */

package preflight

import (
	"net"
	"strings"

	"github.com/yugabyte/yugabyte-db/managed/yba-installer/common"

	log "github.com/yugabyte/yugabyte-db/managed/yba-installer/logging"
)

// Actions that require root for us to execute are not possible under a non-root
// installation. We will simply add them to the preflight checks so that the user
// will know to perform these necessary installs prior to running the install command
// for YBA Installer themselves.

func checkPythonInstalled() {

	command := "bash"
	args := []string{"-c", "command -v python3"}
	output, _ := common.ExecuteBashCommand(command, args)

	outputTrimmed := strings.Replace(output, "\r\n", "", -1)

	if outputTrimmed == "" {

		log.Fatal("Python 3 not installed on the system, " +
			"please install Python 3 before continuing. As an example, " +
			"CentOS users can invoke the command \"sudo yum install -y " +
			"python3\"" + " in order to do so.")

	}

}

func checkUserLevelSystemdCapable() {

	if strings.Contains(common.DetectOS(), "CentOS") {

		log.Fatal("CentOS does not support systemd user mode, so you will " +
			"not be able to run a non-root install with systemd managed services." +
			"Please set the serviceManagementMode variable in yba-installer-input.yml " +
			"to be equal to \"crontab\" instead.")

	}

	currentUser := strings.ReplaceAll(strings.TrimSuffix(common.GetCurrentUser(), "\n"), " ", "")

	command := "bash"
	args := []string{"-c", "ls /var/lib/systemd/linger"}
	output, _ := common.ExecuteBashCommand(command, args)

	outputTrimmed := strings.Replace(output, "\r\n", "", -1)

	if !strings.Contains(outputTrimmed, currentUser) {

		log.Fatal("enable-linger has not been executed for the current user using " +
			"loginctl, so we you will not be able run a non-root install with systemd " +
			"managed services. Please run the command \"sudo loginctl enable-linger " +
			currentUser + " before continuing, or set the serviceManagementMode variable " +
			"in yba-installer-input.yml to be equal to \"crontab\" instead.")

	}

}

func checkBindToPort(port string) {

	_, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal("Unable to bind to TCP port " + port +
			", please make sure that the port " + port +
			" is available. Port " + port + " can be made available using the " +
			" command sudo firewall-cmd --zone=public --permanent --add-port=" + port + "/tcp.")
	}

}