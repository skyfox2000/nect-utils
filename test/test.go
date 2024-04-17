package main

import (
	"fmt"

	"github.com/skyfox2000/nect-utils/json"
)

func main() {

	// Sample input string
	input := `[
		{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
		{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device", "children": [
		{"Code": "KXXY001-XT01-DCC02-C077", "ContainerMode": null, "Icon": "fa-cube", "Id": "6d1580ef13864473a9ef387ef08bb94e", "Name": "电芯", "ParentId": "8329a36d7823408c962b1abb5d25a61a", "SourceId": "000ddfb70ad842dc85d7c5848239f5e4", "TargetId": "a3de3d13fe524819a1398dfce2a82dd0", "Type": "Device"},
		{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device", "children": [
			{"Code": "KXXY001-XT01-DCC02-C077", "ContainerMode": null, "Icon": "fa-cube", "Id": "6d1580ef13864473a9ef387ef08bb94e", "Name": "电芯", "ParentId": "8329a36d7823408c962b1abb5d25a61a", "SourceId": "000ddfb70ad842dc85d7c5848239f5e4", "TargetId": "a3de3d13fe524819a1398dfce2a82dd0", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"},
			{"Code": "KXXY001-XT01-DCC01-C039", "ContainerMode": null, "Icon": "fa-cube", "Id": "5c4aecca92994478b59fccd7e2d0a44f", "Name": "电芯", "ParentId": "207a699d55ea48a990a8567425e678bd", "SourceId": "004a937b5e11427f9a89917a70e04b90", "TargetId": "0de10f5ec5be434c84c1bd5695a4537a", "Type": "Device"}]}]}
	]`

	// Convert input string to JSON object
	result, _ := json.JSON.Parse(input)

	output := json.JSON.Log(result, 2, 3)
	fmt.Println(output)
}
