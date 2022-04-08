package serverctl

import (
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gravitl/netmaker/database"
	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/logic"
	"github.com/gravitl/netmaker/logic/acls"
	"github.com/gravitl/netmaker/logic/acls/nodeacls"
	"github.com/gravitl/netmaker/netclient/ncutils"
	"github.com/gravitl/netmaker/servercfg"
)

// COMMS_NETID - name of the comms network
var COMMS_NETID string

const (
	// NETMAKER_BINARY_NAME - name of netmaker binary
	NETMAKER_BINARY_NAME = "netmaker"
)

func gracefulCommsWait() {
	output, _ := ncutils.RunCmd("wg", false)
	starttime := time.Now()
	ifaceReady := strings.Contains(output, COMMS_NETID)
	for !ifaceReady && !(time.Now().After(starttime.Add(time.Second << 4))) {
		output, _ = ncutils.RunCmd("wg", false)
		SyncServerNetwork(COMMS_NETID)
		time.Sleep(time.Second)
		ifaceReady = strings.Contains(output, COMMS_NETID)
	}
	logger.Log(1, "comms network", COMMS_NETID, "ready")
}

// InitServerNetclient - intializes the server netclient
// 1. Check if config directory exists, if not attempt to make
// 2. Check current networks and run pull to get interface up to date in case of restart
func InitServerNetclient() error {
	netclientDir := ncutils.GetNetclientPath()
	_, err := os.Stat(netclientDir + "/config")
	if os.IsNotExist(err) {
		os.MkdirAll(netclientDir+"/config", 0700)
	} else if err != nil {
		logger.Log(1, "could not find or create", netclientDir)
		return err
	}

	var networks, netsErr = logic.GetNetworks()
	if netsErr == nil || database.IsEmptyRecord(netsErr) {
		for _, network := range networks {
			var currentServerNode, nodeErr = logic.GetNetworkServerLocal(network.NetID)
			if nodeErr == nil {
				if currentServerNode.Version != servercfg.Version {
					currentServerNode.Version = servercfg.Version
					logic.UpdateNode(&currentServerNode, &currentServerNode)
				}
				if err = logic.ServerPull(&currentServerNode, true); err != nil {
					logger.Log(1, "failed pull for network", network.NetID, ", on server node", currentServerNode.ID)
				}
			}
		}
	}

	return nil
}

// SyncServerNetwork - ensures a wg interface and node exists for server
func SyncServerNetwork(network string) error {
	serverNetworkSettings, err := logic.GetNetwork(network)
	if err != nil {
		return err
	}
	localnets, err := net.Interfaces()
	if err != nil {
		return err
	}

	ifaceExists := false
	for _, localnet := range localnets {
		if serverNetworkSettings.DefaultInterface == localnet.Name {
			ifaceExists = true
		}
	}

	serverNode, err := logic.GetNetworkServerLocal(network)
	if !ifaceExists && (err == nil && serverNode.ID != "") {
		return logic.ServerUpdate(&serverNode, true)
	} else if !ifaceExists {
		_, err := logic.ServerJoin(&serverNetworkSettings)
		if err != nil {
			if err == nil {
				err = errors.New("network add failed for " + serverNetworkSettings.NetID)
			}
			if !strings.Contains(err.Error(), "macaddress_unique") { // ignore macaddress unique error throws
				logger.Log(1, "error adding network", serverNetworkSettings.NetID, "during sync:", err.Error())
			}
		}
	}

	// remove networks locally that do not exist in database
	/*
		for _, localnet := range localnets {
			if strings.Contains(localnet.Name, "nm-") {
				var exists = ""
				if serverNetworkSettings.DefaultInterface == localnet.Name {
					exists = serverNetworkSettings.NetID
				}
				if exists == "" {
					err := logic.DeleteNodeByID(serverNode, true)
					if err != nil {
						if err == nil {
							err = errors.New("network delete failed for " + exists)
						}
						logger.Log(1, "error removing network", exists, "during sync", err.Error())
					}
				}
			}
		}
	*/
	return nil
}

// SetDefaultACLS - runs through each network to see if ACL's are set. If not, goes through each node in network and adds the default ACL
func SetDefaultACLS() error {
	// upgraded systems will not have ACL's set, which is why we need this function
	nodes, err := logic.GetAllNodes()
	if err != nil {
		return err
	}
	for i := range nodes {
		currentNodeACL, err := nodeacls.FetchNodeACL(nodeacls.NetworkID(nodes[i].Network), nodeacls.NodeID(nodes[i].ID))
		if (err != nil && (database.IsEmptyRecord(err) || strings.Contains(err.Error(), "no node ACL present"))) || currentNodeACL == nil {
			if _, err = nodeacls.CreateNodeACL(nodeacls.NetworkID(nodes[i].Network), nodeacls.NodeID(nodes[i].ID), acls.Allowed); err != nil {
				logger.Log(1, "could not create a default ACL for node", nodes[i].ID)
			}
		}
	}
	return nil
}
