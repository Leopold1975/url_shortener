//go:generate protoc --go_out=. --go-grpc_out=. ./api/shortener.v1.proto
//go:generate openapi-generator-cli generate -i ./api/shortener.v1.openapi.yaml -g go --additional-properties=withGoMod=false -o internal/shortener/server/restserver/openapi --git-repo-id url_shortener/internal/shortener/server/restserver/openapi --git-user-id Leopold1975
//go:generate  openapi-generator-cli generate -i ./api/shortener.v1.openapi.yaml -g html -o ./docsapi

package generate
