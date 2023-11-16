GENMOCK=make gen-mock

gen-all-mock:
	$(GENMOCK) source=kafka.go folder=mock_kafka pkg=mockkafka
	$(GENMOCK) source=client.go folder=mock_client pkg=mockclient

gen-mock:
	mockgen -source=$(source) -package=$(pkg) -destination=testings/$(folder)/$(pkg).go