package models

import (
	//  "../mongoconn"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Network Struct
//At  some point, need to replace all instances of Name with something else like  Identifier
type Network struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AddressRange string             `json:"addressrange" bson:"addressrange" validate:"required,cidr"`
	// bug in validator --- required_with does not work with bools  issue#683
	//	AddressRange6          string             `json:"addressrange6" bson:"addressrange6" validate:"required_with=isdualstack true,cidrv6"`
	AddressRange6 string `json:"addressrange6" bson:"addressrange6" validate:"addressrange6_valid"`
	//can't have min=1 with omitempty
	DisplayName         string      `json:"displayname,omitempty" bson:"displayname,omitempty" validate:"omitempty,min=1,max=20,displayname_valid"`
	NetID               string      `json:"netid" bson:"netid" validate:"required,min=1,max=12,netid_valid"`
	NodesLastModified   int64       `json:"nodeslastmodified" bson:"nodeslastmodified"`
	NetworkLastModified int64       `json:"networklastmodified" bson:"networklastmodified"`
	DefaultInterface    string      `json:"defaultinterface" bson:"defaultinterface"`
	DefaultListenPort   int32       `json:"defaultlistenport,omitempty" bson:"defaultlistenport,omitempty" validate:"omitempty,min=1024,max=65535"`
	NodeLimit	    int32       `json:"nodelimit" bson:"nodelimit"`
	DefaultPostUp       string      `json:"defaultpostup" bson:"defaultpostup"`
	DefaultPostDown     string      `json:"defaultpostdown" bson:"defaultpostdown"`
	KeyUpdateTimeStamp  int64       `json:"keyupdatetimestamp" bson:"keyupdatetimestamp"`
	DefaultKeepalive    int32       `json:"defaultkeepalive" bson:"defaultkeepalive" validate:"omitempty,max=1000"`
	DefaultSaveConfig   *bool       `json:"defaultsaveconfig" bson:"defaultsaveconfig"`
	AccessKeys          []AccessKey `json:"accesskeys" bson:"accesskeys"`
	AllowManualSignUp   *bool       `json:"allowmanualsignup" bson:"allowmanualsignup"`
	IsLocal             *bool       `json:"islocal" bson:"islocal"`
	IsDualStack         *bool       `json:"isdualstack" bson:"isdualstack"`
	IsIPv4         string       `json:"isipv4" bson:"isipv4"`
	IsIPv6         string       `json:"isipv6" bson:"isipv6"`
	IsGRPCHub         string       `json:"isgrpchub" bson:"isgrpchub"`
	LocalRange          string      `json:"localrange" bson:"localrange" validate:"omitempty,cidr"`
	NextAvailableIP         string       `json:"nextavailableip" bson:"nextavailableip"`
	//can't have min=1 with omitempty
	DefaultCheckInInterval int32 `json:"checkininterval,omitempty" bson:"checkininterval,omitempty" validate:"omitempty,numeric,min=2,max=100000"`
}
type NetworkUpdate struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AddressRange string             `json:"addressrange" bson:"addressrange" validate:"omitempty,cidr"`

	// bug in validator --- required_with does not work with bools  issue#683
	//	AddressRange6          string             `json:"addressrange6" bson:"addressrange6" validate:"required_with=isdualstack true,cidrv6"`
	AddressRange6 string `json:"addressrange6" bson:"addressrange6" validate:"omitempty,cidr"`
	//can't have min=1 with omitempty
	DisplayName         string      `json:"displayname,omitempty" bson:"displayname,omitempty" validate:"omitempty,netid_valid,min=1,max=20"`
	NetID               string      `json:"netid" bson:"netid" validate:"omitempty,netid_valid,min=1,max=15"`
	NodesLastModified   int64       `json:"nodeslastmodified" bson:"nodeslastmodified"`
	NetworkLastModified int64       `json:"networklastmodified" bson:"networklastmodified"`
	DefaultInterface    string      `json:"defaultinterface" bson:"defaultinterface"`
	DefaultListenPort   int32       `json:"defaultlistenport,omitempty" bson:"defaultlistenport,omitempty" validate:"omitempty,min=1024,max=65535"`
	NodeLimit	    int32       `json:"nodelimit" bson:"nodelimit"`
	DefaultPostUp       string      `json:"defaultpostup" bson:"defaultpostup"`
	DefaultPostDown     string      `json:"defaultpostdown" bson:"defaultpostdown"`
	KeyUpdateTimeStamp  int64       `json:"keyupdatetimestamp" bson:"keyupdatetimestamp"`
	DefaultKeepalive    int32       `json:"defaultkeepalive" bson:"defaultkeepalive" validate:"omitempty,max=1000"`
	DefaultSaveConfig   *bool       `json:"defaultsaveconfig" bson:"defaultsaveconfig"`
	AccessKeys          []AccessKey `json:"accesskeys" bson:"accesskeys"`
	AllowManualSignUp   *bool       `json:"allowmanualsignup" bson:"allowmanualsignup"`
	IsLocal             *bool       `json:"islocal" bson:"islocal"`
	IsDualStack         *bool       `json:"isdualstack" bson:"isdualstack"`
        IsIPv4         string       `json:"isipv4" bson:"isipv4"`
        IsIPv6         string       `json:"isipv6" bson:"isipv6"`
        IsGRPCHub         string       `json:"isgrpchub" bson:"isgrpchub"`
	LocalRange          string      `json:"localrange" bson:"localrange" validate:"omitempty,cidr"`
	//can't have min=1 with omitempty
	DefaultCheckInInterval int32 `json:"checkininterval,omitempty" bson:"checkininterval,omitempty" validate:"omitempty,numeric,min=2,max=100000"`
}

//TODO:
//Not  sure if we  need the below two functions. Got rid  of one of the calls. May want  to revisit
func (network *Network) SetNodesLastModified() {
	network.NodesLastModified = time.Now().Unix()
}

func (network *Network) SetNetworkLastModified() {
	network.NetworkLastModified = time.Now().Unix()
}

func (network *Network) SetDefaults() {
	if network.DisplayName == "" {
		network.DisplayName = network.NetID
	}
	if network.DefaultInterface == "" {
		if len(network.NetID) < 13 {
			network.DefaultInterface = "nm-" + network.NetID
		} else {
			network.DefaultInterface = network.NetID
		}
	}
	if network.DefaultListenPort == 0 {
		network.DefaultListenPort = 51821
	}
        if network.NodeLimit == 0 {
                network.NodeLimit = 999999999
        }
	if network.DefaultPostDown == "" {

	}
	if network.DefaultSaveConfig == nil {
		defaultsave := true
		network.DefaultSaveConfig = &defaultsave
	}
	if network.DefaultKeepalive == 0 {
		network.DefaultKeepalive = 20
	}
	if network.DefaultPostUp == "" {
	}
	//Check-In Interval for Nodes, In Seconds
	if network.DefaultCheckInInterval == 0 {
		network.DefaultCheckInInterval = 30
	}
	if network.AllowManualSignUp == nil {
		signup := false
		network.AllowManualSignUp = &signup
	}
	if (network.IsDualStack != nil) && *network.IsDualStack {
		network.IsIPv6 = "yes"
		network.IsIPv4 = "yes"
	} else if network.IsGRPCHub != "yes" {
                network.IsIPv6 = "no"
                network.IsIPv4 = "yes"
	}
}
