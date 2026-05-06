#!/usr/bin/env sh

TARGET_BROKER=localhost:9092

kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic AGENT
kafka-topics --bootstrap-server $TARGET_BROKER --delete --topic AGENT-RESULT

kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic AGENT
kafka-topics --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic AGENT-RESULT

# PGPASSWORD=postgres psql -U postgres -h localhost -p 5432 -d postgres -f script/create.sql