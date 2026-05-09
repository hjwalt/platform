#!/usr/bin/env sh

TARGET_BROKER=localhost:9092

kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "agent-consumer"
kafka-consumer-groups --bootstrap-server $TARGET_BROKER --delete --group "result-consumer"
