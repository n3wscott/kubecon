package person

type Hello struct {
	Name string `json:"name"`
}

/*

curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 0.2" \
  -H "ce-source: curl-command" \
  -H "ce-type: curl.demo" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Earl"}' \
  http://localhost:8080/


curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 0.2" \
  -H "ce-source: curl-command" \
  -H "ce-type: curl.demo" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Earl"}' \
  http://sockeye.default.sps.n3wscott.com/

*/
