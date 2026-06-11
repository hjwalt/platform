#!/usr/bin/env sh

TARGET_BROKER=localhost:9092

kafka-topics.sh --bootstrap-server $TARGET_BROKER --delete --topic AGENT
kafka-topics.sh --bootstrap-server $TARGET_BROKER --delete --topic AGENT-RESULT

kafka-topics.sh --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic AGENT
kafka-topics.sh --bootstrap-server $TARGET_BROKER --create --partitions 10 --topic AGENT-RESULT

rm -f ./tmp/memory/*

# PGPASSWORD=postgres psql -U postgres -h localhost -p 5432 -d postgres -f script/create.sql