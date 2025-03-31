```
protoc -I. -I api/third_party --go-gin_out=. --go_out=. --go_opt=module=gomod.sunmi.com/gomoddepend/appstore-common --go-gin_opt=module=gomod.sunmi.com/gomoddepend/appstore-common api/common.proto

protoc --go-errors_out=fe_ecode=./api/docs/fe_ecode:./ api/common_ecode.proto


protoc -I. -I api/third_party --go-gin_out=. --go_out=. api/partners.proto


# 生成到openapi文件夹
for protoName in "openapi" "openapi_xyt" "partners" "partners_private" "jobs" "directedinstall" "directedinstall_group" "terminal" "mgt";
do
  if [[ $HostName == 'SMSHA1C02DF200ML85.local' ]]; then
    go run -mod=mod swagger.go api/$protoName.proto api/docs/$protoName/$protoName.swagger.proto
  else
    go run swagger.go api/$protoName.proto api/docs/$protoName/$protoName.swagger.proto
  fi
  protoc --proto_path=. \
          --proto_path=./api/third_party \
          --openapi_out=fq_schema_naming=true,default_response=false,output_mode=source_relative:. \
          api/docs/$protoName/$protoName.swagger.proto

  yq -Poj api/docs/$protoName/$protoName.swagger.yaml > api/docs/$protoName/$protoName.swagger.json
done
```