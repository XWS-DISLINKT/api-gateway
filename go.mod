module api-gateway

go 1.18

replace github.com/XWS-DISLINKT/dislinkt/common => ../common

replace github.com/XWS-DISLINKT/dislinkt/tracer => ../tracer

require (
	github.com/XWS-DISLINKT/dislinkt/common v1.0.0
	github.com/XWS-DISLINKT/dislinkt/tracer v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/opentracing/opentracing-go v1.2.0
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
)
