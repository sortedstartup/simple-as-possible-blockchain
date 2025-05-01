package main

// TODO: make a plugin for vscode
// This generates the go code from .proto files
// For simple cases like this it avoids the need to have a Makefile

//go:generate protoc --go_opt=module=sortedfamily.com/financeservice --go-grpc_opt=module=sortedfamily.com/financeservice --go_out=. --go-grpc_out=. --proto_path=../../proto financeservice.proto

// This generates JS code from .proto files
//go:generate protoc --ts_opt=no_namespace --ts_opt=unary_rpc_promise=true --ts_opt=target=web --ts_out=../../frontend/ionic/src/proto/ --proto_path=../../proto financeservice.proto

// // This is a hack to avoid using grpc-js which is not needed in the browser
// // If we can move to connect RPC auto generation this is not needed
// // go:generate sh -c "sed -i  's|@grpc/grpc-js|grpc-web|g' ../../frontend/ionic/src/proto/financeservice.ts"

// // This to avoid any errors during `npm run build`
// // go:generate sh -c "sed -i '1i\\// @ts-nocheck' ../../frontend/ionic/src/proto/financeservice.ts"

// // go:generate sqlc -f db/scripts/sqlc.yaml generate