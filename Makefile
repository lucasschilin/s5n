test-message: PORT ?= 8080
test-message: MESSAGE_ID ?= 1234AAB
test-message: FROM ?= @lucasschilin
test-message: BODY ?= apenas um teste
test-message:
	curl -X POST http://localhost:$(PORT)/webhook \
		-H "Content-Type: application/json" \
		-d '{"message_id":"$(MESSAGE_ID)","from":"$(FROM)","body":"$(BODY)"}'