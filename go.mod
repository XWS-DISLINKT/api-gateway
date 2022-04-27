module api-gateway

go 1.18

replace github.com/XWS-DISLINKT/dislinkt/common => ../common

require (
	github.com/XWS-DISLINKT/dislinkt/common v1.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	google.golang.org/grpc v1.46.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
