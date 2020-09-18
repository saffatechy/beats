#!/bin/bash -e

if test $# -eq 1 -a "$1" = setup; then
    for i in filebeat/magefile.go LICENSE.txt ; do
        test -e "$i" || (echo "Required file $i not found. Run this script from Beats repo base directory"; exit 1)
    done

    for i in filebeat.yml.test sample.log ; do
        test -e "$i" || (echo "Required file $i not found in current dir"; exit 1)
    done
    cp sample.log /tmp/
    cp filebeat.yml.test filebeat/
    exit 0
fi

ES='http://localhost:9200'

# cleanup
curl -X DELETE "$ES/filebeat-*" || true
curl -X DELETE "$ES/_template/filebeat*" || true

# build
pushd filebeat
mage update build

# setup
rm -rf data || true
cp filebeat.yml.test filebeat.yml

# run
./filebeat setup --template
./filebeat run -e &
FBPID=$!
popd

# wait for it to be populated
while true
do
    sleep 3
    N=$(curl "$ES/_cat/indices/filebeat*?h=docs.count" || echo 0)
    test "$N" -gt 0 && break
done

kill $FBPID
curl "$ES/filebeat-8.0.0/_mapping" | grep -q geo_point
